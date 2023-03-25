package tests

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hmdyt/madago/decoder"
	"github.com/hmdyt/madago/domain/entities"
)

func TestDecodeEvents(t *testing.T) {
	wantHits := []entities.Hit{
		{
			Clock: 50,
			IsHit: func() [128]bool {
				var ret [128]bool
				for i := range ret {
					ret[i] = true
				}
				return ret
			}(),
		},
		{
			Clock: 51,
			IsHit: func() [128]bool {
				var ret [128]bool
				for i := range ret {
					ret[i] = true
				}
				return ret
			}(),
		},
		{
			Clock: 52,
			IsHit: func() [128]bool {
				var ret [128]bool
				for i := range ret {
					ret[i] = true
				}
				return ret
			}(),
		},
	}

	// テストケース
	tests := []struct {
		name      string
		inputs    [][]byte
		wantEvent []*entities.Event
	}{
		{
			name: "成功; FlushADCあり, hit3clock (2events)",
			inputs: [][]byte{
				{0x00, 0x00, 0x00, 0x00}, // trash data
				makeTestEventHeader(),
				makeTestCounters(),
				makeTestFADC(),
				makeTestVersionAndDepth(),
				makeTestHit(),
				makeTestEventFooter(),
				makeTestEventHeader(),
				makeTestCounters(),
				makeTestFADC(),
				makeTestVersionAndDepth(),
				makeTestHit(),
				makeTestEventFooter(),
			},
			wantEvent: []*entities.Event{
				{
					Header: entities.EventHeader{
						Trigger:  10,
						Clock:    11,
						InputCh2: 12,
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
					Hits: wantHits,
				},
				{
					Header: entities.EventHeader{
						Trigger:  10,
						Clock:    11,
						InputCh2: 12,
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
					Hits: wantHits,
				},
			},
		},
		{
			name: "成功; FlushADC ch2欠け, hit3clock (2events)",
			inputs: [][]byte{
				makeTestEventHeader(),
				makeTestCounters(),
				makeTestLuckFADC(),
				makeTestVersionAndDepth(),
				makeTestHit(),
				makeTestEventFooter(),
				makeTestEventHeader(),
				makeTestCounters(),
				makeTestLuckFADC(),
				makeTestVersionAndDepth(),
				makeTestHit(),
				makeTestEventFooter(),
			},
			wantEvent: []*entities.Event{
				{
					Header: entities.EventHeader{
						Trigger:  10,
						Clock:    11,
						InputCh2: 12,
						FlushADC: func() [4][1024]uint16 {
							var ret [4][1024]uint16
							for clock := 0; clock < 1024; clock++ {
								for ch := 0; ch < 4; ch++ {
									if ch == 2 {
										ret[ch][clock] = 0
									} else {
										ret[ch][clock] = uint16(ch + 13)
									}
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
					Hits: wantHits,
				},
				{
					Header: entities.EventHeader{
						Trigger:  10,
						Clock:    11,
						InputCh2: 12,
						FlushADC: func() [4][1024]uint16 {
							var ret [4][1024]uint16
							for clock := 0; clock < 1024; clock++ {
								for ch := 0; ch < 4; ch++ {
									if ch == 2 {
										ret[ch][clock] = 0
									} else {
										ret[ch][clock] = uint16(ch + 13)
									}
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
					Hits: wantHits,
				},
			},
		},
		{
			name: "成功; FlushADC 全chなし, hit3clock (2events)",
			inputs: [][]byte{
				makeTestEventHeader(),
				makeTestCounters(),
				makeTestVersionAndDepth(),
				makeTestHit(),
				makeTestEventFooter(),
				makeTestEventHeader(),
				makeTestCounters(),
				makeTestVersionAndDepth(),
				makeTestHit(),
				makeTestEventFooter(),
			},
			wantEvent: []*entities.Event{
				{
					Header: entities.EventHeader{
						Trigger:  10,
						Clock:    11,
						InputCh2: 12,
						Version: entities.Version{
							Year:  23,
							Month: 3,
							Sub:   5,
						},
						EncodingClockDepth: entities.EncodingClockDepth(50),
					},
					Hits: wantHits,
				},
				{
					Header: entities.EventHeader{
						Trigger:  10,
						Clock:    11,
						InputCh2: 12,
						Version: entities.Version{
							Year:  23,
							Month: 3,
							Sub:   5,
						},
						EncodingClockDepth: entities.EncodingClockDepth(50),
					},
					Hits: wantHits,
				},
			},
		},
		{
			name: "成功; FlushADCとhitなし (2events)",
			inputs: [][]byte{
				makeTestEventHeader(),
				makeTestCounters(),
				makeTestVersionAndDepth(),
				makeTestEventFooter(),
				makeTestEventHeader(),
				makeTestCounters(),
				makeTestVersionAndDepth(),
				makeTestEventFooter(),
			},
			wantEvent: []*entities.Event{
				{
					Header: entities.EventHeader{
						Trigger:  10,
						Clock:    11,
						InputCh2: 12,
						Version: entities.Version{
							Year:  23,
							Month: 3,
							Sub:   5,
						},
						EncodingClockDepth: entities.EncodingClockDepth(50),
					},
					Hits: []entities.Hit{},
				},
				{
					Header: entities.EventHeader{
						Trigger:  10,
						Clock:    11,
						InputCh2: 12,
						Version: entities.Version{
							Year:  23,
							Month: 3,
							Sub:   5,
						},
						EncodingClockDepth: entities.EncodingClockDepth(50),
					},
					Hits: []entities.Hit{},
				},
			},
		},
	}

	for i, tt := range tests {
		// flatten inputs
		inputBuf := bytes.NewBuffer([]byte{})
		for _, b := range tt.inputs {
			inputBuf.Write(b)
		}
		reader := bufio.NewReader(inputBuf)
		d := decoder.NewDecoder(reader, binary.BigEndian)

		events, err := d.Decode()
		if err != nil {
			t.Fatalf("fail %s: test%d failed DecodeEvent: %s", tt.name, i, err.Error())
		}

		if diff := cmp.Diff(events, tt.wantEvent); diff != "" {
			t.Fatalf("fail %s:return Event is mismatch :\n%s", tt.name, diff)
		}
		t.Logf("pass: %s", tt.name)
	}
}
