// Code generated by bpf2go; DO NOT EDIT.
//go:build 386 || amd64 || amd64p32 || arm || arm64 || mips64le || mips64p32le || mipsle || ppc64le || riscv64
// +build 386 amd64 amd64p32 arm arm64 mips64le mips64p32le mipsle ppc64le riscv64

package tlstapper

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

type tlsTapperGolangReadWrite struct {
	Pid       uint32
	Fd        uint32
	ConnAddr  uint32
	IsRequest bool
	Data      [524288]uint8
	_         [3]byte
}

type tlsTapperTlsChunk struct {
	Pid      uint32
	Tgid     uint32
	Len      uint32
	Start    uint32
	Recorded uint32
	Fd       uint32
	Flags    uint32
	Address  [16]uint8
	Data     [4096]uint8
}

// loadTlsTapper returns the embedded CollectionSpec for tlsTapper.
func loadTlsTapper() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_TlsTapperBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load tlsTapper: %w", err)
	}

	return spec, err
}

// loadTlsTapperObjects loads tlsTapper and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//     *tlsTapperObjects
//     *tlsTapperPrograms
//     *tlsTapperMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadTlsTapperObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadTlsTapper()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// tlsTapperSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type tlsTapperSpecs struct {
	tlsTapperProgramSpecs
	tlsTapperMapSpecs
}

// tlsTapperSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type tlsTapperProgramSpecs struct {
	GolangCryptoTlsReadUprobe   *ebpf.ProgramSpec `ebpf:"golang_crypto_tls_read_uprobe"`
	GolangCryptoTlsWriteUprobe  *ebpf.ProgramSpec `ebpf:"golang_crypto_tls_write_uprobe"`
	GolangNetHttpDialconnUprobe *ebpf.ProgramSpec `ebpf:"golang_net_http_dialconn_uprobe"`
	GolangNetSocketUprobe       *ebpf.ProgramSpec `ebpf:"golang_net_socket_uprobe"`
	SslRead                     *ebpf.ProgramSpec `ebpf:"ssl_read"`
	SslReadEx                   *ebpf.ProgramSpec `ebpf:"ssl_read_ex"`
	SslRetRead                  *ebpf.ProgramSpec `ebpf:"ssl_ret_read"`
	SslRetReadEx                *ebpf.ProgramSpec `ebpf:"ssl_ret_read_ex"`
	SslRetWrite                 *ebpf.ProgramSpec `ebpf:"ssl_ret_write"`
	SslRetWriteEx               *ebpf.ProgramSpec `ebpf:"ssl_ret_write_ex"`
	SslWrite                    *ebpf.ProgramSpec `ebpf:"ssl_write"`
	SslWriteEx                  *ebpf.ProgramSpec `ebpf:"ssl_write_ex"`
	SysEnterAccept4             *ebpf.ProgramSpec `ebpf:"sys_enter_accept4"`
	SysEnterConnect             *ebpf.ProgramSpec `ebpf:"sys_enter_connect"`
	SysEnterRead                *ebpf.ProgramSpec `ebpf:"sys_enter_read"`
	SysEnterWrite               *ebpf.ProgramSpec `ebpf:"sys_enter_write"`
	SysExitAccept4              *ebpf.ProgramSpec `ebpf:"sys_exit_accept4"`
	SysExitConnect              *ebpf.ProgramSpec `ebpf:"sys_exit_connect"`
}

// tlsTapperMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type tlsTapperMapSpecs struct {
	AcceptSyscallContext *ebpf.MapSpec `ebpf:"accept_syscall_context"`
	ChunksBuffer         *ebpf.MapSpec `ebpf:"chunks_buffer"`
	ConnectSyscallInfo   *ebpf.MapSpec `ebpf:"connect_syscall_info"`
	FileDescriptorToIpv4 *ebpf.MapSpec `ebpf:"file_descriptor_to_ipv4"`
	GolangDialWrites     *ebpf.MapSpec `ebpf:"golang_dial_writes"`
	GolangReadWrites     *ebpf.MapSpec `ebpf:"golang_read_writes"`
	GolangSocketDials    *ebpf.MapSpec `ebpf:"golang_socket_dials"`
	Heap                 *ebpf.MapSpec `ebpf:"heap"`
	LogBuffer            *ebpf.MapSpec `ebpf:"log_buffer"`
	PidsMap              *ebpf.MapSpec `ebpf:"pids_map"`
	SslReadContext       *ebpf.MapSpec `ebpf:"ssl_read_context"`
	SslWriteContext      *ebpf.MapSpec `ebpf:"ssl_write_context"`
}

// tlsTapperObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadTlsTapperObjects or ebpf.CollectionSpec.LoadAndAssign.
type tlsTapperObjects struct {
	tlsTapperPrograms
	tlsTapperMaps
}

func (o *tlsTapperObjects) Close() error {
	return _TlsTapperClose(
		&o.tlsTapperPrograms,
		&o.tlsTapperMaps,
	)
}

// tlsTapperMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadTlsTapperObjects or ebpf.CollectionSpec.LoadAndAssign.
type tlsTapperMaps struct {
	AcceptSyscallContext *ebpf.Map `ebpf:"accept_syscall_context"`
	ChunksBuffer         *ebpf.Map `ebpf:"chunks_buffer"`
	ConnectSyscallInfo   *ebpf.Map `ebpf:"connect_syscall_info"`
	FileDescriptorToIpv4 *ebpf.Map `ebpf:"file_descriptor_to_ipv4"`
	GolangDialWrites     *ebpf.Map `ebpf:"golang_dial_writes"`
	GolangReadWrites     *ebpf.Map `ebpf:"golang_read_writes"`
	GolangSocketDials    *ebpf.Map `ebpf:"golang_socket_dials"`
	Heap                 *ebpf.Map `ebpf:"heap"`
	LogBuffer            *ebpf.Map `ebpf:"log_buffer"`
	PidsMap              *ebpf.Map `ebpf:"pids_map"`
	SslReadContext       *ebpf.Map `ebpf:"ssl_read_context"`
	SslWriteContext      *ebpf.Map `ebpf:"ssl_write_context"`
}

func (m *tlsTapperMaps) Close() error {
	return _TlsTapperClose(
		m.AcceptSyscallContext,
		m.ChunksBuffer,
		m.ConnectSyscallInfo,
		m.FileDescriptorToIpv4,
		m.GolangDialWrites,
		m.GolangReadWrites,
		m.GolangSocketDials,
		m.Heap,
		m.LogBuffer,
		m.PidsMap,
		m.SslReadContext,
		m.SslWriteContext,
	)
}

// tlsTapperPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadTlsTapperObjects or ebpf.CollectionSpec.LoadAndAssign.
type tlsTapperPrograms struct {
	GolangCryptoTlsReadUprobe   *ebpf.Program `ebpf:"golang_crypto_tls_read_uprobe"`
	GolangCryptoTlsWriteUprobe  *ebpf.Program `ebpf:"golang_crypto_tls_write_uprobe"`
	GolangNetHttpDialconnUprobe *ebpf.Program `ebpf:"golang_net_http_dialconn_uprobe"`
	GolangNetSocketUprobe       *ebpf.Program `ebpf:"golang_net_socket_uprobe"`
	SslRead                     *ebpf.Program `ebpf:"ssl_read"`
	SslReadEx                   *ebpf.Program `ebpf:"ssl_read_ex"`
	SslRetRead                  *ebpf.Program `ebpf:"ssl_ret_read"`
	SslRetReadEx                *ebpf.Program `ebpf:"ssl_ret_read_ex"`
	SslRetWrite                 *ebpf.Program `ebpf:"ssl_ret_write"`
	SslRetWriteEx               *ebpf.Program `ebpf:"ssl_ret_write_ex"`
	SslWrite                    *ebpf.Program `ebpf:"ssl_write"`
	SslWriteEx                  *ebpf.Program `ebpf:"ssl_write_ex"`
	SysEnterAccept4             *ebpf.Program `ebpf:"sys_enter_accept4"`
	SysEnterConnect             *ebpf.Program `ebpf:"sys_enter_connect"`
	SysEnterRead                *ebpf.Program `ebpf:"sys_enter_read"`
	SysEnterWrite               *ebpf.Program `ebpf:"sys_enter_write"`
	SysExitAccept4              *ebpf.Program `ebpf:"sys_exit_accept4"`
	SysExitConnect              *ebpf.Program `ebpf:"sys_exit_connect"`
}

func (p *tlsTapperPrograms) Close() error {
	return _TlsTapperClose(
		p.GolangCryptoTlsReadUprobe,
		p.GolangCryptoTlsWriteUprobe,
		p.GolangNetHttpDialconnUprobe,
		p.GolangNetSocketUprobe,
		p.SslRead,
		p.SslReadEx,
		p.SslRetRead,
		p.SslRetReadEx,
		p.SslRetWrite,
		p.SslRetWriteEx,
		p.SslWrite,
		p.SslWriteEx,
		p.SysEnterAccept4,
		p.SysEnterConnect,
		p.SysEnterRead,
		p.SysEnterWrite,
		p.SysExitAccept4,
		p.SysExitConnect,
	)
}

func _TlsTapperClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//go:embed tlstapper_bpfel.o
var _TlsTapperBytes []byte
