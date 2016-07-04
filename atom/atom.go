package atom

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

type AtomType interface {
	Size() uint32
	Type() string
	Print()
}

type Atom struct {
	size     uint32
	atype    string
	children []*AtomType
	Print    func()
	Type     func()
}

var atom_parsers map[string]func(Atom, io.Reader) *AtomType

// ftyp
type ftypAtom struct {
	Atom
	major_brand       string
	minor_version     string
	compatible_brands []string
}

func (fa *ftypAtom) Print() {

	fmt.Printf("size: %d\n", fa.size)
	fmt.Printf("type: %s\n", fa.atype)
	fmt.Printf("Major_brand: %s\n", fa.major_brand)
	fmt.Printf("Minor_version: %s\n", fa.minor_version)

	for _, v := range fa.compatible_brands {
		fmt.Printf("Compatible_brands: %s\n", v)
	}
}

func (fa *ftypAtom) Size() uint32 {
	return fa.size
}

func (fa *ftypAtom) Type() string {
	return fa.atype
}

func parse_ftyp(a Atom, r io.Reader) *AtomType {

	var ftyp ftypAtom

	ftyp.size = a.size
	ftyp.atype = a.atype

	buf := make([]byte, 4)

	r.Read(buf)
	ftyp.major_brand = string(buf)

	r.Read(buf)
	ftyp.minor_version = string(buf)

	remain_size := ftyp.size - 16

	for i := uint32(0); i < remain_size; i += 4 {
		r.Read(buf)
		ftyp.compatible_brands = append(ftyp.compatible_brands, string(buf))
	}

	var ret AtomType = &ftyp

	return &ret
}

// moov
type moovAtom struct {
	Atom
}

func (ma *moovAtom) Print() {
	fmt.Printf("size: %d\n", ma.size)
	fmt.Printf("type: %s\n", ma.atype)
}

func (ma *moovAtom) Size() uint32 {
	return ma.size
}

func (ma *moovAtom) Type() string {
	return ma.atype
}

func parse_moov(a Atom, r io.Reader) *AtomType {

	var moov moovAtom

	moov.size = a.size
	moov.atype = a.atype

	moov.children = Parse_atom(r)

	var ret AtomType = &moov

	return &ret
}

func (a *Atom) Parse(r io.Reader) {

}

// mvhd
type mvhdAtom struct {
	Atom
	version            uint32
	flags              uint32
	creation_time      time.Time
	modification_time  time.Time
	time_scale         uint32
	duration           uint32
	preferred_rate     uint32
	matrix_structure   [3][3]uint32
	preview_time       time.Time
	preview_duration   uint32
	poster_time        time.Time
	selection_time     time.Time
	selection_duration uint32
	current_time       time.Time
	next_track_id      uint32
}

func (ma *mvhdAtom) Print() {
	fmt.Printf("size: %d\n", ma.size)
	fmt.Printf("type: %s\n", ma.atype)
	fmt.Printf("version: %d\n", ma.version)
	fmt.Printf("flags: %x\n", ma.flags)
}

func (ma *mvhdAtom) Size() uint32 {
	return ma.size
}

func (ma *mvhdAtom) Type() string {
	return ma.atype
}

func parse_mvhd(a Atom, r io.Reader) *AtomType {
	var mvhd mvhdAtom

	mvhd.size = a.size
	mvhd.atype = a.atype

	var tmp uint32
	binary.Read(r, binary.LittleEndian, &tmp)
	mvhd.version = (tmp >> 24) & 0xff
	mvhd.flags = tmp & 0xffffff

	fmt.Printf("(%s) %d bytes were ignored\n", a.atype, a.size-12)

	buf := make([]byte, a.size-12)
	r.Read(buf)

	var ret AtomType = &mvhd

	return &ret
}

// free
type freeAtom struct {
	Atom
}

func (fa *freeAtom) Print() {
	fmt.Printf("size: %d\n", fa.size)
	fmt.Printf("type: %s\n", fa.atype)
}

func (fa *freeAtom) Size() uint32 {
	return fa.size
}

func (fa *freeAtom) Type() string {
	return fa.atype
}

func parse_free(a Atom, r io.Reader) *AtomType {
	var free freeAtom

	free.size = a.size
	free.atype = a.atype

	var ret AtomType = &free

	return &ret
}

// general
type generalAtom struct {
	Atom
}

func (ga *generalAtom) Print() {
	fmt.Printf("size: %d\n", ga.size)
	fmt.Printf("type: %s\n", ga.atype)
}

func (ga *generalAtom) Size() uint32 {
	return ga.size
}

func (ga *generalAtom) Type() string {
	return ga.atype
}

func parse_general(a Atom, r io.Reader) *AtomType {

	var ga generalAtom

	ga.size = a.size
	ga.atype = a.atype

	// skip
	fmt.Printf("(%s) %d bytes were ignored\n", a.atype, a.size-8)
	buf := make([]byte, a.size-8)
	r.Read(buf)
	//

	var ret AtomType = &ga

	return &ret

}

func Parse_atom(r io.Reader) []*AtomType {

	var atom Atom

	var atoms = make([]*AtomType, 0)

	buf := make([]byte, 4)

	for binary.Read(r, binary.BigEndian, &atom.size) == nil {

		r.Read(buf)
		atom.atype = string(buf)

		if atom_parsers[atom.atype] != nil {

			mp4 := *atom_parsers[atom.atype](atom, r)
			mp4.Print()

			atoms = append(atoms, &mp4)

		} else {

			mp4 := *parse_general(atom, r)
			mp4.Print()

			atoms = append(atoms, &mp4)
		}
	}

	return atoms
}

func init() {

	atom_parsers = map[string]func(Atom, io.Reader) *AtomType{
		"ftyp": parse_ftyp,
		"moov": parse_moov,
		"mvhd": parse_mvhd,
		"free": parse_free,
	}
}
