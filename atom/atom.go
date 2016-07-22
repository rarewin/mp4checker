package atom

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"time"
)

type Atom struct {
	size     uint32
	atype    string
	level    int
	parent   *Atom
	children []Atom
	elements map[string]interface{}
}

func (self Atom) Print() {

	indent := strings.Repeat(" ", self.level)

	fmt.Printf("%stype: %s\n", indent, self.atype)
	fmt.Printf("%ssize: %d\n", indent, self.size)

	for k, v := range self.elements {

		fmt.Printf("%s%s: ", indent, k)

		switch k {

		case "minor_version":
			fmt.Printf("%#08x\n", v)

		case "flags":
			fmt.Printf("%x\n", v)

		case "track_width", "track_height":
			fmt.Printf("%.4f\n", (float32(v.(uint32)) / 65536.0))

		case "matrix_structure":
			matrix := v.([3][3]uint32)

			fmt.Printf("\n")

			for i := 0; i < 3; i++ {

				fmt.Printf("\t| %3.4f  %3.4f  %3.4f |\n",
					(float32(matrix[i][0]) / 65536.0),
					(float32(matrix[i][1]) / 65536.0),
					(float32(matrix[i][2]) / 1073741824.0))

			}

		default:
			fmt.Printf("%v\n", v)
		}
	}

	fmt.Printf("\n")

	for i := 0; i < len(self.children); i++ {
		self.children[i].Print()
	}
}

var atom_parsers map[string]func(*Atom, io.Reader) *Atom
var diff_time time.Duration

// read 32-bit value
func read32(r io.Reader) uint32 {

	var tmp32 uint32

	binary.Read(r, binary.BigEndian, &tmp32)

	return tmp32
}

// read 16-bit value
func read16(r io.Reader) uint16 {

	var tmp16 uint16

	binary.Read(r, binary.BigEndian, &tmp16)

	return tmp16
}

// ftyp
func parse_ftyp(a *Atom, r io.Reader) *Atom {

	buf := make([]byte, 4)
	el := make(map[string]interface{})

	r.Read(buf)
	el["major_brand"] = string(buf)
	el["minor_version"] = read32(r)

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

// mvhd
func parse_mvhd(a *Atom, r io.Reader) *Atom {

	el := make(map[string]interface{})
	var tmp32 uint32

	tmp32 = read32(r)
	el["version"] = (tmp32 >> 24) & 0xff
	el["flags"] = tmp32 & 0xffffff

	el["creation_time"] = time.Unix(int64(read32(r)), 0).Add(diff_time)
	el["modification_time"] = time.Unix(int64(read32(r)), 0).Add(diff_time)
	el["time_scale"] = read32(r)
	el["duration"] = read32(r)
	el["preferred_rate"] = read32(r)
	el["preferred_volume"] = read16(r)

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

	el["preview_time"] = read32(r)
	el["preview_duration"] = read32(r)
	el["poster_time"] = read32(r)
	el["selection_time"] = read32(r)
	el["selection_duration"] = read32(r)
	el["current_time"] = read32(r)
	el["next_track_ID"] = read32(r)

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

// tkhd atom
func parse_tkhd(a *Atom, r io.Reader) *Atom {

	el := make(map[string]interface{})
	var tmp32 uint32

	tmp32 = read32(r)
	el["version"] = (tmp32 >> 24) & 0xff
	el["flags"] = tmp32 & 0xffffff

	el["creation_time"] = time.Unix(int64(read32(r)), 0).Add(diff_time)
	el["modification_time"] = time.Unix(int64(read32(r)), 0).Add(diff_time)

	// track ID
	el["track_ID"] = read32(r)

	// Reserved(32 bits)
	read32(r)

	// Duration
	el["duration"] = read32(r)

	// reserved
	buf := make([]byte, 8)
	r.Read(buf)

	// Layer
	el["layer"] = read16(r)
	el["alternate_group"] = read16(r)
	el["volume"] = read16(r)

	// Reserved
	read16(r)

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

	el["track_width"] = read32(r)
	el["track_height"] = read32(r)

	a.elements = el

	return a
}

// elst atom
func parse_elst(a *Atom, r io.Reader) *Atom {

	el := make(map[string]interface{})
	var tmp32 uint32

	tmp32 = read32(r)
	el["version"] = (tmp32 >> 24) & 0xff
	el["flags"] = tmp32 & 0xffffff

	el["number_of_entries"] = read32(r)

	entries := make([][3]uint32, el["number_of_entries"].(uint32))
	for i := uint32(0); i < el["number_of_entries"].(uint32); i++ {

		entries[i][0] = read32(r)
		entries[i][1] = read32(r)
		entries[i][2] = read32(r)

	}

	el["edit_list_table"] = entries
	fmt.Println(entries)

	a.elements = el

	return a
}

// general atoms with children
func parse_general_with_children(a *Atom, r io.Reader) *Atom {

	a.children = Parse_atom(r, a.level+1)

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

// parse each atoms
func Parse_atom(r io.Reader, level int) []Atom {

	var atoms = make([]Atom, 0)
	var size uint32

	buf := make([]byte, 4)

	for binary.Read(r, binary.BigEndian, &size) == nil {

		atom := new(Atom)
		atom.size = size
		atom.level = level

		if level == 0 {
			atom.parent = nil
		}

		r.Read(buf)
		atom.atype = string(buf)

		if atom_parsers[atom.atype] != nil {
			atom_parsers[atom.atype](atom, r)
		} else {
			parse_general(atom, r)
		}

		atoms = append(atoms, *atom)
	}

	return atoms
}

// initialize
func init() {

	// initialize parser table
	atom_parsers = map[string]func(*Atom, io.Reader) *Atom{
		"moov": parse_general_with_children,
		"trak": parse_general_with_children,
		"edts": parse_general_with_children,
		"ftyp": parse_ftyp,
		"mvhd": parse_mvhd,
		"free": parse_free,
		"tkhd": parse_tkhd,
		"elst": parse_elst,
	}

	diff_time = time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC).Sub(time.Unix(0, 0))
}
