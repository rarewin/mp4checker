package atom

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

type Atom struct {
	size     uint32
	atype    string
	children []Atom
	elements map[string]interface{}
}

var atom_parsers map[string]func(*Atom, io.Reader) *Atom
var diff_time time.Duration

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
	var tmp32 uint32
	var tmp16 uint16

	binary.Read(r, binary.BigEndian, &tmp32)
	el["version"] = (tmp32 >> 24) & 0xff
	el["flags"] = tmp32 & 0xffffff

	binary.Read(r, binary.BigEndian, &tmp32)
	el["creation_time"] = time.Unix(int64(tmp32), 0).Add(diff_time)

	binary.Read(r, binary.BigEndian, &tmp32)
	el["modification_time"] = time.Unix(int64(tmp32), 0).Add(diff_time)

	binary.Read(r, binary.BigEndian, &tmp32)
	el["time_scale"] = tmp32

	binary.Read(r, binary.BigEndian, &tmp32)
	el["duration"] = tmp32

	binary.Read(r, binary.BigEndian, &tmp32)
	el["preferred_rate"] = tmp32

	binary.Read(r, binary.BigEndian, &tmp16)
	el["preferred_volume"] = tmp16

	// reserved
	buf := make([]byte, 10)
	r.Read(buf)

	// matrix structure
	var matrix [3][3]uint32
	binary.Read(r, binary.BigEndian, &matrix[0][0])
	binary.Read(r, binary.BigEndian, &matrix[0][1])
	binary.Read(r, binary.BigEndian, &matrix[0][2])
	binary.Read(r, binary.BigEndian, &matrix[1][0])
	binary.Read(r, binary.BigEndian, &matrix[1][1])
	binary.Read(r, binary.BigEndian, &matrix[1][2])
	binary.Read(r, binary.BigEndian, &matrix[2][0])
	binary.Read(r, binary.BigEndian, &matrix[2][1])
	binary.Read(r, binary.BigEndian, &matrix[2][2])
	el["matrix_structure"] = matrix

	// preview time
	binary.Read(r, binary.BigEndian, &tmp32)
	el["preview_time"] = tmp32

	// preview duration
	binary.Read(r, binary.BigEndian, &tmp32)
	el["preview_duration"] = tmp32

	// poster_time
	binary.Read(r, binary.BigEndian, &tmp32)
	el["poster_time"] = tmp32

	// selection time
	binary.Read(r, binary.BigEndian, &tmp32)
	el["selection_time"] = tmp32

	// selection duration
	binary.Read(r, binary.BigEndian, &tmp32)
	el["selection_duration"] = tmp32

	// current time
	binary.Read(r, binary.BigEndian, &tmp32)
	el["current_time"] = tmp32

	// next track ID
	binary.Read(r, binary.BigEndian, &tmp32)
	el["next_track_ID"] = tmp32

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

func Parse_atom(r io.Reader) []Atom {

	var atoms = make([]Atom, 0)
	var size uint32

	buf := make([]byte, 4)

	for binary.Read(r, binary.BigEndian, &size) == nil {

		atom := new(Atom)
		atom.size = size

		r.Read(buf)
		atom.atype = string(buf)

		if atom_parsers[atom.atype] != nil {
			atom_parsers[atom.atype](atom, r)
		} else {
			parse_general(atom, r)
		}

		Print_atom(atom)
		atoms = append(atoms, *atom)
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

	diff_time = time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC).Sub(time.Unix(0, 0))
}
