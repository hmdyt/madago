package decoder

import (
	"encoding/binary"
	"io"

	"github.com/hmdyt/madago/domain/entconst"
	"github.com/hmdyt/madago/domain/entities"
)

type Decoder struct {
	reader       io.Reader
	endian       binary.ByteOrder
	buf          []byte
	currentEvent *entities.Event
	events       []*entities.Event
}

func NewDecoder(reader io.Reader, endian binary.ByteOrder) *Decoder {
	return &Decoder{
		reader:       reader,
		endian:       endian,
		buf:          make([]byte, 128),
		currentEvent: &entities.Event{},
		events:       make([]*entities.Event, 0, 1000),
	}
}

func (d *Decoder) Decode() ([]*entities.Event, error) {
	if err := d.DecodeEvent(); err != nil {
		return []*entities.Event{}, err
	}

	d.events = append(d.events, d.currentEvent)

	return d.events, nil
}

func (d *Decoder) DecodeEvent() error {
	if err := d.SkipEventHeaderSymbol(); err != nil {
		return err
	}

	if err := d.ReadEventCounter(); err != nil {
		return err
	}

	if err := d.ReadClockCounter(); err != nil {
		return err
	}

	return nil
}

func (d *Decoder) SkipEventHeaderSymbol() error {
	var b [4]byte
	err := binary.Read(d.reader, d.endian, &b)
	if err != nil {
		return err
	}
	if !entconst.IsValidHeaderSymbol(b) {
		return entities.InvalidHeaderError{Got: b}
	}

	return nil
}

func (d *Decoder) ReadEventCounter() error {
	err := binary.Read(d.reader, d.endian, &d.currentEvent.Header.Counter)
	if err != nil {
		return err
	}

	return nil
}

func (d *Decoder) ReadClockCounter() error {
	err := binary.Read(d.reader, d.endian, &d.currentEvent.Header.Clock)
	if err != nil {
		return err
	}

	return nil
}
