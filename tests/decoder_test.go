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
		input     []byte
		wantEvent []*entities.Event
	}{
		{
			input: []byte{
				0xeb, 0x90, 0x19, 0x64,
				0x00, 0x00, 0x00, 0x0a,
				0x00, 0x00, 0x00, 0x0b,
			},
			wantEvent: []*entities.Event{
				{
					Header: entities.EventHeader{
						Counter: 10,
						Clock:   11,
					},
				},
			},
		},
	}

	for i, tt := range tests {
		reader := bytes.NewBuffer(tt.input)
		d := decoder.NewDecoder(reader, binary.BigEndian)

		events, err := d.Decode()
		if err != nil {
			t.Errorf("test%d failed DecodeEvent: %s", i, err.Error())
		}

		if diff := cmp.Diff(events, tt.wantEvent); diff != "" {
			t.Errorf("return Event is mismatch :\n%s", diff)
		}
	}
}
