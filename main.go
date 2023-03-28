package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"log"
	"os"

	"github.com/hmdyt/madago/bar"
	"github.com/hmdyt/madago/decoder"
	"github.com/hmdyt/madago/domain/entities"
	"github.com/hmdyt/madago/encoder/root"
	"github.com/hmdyt/madago/usecases"
	"go-hep.org/x/hep/groot"
)

func main() {
	flag.Parse()
	path00 := flag.Arg(0)
	path01 := flag.Arg(1)
	path03 := flag.Arg(2)
	path10 := flag.Arg(3)
	path11 := flag.Arg(4)
	path13 := flag.Arg(5)
	outputPath := flag.Arg(6)
	filePaths := map[entities.BoardID]string{
		entities.GBKB00: path00,
		entities.GBKB01: path01,
		entities.GBKB03: path03,
		entities.GBKB10: path10,
		entities.GBKB11: path11,
		entities.GBKB13: path13,
	}

	// decode
	madaEvents := make(map[entities.BoardID][]*entities.MadaEvent, len(filePaths))
	for boardID, path := range filePaths {
		file, err := os.Open(path)
		if err != nil {
			log.Fatalf("file open %s: %s", path, err.Error())
		}
		fileInfo, err := file.Stat()
		if err != nil {
			log.Fatalf("file stat %s: %s", path, err.Error())
		}

		b := bar.DecoderProgressBar(fileInfo.Size())
		r := bufio.NewReader(b.NewProxyReader(file))
		d := decoder.NewDecoder(r, binary.BigEndian)
		madaEvents[boardID], err = d.Decode()
		if err != nil {
			log.Fatalln(err)
		}

		b.Finish()
		if err := file.Close(); err != nil {
			log.Fatalf("file close at decoder %s: %s", path, err.Error())
		}
	}

	// merge
	// TODO: implement progress bar
	cmd := usecases.MadaMergeCmd{
		MadaEventMap: madaEvents,
	}
	rawEvents := usecases.MergeMadaEvents(cmd)

	// encode
	f, err := groot.Create(outputPath)
	if err != nil {
		log.Fatalf("groot.Create : %s", err.Error())
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("f.Close : %s", err.Error())
		}
	}()
	encoderProgressBar := bar.EncoderProgressBar(len(rawEvents))
	rootEncoder := root.NewRawEncoder(f, encoderProgressBar)
	if err := rootEncoder.Encode(rawEvents); err != nil {
		log.Fatalf("RootEncoder.Write : %s", err.Error())
	}
	encoderProgressBar.Finish()
}
