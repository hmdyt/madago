package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"log"
	"os"

	"github.com/hmdyt/madago/bar"
	"github.com/hmdyt/madago/decoder"
	"github.com/hmdyt/madago/encoder/root"
	"go-hep.org/x/hep/groot"
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
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("file close %s: %s", path, err.Error())
		}
	}()

	// decode
	decoderProgressBar := bar.DecoderProgressBar(fileInfo.Size())
	reader := bufio.NewReader(decoderProgressBar.NewProxyReader(file))
	d := decoder.NewDecoder(reader, binary.BigEndian)
	events, err := d.Decode()
	if err != nil {
		log.Fatalln(err)
	}
	decoderProgressBar.Finish()

	// open file
	f, err := groot.Create("tree2.root")
	if err != nil {
		log.Fatalf("groot.Create : %s", err.Error())
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("f.Close : %s", err.Error())
		}
	}()

	// encode
	encoderProgressBar := bar.EncoderProgressBar(len(events))
	rootEncoder, err := root.NewEncoder(f, encoderProgressBar)
	if err != nil {
		log.Fatalf("NewEncoder : %s", err.Error())
	}

	if err := rootEncoder.Write(events); err != nil {
		log.Fatalf("RootEncoder.Write : %s", err.Error())
	}
	encoderProgressBar.Finish()
}
