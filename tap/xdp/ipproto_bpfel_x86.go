// Code generated by bpf2go; DO NOT EDIT.
//go:build 386 || amd64
// +build 386 amd64

package ebpf

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

// loadIpproto returns the embedded CollectionSpec for ipproto.
func loadIpproto() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_IpprotoBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load ipproto: %w", err)
	}

	return spec, err
}

// loadIpprotoObjects loads ipproto and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//     *ipprotoObjects
//     *ipprotoPrograms
//     *ipprotoMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadIpprotoObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadIpproto()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// ipprotoSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type ipprotoSpecs struct {
	ipprotoProgramSpecs
	ipprotoMapSpecs
}

// ipprotoSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type ipprotoProgramSpecs struct {
	XdpSockProg *ebpf.ProgramSpec `ebpf:"xdp_sock_prog"`
}

// ipprotoMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type ipprotoMapSpecs struct {
	QidconfMap *ebpf.MapSpec `ebpf:"qidconf_map"`
	XsksMap    *ebpf.MapSpec `ebpf:"xsks_map"`
}

// ipprotoObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadIpprotoObjects or ebpf.CollectionSpec.LoadAndAssign.
type ipprotoObjects struct {
	ipprotoPrograms
	ipprotoMaps
}

func (o *ipprotoObjects) Close() error {
	return _IpprotoClose(
		&o.ipprotoPrograms,
		&o.ipprotoMaps,
	)
}

// ipprotoMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadIpprotoObjects or ebpf.CollectionSpec.LoadAndAssign.
type ipprotoMaps struct {
	QidconfMap *ebpf.Map `ebpf:"qidconf_map"`
	XsksMap    *ebpf.Map `ebpf:"xsks_map"`
}

func (m *ipprotoMaps) Close() error {
	return _IpprotoClose(
		m.QidconfMap,
		m.XsksMap,
	)
}

// ipprotoPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadIpprotoObjects or ebpf.CollectionSpec.LoadAndAssign.
type ipprotoPrograms struct {
	XdpSockProg *ebpf.Program `ebpf:"xdp_sock_prog"`
}

func (p *ipprotoPrograms) Close() error {
	return _IpprotoClose(
		p.XdpSockProg,
	)
}

func _IpprotoClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//go:embed ipproto_bpfel_x86.o
var _IpprotoBytes []byte
