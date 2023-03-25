package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
	"github.com/hmdyt/madago/decoder"
	"github.com/hmdyt/madago/encoder/root"
)

func main() {
	flag.Parse()
	path := flag.Arg(0)

	// open file
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("file open %s: %s", path, err.Error())
	}
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("file stat %s: %s", path, err.Error())
	}
	defer file.Close()

	// decode
	bar := pb.Full.Start64(fileInfo.Size())
	reader := bufio.NewReader(bar.NewProxyReader(file))
	d := decoder.NewDecoder(reader, binary.BigEndian)
	events, err := d.Decode()
	if err != nil {
		log.Fatalln(err)
	}
	bar.Finish()

	rootEncoder, err := root.NewEncoder("tree.root")
	if err != nil {
		log.Fatalf("NewEncoder : %s", err.Error())
	}

	if err := rootEncoder.Write(events); err != nil {
		log.Fatalf("RootEncoder.Write : %s", err.Error())
	}

}
