package atom

import (
	"encoding/binary"
	"fmt"
	"io"
)

var atom_parsers = map[string]func(r io.Reader) (a Atom){
	"ftyp": parse_ftyp,
}

type Atom struct {
	size  uint32
	atype string
}

func parse_ftyp(r io.Reader) (a Atom) {

	var ftyp Atom

	return ftyp
}

func Print_atom(a Atom) {

	fmt.Printf("size: %d\n", a.size)
	fmt.Printf("type: %s\n", a.atype)
}

func (a *Atom) Parse(r io.Reader) {

}

func Parse_atom(r io.Reader) {

	var atom Atom

	buf := make([]byte, 4)

	binary.Read(r, binary.BigEndian, &atom.size)

	r.Read(buf)
	atom.atype = string(buf)

	Print_atom(atom)
}
