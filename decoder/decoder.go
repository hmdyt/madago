package decoder

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/hmdyt/madago/domain/entconst"
	"github.com/hmdyt/madago/domain/entities"
)

type Decoder struct {
	reader       *bufio.Reader
	endian       binary.ByteOrder
	flushAdcBuf  []uint16
	currentEvent *entities.Event
	events       []*entities.Event
}

func NewDecoder(reader *bufio.Reader, endian binary.ByteOrder) *Decoder {
	return &Decoder{
		reader:       reader,
		endian:       endian,
		flushAdcBuf:  make([]uint16, 4*entconst.Clock),
		currentEvent: &entities.Event{Hits: make([]entities.Hit, 0, entconst.Clock)},
		events:       make([]*entities.Event, 0, 1000),
	}
}

// Decode entrypoint
func (d *Decoder) Decode() ([]*entities.Event, error) {
	if err := d.DecodeEvent(); err != nil {
		return []*entities.Event{}, err
	}

	d.events = append(d.events, d.currentEvent)
	d.clearCurrentEvent()

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

	// footer ("uPIC") に当たるまでhitを読み続ける
	for {
		peekBytes, err := d.reader.Peek(4)
		if err != nil {
			return err
		}
		switch {
		case entconst.IsHitHeaderSymbol(peekBytes[0] >> 4):
			if _, err := d.reader.Discard(2); err != nil {
				return err
			}
			if err := d.ReadHit(); err != nil {
				return err
			}
		case entconst.IsEventFooterSymbol(peekBytes):
			if _, err := d.reader.Discard(4); err != nil {
				return err
			}
			return nil
		default:
			return entities.InvalidHeaderError{Got: peekBytes}
		}
	}
}

func (d *Decoder) SkipEventHeaderSymbol() error {
	var b [4]byte
	if err := binary.Read(d.reader, d.endian, &b); err != nil {
		return err
	}
	if !entconst.IsEventHeaderSymbol(b) {
		return entities.InvalidHeaderError{Got: b[:]}
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
			if !entconst.IsAdcHeaderSymbol(ch, header) {
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

func (d *Decoder) ReadHit() error {
	var hit entities.Hit

	// clock
	if err := binary.Read(d.reader, d.endian, &hit.Clock); err != nil {
		return err
	}

	// hit: 64bitずつよむ
	var buf uint64

	// 127 - 64 ch
	if err := binary.Read(d.reader, d.endian, &buf); err != nil {
		return err
	}
	for i := 0; i < 64; i++ {
		// 下から i bit 目をboolとして取り出す
		isHitInt := (buf >> i) & 0b1
		switch isHitInt {
		case 0:
			hit.IsHit[64+i] = false
		case 1:
			hit.IsHit[64+i] = true
		default:
			return errors.New(fmt.Sprintf("invalid IsHit, got=%d", isHitInt))
		}
	}

	// 63 - 0 ch
	if err := binary.Read(d.reader, d.endian, &buf); err != nil {
		return err
	}
	for i := 0; i < 64; i++ {
		// 下から i bit 目をboolとして取り出す
		isHitInt := (buf >> i) & 0b1
		switch isHitInt {
		case 0:
			hit.IsHit[i] = false
		case 1:
			hit.IsHit[i] = true
		default:
			return errors.New(fmt.Sprintf("invalid IsHit, got=%d", isHitInt))
		}
	}

	d.currentEvent.Hits = append(d.currentEvent.Hits, hit)
	return nil
}

func (d *Decoder) clearCurrentEvent() {
	d.currentEvent = &entities.Event{Hits: make([]entities.Hit, 0, entconst.Clock)}
}
