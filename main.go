package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hmdyt/madago/decoder"
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

	fmt.Println(len(events))
}
