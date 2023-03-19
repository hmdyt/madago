package tests

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hmdyt/madago/decoder"
	"github.com/hmdyt/madago/domain/entities"
)

func TestDecodeEvent(t *testing.T) {
	tests := []struct {
		inputs    [][]byte
		wantEvent []*entities.Event
	}{
		{
			inputs: [][]byte{
				{
					0xeb, 0x90, 0x19, 0x64, // header
					0x00, 0x00, 0x00, 0x0a, // event counter
					0x00, 0x00, 0x00, 0x0b, // clock counter
				},
				// flush adc: 全clockで同じ値, ch0 -> 13, ch1 -> 14, ch2 -> 15, ch3 -> 16
				func() []byte {
					ret := make([]byte, 0, 2*4*1024)
					for clock := 0; clock < 1024; clock++ {
						for ch := 0; ch < 4; ch++ {
							header := uint8(ch+4) << 4 // fixed value
							adcValue := uint8(ch + 13) // test random value
							ret = append(ret, header, adcValue)
						}
					}
					return ret
				}(),
				func() []byte {
					ret := make([]byte, 0, 4)
					versionYear := uint8(23)
					versionMonth := uint8(3)
					versionSub := uint8(5) << 4
					encodingClockDepth := uint8(50)
					ret = append(ret, versionYear, versionMonth, versionSub, encodingClockDepth)
					return ret
				}(),
			},
			wantEvent: []*entities.Event{
				{
					Header: entities.EventHeader{
						Counter: 10,
						Clock:   11,
						FlushADC: func() [4][1024]uint16 {
							var ret [4][1024]uint16
							for clock := 0; clock < 1024; clock++ {
								for ch := 0; ch < 4; ch++ {
									ret[ch][clock] = uint16(ch + 13)
								}
							}
							return ret
						}(),
						Version: entities.Version{
							Year:  23,
							Month: 3,
							Sub:   5,
						},
						EncodingClockDepth: entities.EncodingClockDepth(50),
					},
				},
			},
		},
	}

	for i, tt := range tests {
		// flatten inputs
		reader := bytes.NewBuffer([]byte{})
		for _, b := range tt.inputs {
			reader.Write(b)
		}
		d := decoder.NewDecoder(reader, binary.BigEndian)

		events, err := d.Decode()
		if err != nil {
			t.Fatalf("test%d failed DecodeEvent: %s", i, err.Error())
		}

		if diff := cmp.Diff(events[0].Header.FlushADC, tt.wantEvent[0].Header.FlushADC); diff != "" {
			t.Fatalf("return Event is mismatch :\n%s", diff)
		}

		if diff := cmp.Diff(events, tt.wantEvent); diff != "" {
			t.Fatalf("return Event is mismatch :\n%s", diff)
		}
	}
}
