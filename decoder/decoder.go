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
	flushAdcBuf  []uint16
	currentEvent *entities.Event
	events       []*entities.Event
}

func NewDecoder(reader io.Reader, endian binary.ByteOrder) *Decoder {
	return &Decoder{
		reader:       reader,
		endian:       endian,
		flushAdcBuf:  make([]uint16, 4*entconst.Clock),
		currentEvent: &entities.Event{},
		events:       make([]*entities.Event, 0, 1000),
	}
}

// Decode entrypoint
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

	if err := d.ReadFlushAdc(); err != nil {
		return err
	}

	if err := d.ReadVersionAndDepth(); err != nil {
		return err
	}

	return nil
}

func (d *Decoder) SkipEventHeaderSymbol() error {
	var b [4]byte
	if err := binary.Read(d.reader, d.endian, &b); err != nil {
		return err
	}
	if !entconst.IsValidHeaderSymbol(b) {
		return entities.InvalidHeaderError{Got: b}
	}

	return nil
}

func (d *Decoder) ReadEventCounter() error {
	if err := binary.Read(d.reader, d.endian, &d.currentEvent.Header.Counter); err != nil {
		return err
	}

	return nil
}

func (d *Decoder) ReadClockCounter() error {
	if err := binary.Read(d.reader, d.endian, &d.currentEvent.Header.Clock); err != nil {
		return err
	}

	return nil
}

// ReadFlushAdc 4ch * 1024 clockを一気読み
func (d *Decoder) ReadFlushAdc() error {
	if err := binary.Read(d.reader, d.endian, &d.flushAdcBuf); err != nil {
		return err
	}

	// 1 clock 2byte*4 = uint16 * 4
	for clock := 0; clock < entconst.Clock; clock++ {
		flushADC4ChBuff := d.flushAdcBuf[4*clock : 4*clock+4]
		for ch := 0; ch < 4; ch++ {
			header := flushADC4ChBuff[ch] >> 12                  // 上位4bit
			adcValue := flushADC4ChBuff[ch] & 0b0000001111111111 // 下位10bit
			if !entconst.IsValidAdcHeaderSymbol(ch, header) {
				return entities.InvalidFlushAdcHeaderError{
					Got:      header,
					Expected: entconst.FlushAdcHeaderSymbol()[ch],
				}
			}
			d.currentEvent.Header.FlushADC[ch][clock] = adcValue
		}
	}

	return nil
}
func (d *Decoder) ReadVersionAndDepth() error {
	if err := binary.Read(d.reader, d.endian, &d.currentEvent.Header.Version.Year); err != nil {
		return err
	}

	if err := binary.Read(d.reader, d.endian, &d.currentEvent.Header.Version.Month); err != nil {
		return err
	}

	var buf uint16
	if err := binary.Read(d.reader, d.endian, &buf); err != nil {
		return err
	}
	d.currentEvent.Header.Version.Sub = uint8(buf >> 12)                                             // 上位4bit
	d.currentEvent.Header.EncodingClockDepth = entities.EncodingClockDepth(buf & 0b0000011111111111) // 下位11bit

	return nil
}
