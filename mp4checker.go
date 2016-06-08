package main

import (
	"flag"
	"github.com/rarewin/mp4checker/atom"
	"os"
)

type Options struct {
	help *bool
}

var options Options

func init() {
	options.help = flag.Bool("h", false, "Print this help")
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

	atom.Parse_atom(file_i)

}
