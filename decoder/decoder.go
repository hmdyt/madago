package decoder

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"

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
	// headerまで読み飛ばす
	for i := 0; ; i++ {
		b, err := d.reader.Peek(4)
		if err != nil {
			return []*entities.Event{}, err
		}

		if entconst.IsEventHeaderSymbol(b) {
			break
		} else {
			if _, err := d.reader.Discard(1); err != nil {
				return []*entities.Event{}, err
			}
		}
	}

	// event read loop
	for {
		if err := d.DecodeEvent(); err == nil {
			d.events = append(d.events, d.currentEvent)
			d.clearCurrentEvent()
		} else if errors.Is(err, io.ErrUnexpectedEOF) {
			break
		} else if errors.Is(err, io.EOF) {
			break
		} else {
			log.Fatalf("event read loop end: %#v \n", err)
		}
	}

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

	if err := d.ReadInputCh2Counter(); err != nil {
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
	if !entconst.IsEventHeaderSymbol(b[:]) {
		return entities.InvalidHeaderError{Got: b[:]}
	}

	return nil
}

func (d *Decoder) ReadEventCounter() error {
	if err := binary.Read(d.reader, d.endian, &d.currentEvent.Header.Trigger); err != nil {
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

func (d *Decoder) ReadInputCh2Counter() error {
	if err := binary.Read(d.reader, d.endian, &d.currentEvent.Header.InputCh2); err != nil {
		return err
	}

	return nil
}

// ReadFlushAdc 4ch * 1024 clock
func (d *Decoder) ReadFlushAdc() error {
	var clock_counter [4]int
	for {
		ch, err := d.peekFlushAdcHeader()
		if err != nil {
			return err
		} else if ch == nil {
			return nil
		} else {
			var buf uint16
			if err := binary.Read(d.reader, d.endian, &buf); err != nil {
				return err
			}
			adcValue := buf & 0b0000001111111111 // 下位10bit
			d.currentEvent.Header.FlushADC[*ch][clock_counter[*ch]] = adcValue
			clock_counter[*ch]++
		}
	}
}

// peekFlushAdcHeader FADCのheaderをみて, 対応するchannelを返す。対応がない場合はnil
func (d *Decoder) peekFlushAdcHeader() (*int, error) {
	buf, err := d.reader.Peek(1)
	if err != nil {
		return nil, err
	}

	header := buf[0] >> 4
	for ch := 0; ch < 4; ch++ {
		if entconst.IsAdcHeaderSymbol(ch, uint16(header)) {
			return &ch, nil
		}
	}

	return nil, nil
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
