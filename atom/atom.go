package atom

import (
	"encoding/binary"
	"fmt"
	"io"
	//	"time"
)

type Atom struct {
	size     uint32
	atype    string
	children []*Atom
	elements map[string]interface{}
}

var atom_parsers map[string]func(*Atom, io.Reader) *Atom

// ftyp
func parse_ftyp(a *Atom, r io.Reader) *Atom {

	buf := make([]byte, 4)
	el := make(map[string]interface{})

	r.Read(buf)
	el["major_brand"] = string(buf)

	r.Read(buf)
	el["minor_version"] = string(buf)

	remain_size := a.size - 16
	tmp := make([]string, remain_size/4+1)

	for i := uint32(0); i < remain_size; i += 4 {
		r.Read(buf)
		tmp = append(tmp, string(buf))
	}
	el["compatible_brands"] = tmp

	a.elements = el

	return a
}

// moov
func parse_moov(a *Atom, r io.Reader) *Atom {

	a.children = Parse_atom(r)

	return a
}

// mvhd
func parse_mvhd(a *Atom, r io.Reader) *Atom {

	el := make(map[string]interface{})
	var tmp uint32

	binary.Read(r, binary.LittleEndian, &tmp)
	el["version"] = (tmp >> 24) & 0xff
	el["flags"] = tmp & 0xffffff

	fmt.Printf("(%s) %d bytes were ignored\n", a.atype, a.size-12)

	buf := make([]byte, a.size-12)
	r.Read(buf)

	a.elements = el

	return a
}

// free
func parse_free(a *Atom, r io.Reader) *Atom {

	fmt.Printf("(%s) %d bytes were ignored\n", a.atype, a.size-8)
	buf := make([]byte, a.size-8)
	r.Read(buf)

	return a
}

// general
func parse_general(a *Atom, r io.Reader) *Atom {

	// skip
	fmt.Printf("(%s) %d bytes were ignored\n", a.atype, a.size-8)
	buf := make([]byte, a.size-8)
	r.Read(buf)
	//

	return a
}

func Print_atom(a *Atom) {

	fmt.Printf("type: %s\n", a.atype)
	fmt.Printf("size: %d\n", a.size)

	for k, v := range a.elements {
		fmt.Printf("%s: %v\n", k, v)
	}

	fmt.Printf("\n")
}

func Parse_atom(r io.Reader) []*Atom {

	var atom Atom

	var atoms = make([]*Atom, 0)

	buf := make([]byte, 4)

	for binary.Read(r, binary.BigEndian, &atom.size) == nil {

		var mp4 *Atom

		r.Read(buf)
		atom.atype = string(buf)

		if atom_parsers[atom.atype] != nil {
			mp4 = atom_parsers[atom.atype](&atom, r)
		} else {
			mp4 = parse_general(&atom, r)
		}

		Print_atom(mp4)
		atoms = append(atoms, mp4)
	}

	return atoms
}

func init() {

	atom_parsers = map[string]func(*Atom, io.Reader) *Atom{
		"ftyp": parse_ftyp,
		"moov": parse_moov,
		"mvhd": parse_mvhd,
		"free": parse_free,
	}
}
