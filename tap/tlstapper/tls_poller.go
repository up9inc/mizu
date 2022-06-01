package tlstapper

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"sync"
	"time"
	"unsafe"

	"encoding/binary"
	"encoding/hex"
	"os"
	"strconv"
	"strings"

	"github.com/cilium/ebpf/perf"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/go-errors/errors"
	"github.com/hashicorp/golang-lru/simplelru"
	"github.com/up9inc/mizu/logger"
	"github.com/up9inc/mizu/tap/api"
	orderedmap "github.com/wk8/go-ordered-map"
)

const (
	fdCachedItemAvgSize = 40
	fdCacheMaxItems     = 500000 / fdCachedItemAvgSize
	golangMapLimit      = 1 << 10 // 1024
)

type tlsPoller struct {
	tls                *TlsTapper
	readers            map[string]*tlsReader
	closedReaders      chan string
	reqResMatcher      api.RequestResponseMatcher
	chunksReader       *perf.Reader
	golangReader       *ringbuf.Reader
	golangReadWriteMap *orderedmap.OrderedMap
	extension          *api.Extension
	procfs             string
	pidToNamespace     sync.Map
	fdCache            *simplelru.LRU // Actual typs is map[string]addressPair
	evictedCounter     int
}

func newTlsPoller(tls *TlsTapper, extension *api.Extension, procfs string) (*tlsPoller, error) {
	poller := &tlsPoller{
		tls:           tls,
		readers:       make(map[string]*tlsReader),
		closedReaders: make(chan string, 100),
		reqResMatcher: extension.Dissector.NewResponseRequestMatcher(),
		extension:     extension,
		chunksReader:  nil,
		procfs:        procfs,
	}

	fdCache, err := simplelru.NewLRU(fdCacheMaxItems, poller.fdCacheEvictCallback)

	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	poller.fdCache = fdCache
	return poller, nil
}

func (p *tlsPoller) init(bpfObjects *tlsTapperObjects, bufferSize int) error {
	var err error

	p.chunksReader, err = perf.NewReader(bpfObjects.ChunksBuffer, bufferSize)

	if err != nil {
		return errors.Wrap(err, 0)
	}

	p.golangReader, err = ringbuf.NewReader(bpfObjects.GolangReadWrites)

	if err != nil {
		return errors.Wrap(err, 0)
	}

	p.golangReadWriteMap = orderedmap.New()

	return nil
}

func (p *tlsPoller) close() error {
	return p.chunksReader.Close()
}

func (p *tlsPoller) pollSsllib(emitter api.Emitter, options *api.TrafficFilteringOptions, streamsMap api.TcpStreamMap) {
	chunks := make(chan *tlsChunk)

	go p.pollChunksPerfBuffer(chunks)

	for {
		select {
		case chunk, ok := <-chunks:
			if !ok {
				return
			}

			if err := p.handleTlsChunk(chunk, p.extension, emitter, options, streamsMap); err != nil {
				LogError(err)
			}
		case key := <-p.closedReaders:
			delete(p.readers, key)
		}
	}
}

func (p *tlsPoller) pollGolang(emitter api.Emitter, options *api.TrafficFilteringOptions, streamsMap api.TcpStreamMap) {
	go p.pollGolangReadWrite(p.golangReader, emitter, options, streamsMap)
}

func (p *tlsPoller) pollGolangReadWrite(rd *ringbuf.Reader, emitter api.Emitter, options *api.TrafficFilteringOptions,
	streamsMap api.TcpStreamMap) {
	nativeEndian := p.getByteOrder()
	// tlsTapperGolangReadWrite is generated by bpf2go.
	var b tlsTapperGolangReadWrite
	for {
		record, err := rd.Read()
		if err != nil {
			if errors.Is(err, ringbuf.ErrClosed) {
				log.Println("received signal, exiting..")
				return
			}
			log.Printf("reading from reader: %s", err)
			continue
		}

		// Parse the ringbuf event entry into a tlsTapperGolangReadWrite structure.
		if err := binary.Read(bytes.NewBuffer(record.RawSample), nativeEndian, &b); err != nil {
			log.Printf("parsing ringbuf event: %s", err)
			continue
		}

		if p.golangReadWriteMap.Len()+1 > golangMapLimit {
			pair := p.golangReadWriteMap.Oldest()
			p.golangReadWriteMap.Delete(pair.Key)
		}

		pid := uint64(b.Pid)
		identifier := pid<<32 + uint64(b.ConnAddr)

		var connection *golangConnection
		var _connection interface{}
		var ok bool
		if _connection, ok = p.golangReadWriteMap.Get(identifier); !ok {
			tlsEmitter := &tlsEmitter{
				delegate:  emitter,
				namespace: p.getNamespace(b.Pid),
			}

			connection = NewGolangConnection(b.Pid, b.ConnAddr, p.extension, tlsEmitter)
			p.golangReadWriteMap.Set(identifier, connection)
			streamsMap.Store(streamsMap.NextId(), connection.Stream)
		} else {
			connection = _connection.(*golangConnection)
		}

		if b.IsGzipChunk {
			connection.Gzipped = true
		}

		if b.IsRequest {
			err := connection.setAddressBySockfd(p.procfs, b.Pid, b.Fd)
			if err != nil {
				log.Printf("Error resolving address pair from fd: %s", err)
				continue
			}

			tcpid := p.buildTcpId(&connection.AddressPair)
			connection.ClientReader.tcpID = &tcpid
			connection.ServerReader.tcpID = &api.TcpID{
				SrcIP:   connection.ClientReader.tcpID.DstIP,
				DstIP:   connection.ClientReader.tcpID.SrcIP,
				SrcPort: connection.ClientReader.tcpID.DstPort,
				DstPort: connection.ClientReader.tcpID.SrcPort,
			}

			go dissect(p.extension, connection.ClientReader, options)
			go dissect(p.extension, connection.ServerReader, options)

			request := make([]byte, len(b.Data[:]))
			copy(request, b.Data[:])
			connection.ClientReader.send(request)
		} else {
			response := make([]byte, len(b.Data[:]))
			copy(response, b.Data[:])
			connection.ServerReader.send(response)
		}
	}
}

func (p *tlsPoller) pollChunksPerfBuffer(chunks chan<- *tlsChunk) {
	logger.Log.Infof("Start polling for tls events")

	for {
		record, err := p.chunksReader.Read()

		if err != nil {
			close(chunks)

			if errors.Is(err, perf.ErrClosed) {
				return
			}

			LogError(errors.Errorf("Error reading chunks from tls perf, aborting TLS! %v", err))
			return
		}

		if record.LostSamples != 0 {
			logger.Log.Infof("Buffer is full, dropped %d chunks", record.LostSamples)
			continue
		}

		buffer := bytes.NewReader(record.RawSample)

		var chunk tlsChunk

		if err := binary.Read(buffer, binary.LittleEndian, &chunk); err != nil {
			LogError(errors.Errorf("Error parsing chunk %v", err))
			continue
		}

		chunks <- &chunk
	}
}

func (p *tlsPoller) handleTlsChunk(chunk *tlsChunk, extension *api.Extension, emitter api.Emitter,
	options *api.TrafficFilteringOptions, streamsMap api.TcpStreamMap) error {
	address, err := p.getSockfdAddressPair(chunk)

	if err != nil {
		address, err = chunk.getAddressPair()

		if err != nil {
			return err
		}
	}

	key := buildTlsKey(address)
	reader, exists := p.readers[key]

	if !exists {
		reader = p.startNewTlsReader(chunk, &address, key, emitter, extension, options, streamsMap)
		p.readers[key] = reader
	}

	reader.newChunk(chunk)

	if os.Getenv("MIZU_VERBOSE_TLS_TAPPER") == "true" {
		p.logTls(chunk, key, reader)
	}

	return nil
}

func (p *tlsPoller) startNewTlsReader(chunk *tlsChunk, address *addressPair, key string,
	emitter api.Emitter, extension *api.Extension, options *api.TrafficFilteringOptions,
	streamsMap api.TcpStreamMap) *tlsReader {

	tcpid := p.buildTcpId(address)

	doneHandler := func(r *tlsReader) {
		p.closeReader(key, r)
	}

	tlsEmitter := &tlsEmitter{
		delegate:  emitter,
		namespace: p.getNamespace(chunk.Pid),
	}

	reader := &tlsReader{
		key:           key,
		chunks:        make(chan *tlsChunk, 1),
		doneHandler:   doneHandler,
		progress:      &api.ReadProgress{},
		tcpID:         &tcpid,
		isClient:      chunk.isRequest(),
		captureTime:   time.Now(),
		extension:     extension,
		emitter:       tlsEmitter,
		counterPair:   &api.CounterPair{},
		reqResMatcher: p.reqResMatcher,
	}

	stream := &tlsStream{
		reader: reader,
	}
	streamsMap.Store(streamsMap.NextId(), stream)

	reader.parent = stream

	go dissect(extension, reader, options)
	return reader
}

func dissect(extension *api.Extension, reader api.TcpReader, options *api.TrafficFilteringOptions) {
	b := bufio.NewReader(reader)

	err := extension.Dissector.Dissect(b, reader, options)

	if err != nil {
		logger.Log.Warningf("Error dissecting TLS %v - %v", reader.GetTcpID(), err)
	}
}

func (p *tlsPoller) closeReader(key string, r *tlsReader) {
	close(r.chunks)
	p.closedReaders <- key
}

func (p *tlsPoller) getSockfdAddressPair(chunk *tlsChunk) (addressPair, error) {
	address, err := getAddressBySockfd(p.procfs, chunk.Pid, chunk.Fd)
	fdCacheKey := fmt.Sprintf("%d:%d", chunk.Pid, chunk.Fd)

	if err == nil {
		if !chunk.isRequest() {
			switchedAddress := addressPair{
				srcIp:   address.dstIp,
				srcPort: address.dstPort,
				dstIp:   address.srcIp,
				dstPort: address.srcPort,
			}
			p.fdCache.Add(fdCacheKey, switchedAddress)
			return switchedAddress, nil
		} else {
			p.fdCache.Add(fdCacheKey, address)
			return address, nil
		}
	}

	fromCacheIfc, ok := p.fdCache.Get(fdCacheKey)

	if !ok {
		return addressPair{}, err
	}

	fromCache, ok := fromCacheIfc.(addressPair)

	if !ok {
		return address, errors.Errorf("Unable to cast %T to addressPair", fromCacheIfc)
	}

	return fromCache, nil
}

func buildTlsKey(address addressPair) string {
	return fmt.Sprintf("%s:%d>%s:%d", address.srcIp, address.srcPort, address.dstIp, address.dstPort)
}

func (p *tlsPoller) buildTcpId(address *addressPair) api.TcpID {
	return api.TcpID{
		SrcIP:   address.srcIp.String(),
		DstIP:   address.dstIp.String(),
		SrcPort: strconv.FormatUint(uint64(address.srcPort), 10),
		DstPort: strconv.FormatUint(uint64(address.dstPort), 10),
		Ident:   "",
	}
}

func (p *tlsPoller) addPid(pid uint32, namespace string) {
	p.pidToNamespace.Store(pid, namespace)
}

func (p *tlsPoller) getNamespace(pid uint32) string {
	namespaceIfc, ok := p.pidToNamespace.Load(pid)

	if !ok {
		return api.UNKNOWN_NAMESPACE
	}

	namespace, ok := namespaceIfc.(string)

	if !ok {
		return api.UNKNOWN_NAMESPACE
	}

	return namespace
}

func (p *tlsPoller) clearPids() {
	p.pidToNamespace.Range(func(key, v interface{}) bool {
		p.pidToNamespace.Delete(key)
		return true
	})
}

func (p *tlsPoller) logTls(chunk *tlsChunk, key string, reader *tlsReader) {
	var flagsStr string

	if chunk.isClient() {
		flagsStr = "C"
	} else {
		flagsStr = "S"
	}

	if chunk.isRead() {
		flagsStr += "R"
	} else {
		flagsStr += "W"
	}

	str := strings.ReplaceAll(strings.ReplaceAll(string(chunk.Data[0:chunk.Recorded]), "\n", " "), "\r", "")

	logger.Log.Infof("[%-44s] %s #%-4d (fd: %d) (recorded %d/%d:%d) - %s - %s",
		key, flagsStr, reader.seenChunks, chunk.Fd,
		chunk.Recorded, chunk.Len, chunk.Start,
		str, hex.EncodeToString(chunk.Data[0:chunk.Recorded]))
}

func (p *tlsPoller) fdCacheEvictCallback(key interface{}, value interface{}) {
	p.evictedCounter = p.evictedCounter + 1

	if p.evictedCounter%1000000 == 0 {
		logger.Log.Infof("Tls fdCache evicted %d items", p.evictedCounter)
	}
}

func (p *tlsPoller) getByteOrder() (byteOrder binary.ByteOrder) {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		byteOrder = binary.LittleEndian
	case [2]byte{0xAB, 0xCD}:
		byteOrder = binary.BigEndian
	default:
		panic("Could not determine native endianness.")
	}

	return
}
