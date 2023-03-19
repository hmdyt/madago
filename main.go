package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"log"
	"os"

	"github.com/hmdyt/madago/decoder"
	"github.com/hmdyt/madago/encoder/root"
)

func main() {
	flag.Parse()
	path := flag.Arg(0)

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("file open %s: %s", path, err.Error())
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	d := decoder.NewDecoder(reader, binary.BigEndian)

	events, err := d.Decode()
	if err != nil {
		log.Fatalln(err)
	}

	rootEncoder, err := root.NewEncoder("tree.root")
	if err != nil {
		log.Fatalf("NewEncoder : %s", err.Error())
	}

	if err := rootEncoder.Write(events); err != nil {
		log.Fatalf("RootEncoder.Write : %s", err.Error())
	}

}
