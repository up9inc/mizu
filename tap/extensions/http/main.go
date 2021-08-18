package main

import (
	"bufio"
	"fmt"
	"io"
	"log"

	"github.com/up9inc/mizu/tap/api"
)

var requestCounter uint
var responseCounter uint

func init() {
	log.Println("Initializing HTTP extension.")
	requestCounter = 0
	responseCounter = 0
}

type dissecting string

func (d dissecting) Register(extension *api.Extension) {
	extension.Name = "http"
	extension.OutboundPorts = []string{"80", "8080", "443"}
	extension.InboundPorts = []string{}
}

func (d dissecting) Ping() {
	log.Printf("pong HTTP\n")
}

func (d dissecting) Dissect(b *bufio.Reader, isClient bool, tcpID *api.TcpID, Emit func(reqResPair *api.RequestResponsePair)) {
	ident := fmt.Sprintf("%s->%s:%s->%s", tcpID.SrcIP, tcpID.DstIP, tcpID.SrcPort, tcpID.DstPort)
	isHTTP2, err := checkIsHTTP2Connection(b, isClient)
	if err != nil {
		SilentError("HTTP/2-Prepare-Connection", "stream %s Failed to check if client is HTTP/2: %s (%v,%+v)", ident, err, err, err)
		// Do something?
	}

	var grpcAssembler *GrpcAssembler
	if isHTTP2 {
		err := prepareHTTP2Connection(b, isClient)
		if err != nil {
			SilentError("HTTP/2-Prepare-Connection-After-Check", "stream %s error: %s (%v,%+v)", ident, err, err, err)
		}
		grpcAssembler = createGrpcAssembler(b)
	}

	for {
		if isHTTP2 {
			err = handleHTTP2Stream(grpcAssembler, tcpID, Emit)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			} else if err != nil {
				SilentError("HTTP/2", "stream %s error: %s (%v,%+v)", ident, err, err, err)
				continue
			}
		} else if isClient {
			err = handleHTTP1ClientStream(b, tcpID, Emit)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			} else if err != nil {
				SilentError("HTTP-request", "stream %s Request error: %s (%v,%+v)", ident, err, err, err)
				continue
			}
		} else {
			err = handleHTTP1ServerStream(b, tcpID, Emit)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			} else if err != nil {
				SilentError("HTTP-response", "stream %s Response error: %s (%v,%+v)", ident, err, err, err)
				continue
			}
		}
	}
}

var Dissector dissecting
