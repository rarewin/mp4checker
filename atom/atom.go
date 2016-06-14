package atom

import (
	"encoding/binary"
	"fmt"
	"io"
)

var atom_parsers = map[string]func(Atom, io.Reader) ftypAtom{
	"ftyp": parse_ftyp,
}

type Atom struct {
	size  uint32
	atype string
	Print func()
}

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

func parse_ftyp(a Atom, r io.Reader) ftypAtom {

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

	return ftyp
}

func (a *Atom) Parse(r io.Reader) {

}

func Parse_atom(r io.Reader) {

	var atom Atom

	buf := make([]byte, 4)

	binary.Read(r, binary.BigEndian, &atom.size)

	r.Read(buf)
	atom.atype = string(buf)

	mp4 := atom_parsers[atom.atype](atom, r)

	mp4.Print()
}
