package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func _main() int {
	var offset int64
	var length int64
	flag.Int64Var(&offset, "o", 0, "Offset bytes")
	flag.Int64Var(&length, "n", 1024, "Maximum byte length to dump")
	flag.Parse()

	args := flag.Args()

	var r io.Reader
	if len(args) == 0 {
		r = os.Stdin
	} else {
		name := args[0]
		if name == "-" {
			r = os.Stdin
		} else {
			file, err := os.Open(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to open file: %s\n", name)
				return 1
			}
			defer file.Close()

			r = file
		}
	}

	if offset > 0 {
		_, err := io.CopyN(ioutil.Discard, r, offset)
		if err != nil {
			if err == io.EOF {
				return 0
			}

			fmt.Fprintf(os.Stderr, "Failed to read: %v\n", err)
			return 1
		}
	}

	d := hex.Dumper(os.Stdout)
	defer d.Close()

	if _, err := io.CopyN(d, r, length); err != nil {
		if err == io.EOF {
			return 0
		}

		fmt.Fprintf(os.Stderr, "Failed to read: %v\n", err)
		return 1
	}

	return 0
}

func main() {
	code := _main()
	if code != 0 {
		os.Exit(code)
	}
}
