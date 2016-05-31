package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
)

type Options struct {
	help *bool
}

var options Options

func init() {
	options.help = flag.Bool("h", false, "Print this help")
}

type Atom struct {
	size  uint32
	atype string
}

func print_atom(a Atom) {

	fmt.Printf("size: %d\n", a.size)
	fmt.Printf("type: %s\n", a.atype)
}

func parse_atom(r io.Reader) {

	var atom Atom

	buf := make([]byte, 4)

	binary.Read(r, binary.BigEndian, &atom.size)

	r.Read(buf)
	atom.atype = string(buf)

	print_atom(atom)
}

func main() {

	// paser command line options
	flag.Parse()

	if *options.help {
		flag.Usage()
		os.Exit(1)
	}

	var args = flag.Args()

	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	var inputfile = args[0]

	file_i, err := os.Open(inputfile)

	if err != nil {
		os.Exit(1)
	}

	parse_atom(file_i)

}
