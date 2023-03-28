package tests

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hmdyt/madago/domain/entities"
	"github.com/hmdyt/madago/usecases"
)

func TestMadaMerge(t *testing.T) {

	wantHits := func(id int) entities.RawHits {
		var ret entities.RawHits
		for i := 0; i < 128; i++ {
			ret[id][i] = true
		}
		return ret
	}

	wantVersion := func(id int) entities.Version {
		return entities.Version{
			Year:  uint8(id),
			Month: uint8(id),
			Sub:   uint8(id),
		}
	}

	tests := []struct {
		name  string
		input usecases.MadaMergeCmd
		want  []*entities.RawEvent
	}{
		{
			name: "成功; 1event, 1board",
			input: usecases.MadaMergeCmd{
				MadaEventMap: map[entities.BoardID][]*entities.MadaEvent{
					entities.GBKB00: []*entities.MadaEvent{
						makeTestMadaEvent(0, 0),
					},
				},
			},
			want: []*entities.RawEvent{
				{
					Trigger:     0,
					ClockMap:    map[entities.BoardID]entities.ClockCounter{entities.GBKB00: 0},
					InputCh2Map: map[entities.BoardID]entities.InputCh2Counter{entities.GBKB00: 0},
					VersionMap: map[entities.BoardID]entities.Version{
						entities.GBKB00: wantVersion(0),
					},
					EncodingClockDepthMap: map[entities.BoardID]entities.EncodingClockDepth{entities.GBKB00: 0},
					FlushAdcMap: map[entities.BoardID]entities.FlushAdc{
						entities.GBKB00: makeTestMadaEvent(0, 0).FlushADC,
					},
					HitsMap: map[entities.BoardID]entities.RawHits{
						entities.GBKB00: wantHits(0),
					},
				},
			},
		},
		{
			name: "成功; 1event, 2board",
			input: usecases.MadaMergeCmd{
				MadaEventMap: map[entities.BoardID][]*entities.MadaEvent{
					entities.GBKB00: []*entities.MadaEvent{
						makeTestMadaEvent(0, 0),
					},
					entities.GBKB01: []*entities.MadaEvent{
						makeTestMadaEvent(0, 1),
					},
				},
			},
			want: []*entities.RawEvent{
				{
					Trigger: 0,
					ClockMap: map[entities.BoardID]entities.ClockCounter{
						entities.GBKB00: 0,
						entities.GBKB01: 1,
					},
					InputCh2Map: map[entities.BoardID]entities.InputCh2Counter{
						entities.GBKB00: 0,
						entities.GBKB01: 1,
					},
					VersionMap: map[entities.BoardID]entities.Version{
						entities.GBKB00: wantVersion(0),
						entities.GBKB01: wantVersion(1),
					},
					EncodingClockDepthMap: map[entities.BoardID]entities.EncodingClockDepth{
						entities.GBKB00: 0,
						entities.GBKB01: 1,
					},
					FlushAdcMap: map[entities.BoardID]entities.FlushAdc{
						entities.GBKB00: makeTestMadaEvent(0, 0).FlushADC,
						entities.GBKB01: makeTestMadaEvent(0, 1).FlushADC,
					},
					HitsMap: map[entities.BoardID]entities.RawHits{
						entities.GBKB00: wantHits(0),
						entities.GBKB01: wantHits(1),
					},
				},
			},
		},
		{
			name: "成功; 2event, 2board",
			input: usecases.MadaMergeCmd{
				MadaEventMap: map[entities.BoardID][]*entities.MadaEvent{
					entities.GBKB00: []*entities.MadaEvent{
						makeTestMadaEvent(0, 0),
						makeTestMadaEvent(1, 1),
					},
					entities.GBKB01: []*entities.MadaEvent{
						makeTestMadaEvent(0, 2),
						makeTestMadaEvent(1, 3),
					},
				},
			},
			want: []*entities.RawEvent{
				{
					Trigger: 0,
					ClockMap: map[entities.BoardID]entities.ClockCounter{
						entities.GBKB00: 0,
						entities.GBKB01: 2,
					},
					InputCh2Map: map[entities.BoardID]entities.InputCh2Counter{
						entities.GBKB00: 0,
						entities.GBKB01: 2,
					},
					VersionMap: map[entities.BoardID]entities.Version{
						entities.GBKB00: wantVersion(0),
						entities.GBKB01: wantVersion(2),
					},
					EncodingClockDepthMap: map[entities.BoardID]entities.EncodingClockDepth{
						entities.GBKB00: 0,
						entities.GBKB01: 2,
					},
					FlushAdcMap: map[entities.BoardID]entities.FlushAdc{
						entities.GBKB00: makeTestMadaEvent(0, 0).FlushADC,
						entities.GBKB01: makeTestMadaEvent(0, 2).FlushADC,
					},
					HitsMap: map[entities.BoardID]entities.RawHits{
						entities.GBKB00: wantHits(0),
						entities.GBKB01: wantHits(2),
					},
				},
				{
					Trigger: 1,
					ClockMap: map[entities.BoardID]entities.ClockCounter{
						entities.GBKB00: 1,
						entities.GBKB01: 3,
					},
					InputCh2Map: map[entities.BoardID]entities.InputCh2Counter{
						entities.GBKB00: 1,
						entities.GBKB01: 3,
					},
					VersionMap: map[entities.BoardID]entities.Version{
						entities.GBKB00: wantVersion(1),
						entities.GBKB01: wantVersion(3),
					},
					EncodingClockDepthMap: map[entities.BoardID]entities.EncodingClockDepth{
						entities.GBKB00: 1,
						entities.GBKB01: 3,
					},
					FlushAdcMap: map[entities.BoardID]entities.FlushAdc{
						entities.GBKB00: makeTestMadaEvent(1, 1).FlushADC,
						entities.GBKB01: makeTestMadaEvent(1, 3).FlushADC,
					},
					HitsMap: map[entities.BoardID]entities.RawHits{
						entities.GBKB00: wantHits(1),
						entities.GBKB01: wantHits(3),
					},
				},
			},
		},
		{
			name: "成功; 3event, 6board",
			input: usecases.MadaMergeCmd{
				MadaEventMap: map[entities.BoardID][]*entities.MadaEvent{
					entities.GBKB00: []*entities.MadaEvent{
						makeTestMadaEvent(0, 0),
						makeTestMadaEvent(1, 1),
						makeTestMadaEvent(2, 2),
					},
					entities.GBKB01: []*entities.MadaEvent{
						makeTestMadaEvent(0, 3),
						makeTestMadaEvent(1, 4),
						makeTestMadaEvent(2, 5),
					},
					entities.GBKB03: []*entities.MadaEvent{
						makeTestMadaEvent(0, 6),
						makeTestMadaEvent(1, 7),
						makeTestMadaEvent(2, 8),
					},
					entities.GBKB10: []*entities.MadaEvent{
						makeTestMadaEvent(0, 9),
						makeTestMadaEvent(1, 10),
						makeTestMadaEvent(2, 11),
					},
					entities.GBKB11: []*entities.MadaEvent{
						makeTestMadaEvent(0, 12),
						makeTestMadaEvent(1, 13),
						makeTestMadaEvent(2, 14),
					},
					entities.GBKB13: []*entities.MadaEvent{
						makeTestMadaEvent(0, 15),
						makeTestMadaEvent(1, 16),
						makeTestMadaEvent(2, 17),
					},
				},
			},
			want: []*entities.RawEvent{
				{
					Trigger: 0,
					ClockMap: map[entities.BoardID]entities.ClockCounter{
						entities.GBKB00: 0,
						entities.GBKB01: 3,
						entities.GBKB03: 6,
						entities.GBKB10: 9,
						entities.GBKB11: 12,
						entities.GBKB13: 15,
					},
					InputCh2Map: map[entities.BoardID]entities.InputCh2Counter{
						entities.GBKB00: 0,
						entities.GBKB01: 3,
						entities.GBKB03: 6,
						entities.GBKB10: 9,
						entities.GBKB11: 12,
						entities.GBKB13: 15,
					},
					VersionMap: map[entities.BoardID]entities.Version{
						entities.GBKB00: wantVersion(0),
						entities.GBKB01: wantVersion(3),
						entities.GBKB03: wantVersion(6),
						entities.GBKB10: wantVersion(9),
						entities.GBKB11: wantVersion(12),
						entities.GBKB13: wantVersion(15),
					},
					EncodingClockDepthMap: map[entities.BoardID]entities.EncodingClockDepth{
						entities.GBKB00: 0,
						entities.GBKB01: 3,
						entities.GBKB03: 6,
						entities.GBKB10: 9,
						entities.GBKB11: 12,
						entities.GBKB13: 15,
					},
					FlushAdcMap: map[entities.BoardID]entities.FlushAdc{
						entities.GBKB00: makeTestMadaEvent(0, 0).FlushADC,
						entities.GBKB01: makeTestMadaEvent(0, 3).FlushADC,
						entities.GBKB03: makeTestMadaEvent(0, 6).FlushADC,
						entities.GBKB10: makeTestMadaEvent(0, 9).FlushADC,
						entities.GBKB11: makeTestMadaEvent(0, 12).FlushADC,
						entities.GBKB13: makeTestMadaEvent(0, 15).FlushADC,
					},
					HitsMap: map[entities.BoardID]entities.RawHits{
						entities.GBKB00: wantHits(0),
						entities.GBKB01: wantHits(3),
						entities.GBKB03: wantHits(6),
						entities.GBKB10: wantHits(9),
						entities.GBKB11: wantHits(12),
						entities.GBKB13: wantHits(15),
					},
				},
				{
					Trigger: 1,
					ClockMap: map[entities.BoardID]entities.ClockCounter{
						entities.GBKB00: 1,
						entities.GBKB01: 4,
						entities.GBKB03: 7,
						entities.GBKB10: 10,
						entities.GBKB11: 13,
						entities.GBKB13: 16,
					},
					InputCh2Map: map[entities.BoardID]entities.InputCh2Counter{
						entities.GBKB00: 1,
						entities.GBKB01: 4,
						entities.GBKB03: 7,
						entities.GBKB10: 10,
						entities.GBKB11: 13,
						entities.GBKB13: 16,
					},
					VersionMap: map[entities.BoardID]entities.Version{
						entities.GBKB00: wantVersion(1),
						entities.GBKB01: wantVersion(4),
						entities.GBKB03: wantVersion(7),
						entities.GBKB10: wantVersion(10),
						entities.GBKB11: wantVersion(13),
						entities.GBKB13: wantVersion(16),
					},
					EncodingClockDepthMap: map[entities.BoardID]entities.EncodingClockDepth{
						entities.GBKB00: 1,
						entities.GBKB01: 4,
						entities.GBKB03: 7,
						entities.GBKB10: 10,
						entities.GBKB11: 13,
						entities.GBKB13: 16,
					},
					FlushAdcMap: map[entities.BoardID]entities.FlushAdc{
						entities.GBKB00: makeTestMadaEvent(1, 1).FlushADC,
						entities.GBKB01: makeTestMadaEvent(1, 4).FlushADC,
						entities.GBKB03: makeTestMadaEvent(1, 7).FlushADC,
						entities.GBKB10: makeTestMadaEvent(1, 10).FlushADC,
						entities.GBKB11: makeTestMadaEvent(1, 13).FlushADC,
						entities.GBKB13: makeTestMadaEvent(1, 16).FlushADC,
					},
					HitsMap: map[entities.BoardID]entities.RawHits{
						entities.GBKB00: wantHits(1),
						entities.GBKB01: wantHits(4),
						entities.GBKB03: wantHits(7),
						entities.GBKB10: wantHits(10),
						entities.GBKB11: wantHits(13),
						entities.GBKB13: wantHits(16),
					},
				},
				{
					Trigger: 2,
					ClockMap: map[entities.BoardID]entities.ClockCounter{
						entities.GBKB00: 2,
						entities.GBKB01: 5,
						entities.GBKB03: 8,
						entities.GBKB10: 11,
						entities.GBKB11: 14,
						entities.GBKB13: 17,
					},
					InputCh2Map: map[entities.BoardID]entities.InputCh2Counter{
						entities.GBKB00: 2,
						entities.GBKB01: 5,
						entities.GBKB03: 8,
						entities.GBKB10: 11,
						entities.GBKB11: 14,
						entities.GBKB13: 17,
					},
					VersionMap: map[entities.BoardID]entities.Version{
						entities.GBKB00: wantVersion(2),
						entities.GBKB01: wantVersion(5),
						entities.GBKB03: wantVersion(8),
						entities.GBKB10: wantVersion(11),
						entities.GBKB11: wantVersion(14),
						entities.GBKB13: wantVersion(17),
					},
					EncodingClockDepthMap: map[entities.BoardID]entities.EncodingClockDepth{
						entities.GBKB00: 2,
						entities.GBKB01: 5,
						entities.GBKB03: 8,
						entities.GBKB10: 11,
						entities.GBKB11: 14,
						entities.GBKB13: 17,
					},
					FlushAdcMap: map[entities.BoardID]entities.FlushAdc{
						entities.GBKB00: makeTestMadaEvent(2, 2).FlushADC,
						entities.GBKB01: makeTestMadaEvent(2, 5).FlushADC,
						entities.GBKB03: makeTestMadaEvent(2, 8).FlushADC,
						entities.GBKB10: makeTestMadaEvent(2, 11).FlushADC,
						entities.GBKB11: makeTestMadaEvent(2, 14).FlushADC,
						entities.GBKB13: makeTestMadaEvent(2, 17).FlushADC,
					},
					HitsMap: map[entities.BoardID]entities.RawHits{
						entities.GBKB00: wantHits(2),
						entities.GBKB01: wantHits(5),
						entities.GBKB03: wantHits(8),
						entities.GBKB10: wantHits(11),
						entities.GBKB11: wantHits(14),
						entities.GBKB13: wantHits(17),
					},
				},
			},
		},
		{
			name: "成功; 3event, 6board 歯抜けあり",
			input: usecases.MadaMergeCmd{
				MadaEventMap: map[entities.BoardID][]*entities.MadaEvent{
					entities.GBKB00: []*entities.MadaEvent{
						makeTestMadaEvent(0, 0),
						makeTestMadaEvent(1, 1),
						makeTestMadaEvent(2, 2),
					},
					entities.GBKB01: []*entities.MadaEvent{
						makeTestMadaEvent(1, 4),
						makeTestMadaEvent(2, 5),
					},
					entities.GBKB03: []*entities.MadaEvent{
						makeTestMadaEvent(0, 6),
						makeTestMadaEvent(2, 8),
					},
					entities.GBKB10: []*entities.MadaEvent{
						makeTestMadaEvent(0, 9),
						makeTestMadaEvent(1, 10),
						makeTestMadaEvent(2, 11),
					},
					entities.GBKB11: []*entities.MadaEvent{
						makeTestMadaEvent(0, 12),
						makeTestMadaEvent(1, 13),
						makeTestMadaEvent(2, 14),
					},
					entities.GBKB13: []*entities.MadaEvent{
						makeTestMadaEvent(0, 15),
						makeTestMadaEvent(2, 17),
					},
				},
			},
			want: []*entities.RawEvent{
				{
					Trigger: 0,
					ClockMap: map[entities.BoardID]entities.ClockCounter{
						entities.GBKB00: 0,
						entities.GBKB03: 6,
						entities.GBKB10: 9,
						entities.GBKB11: 12,
						entities.GBKB13: 15,
					},
					InputCh2Map: map[entities.BoardID]entities.InputCh2Counter{
						entities.GBKB00: 0,
						entities.GBKB03: 6,
						entities.GBKB10: 9,
						entities.GBKB11: 12,
						entities.GBKB13: 15,
					},
					VersionMap: map[entities.BoardID]entities.Version{
						entities.GBKB00: wantVersion(0),
						entities.GBKB03: wantVersion(6),
						entities.GBKB10: wantVersion(9),
						entities.GBKB11: wantVersion(12),
						entities.GBKB13: wantVersion(15),
					},
					EncodingClockDepthMap: map[entities.BoardID]entities.EncodingClockDepth{
						entities.GBKB00: 0,
						entities.GBKB03: 6,
						entities.GBKB10: 9,
						entities.GBKB11: 12,
						entities.GBKB13: 15,
					},
					FlushAdcMap: map[entities.BoardID]entities.FlushAdc{
						entities.GBKB00: makeTestMadaEvent(0, 0).FlushADC,
						entities.GBKB03: makeTestMadaEvent(0, 6).FlushADC,
						entities.GBKB10: makeTestMadaEvent(0, 9).FlushADC,
						entities.GBKB11: makeTestMadaEvent(0, 12).FlushADC,
						entities.GBKB13: makeTestMadaEvent(0, 15).FlushADC,
					},
					HitsMap: map[entities.BoardID]entities.RawHits{
						entities.GBKB00: wantHits(0),
						entities.GBKB03: wantHits(6),
						entities.GBKB10: wantHits(9),
						entities.GBKB11: wantHits(12),
						entities.GBKB13: wantHits(15),
					},
				},
				{
					Trigger: 1,
					ClockMap: map[entities.BoardID]entities.ClockCounter{
						entities.GBKB00: 1,
						entities.GBKB01: 4,
						entities.GBKB10: 10,
						entities.GBKB11: 13,
					},
					InputCh2Map: map[entities.BoardID]entities.InputCh2Counter{
						entities.GBKB00: 1,
						entities.GBKB01: 4,
						entities.GBKB10: 10,
						entities.GBKB11: 13,
					},
					VersionMap: map[entities.BoardID]entities.Version{
						entities.GBKB00: wantVersion(1),
						entities.GBKB01: wantVersion(4),
						entities.GBKB10: wantVersion(10),
						entities.GBKB11: wantVersion(13),
					},
					EncodingClockDepthMap: map[entities.BoardID]entities.EncodingClockDepth{
						entities.GBKB00: 1,
						entities.GBKB01: 4,
						entities.GBKB10: 10,
						entities.GBKB11: 13,
					},
					FlushAdcMap: map[entities.BoardID]entities.FlushAdc{
						entities.GBKB00: makeTestMadaEvent(1, 1).FlushADC,
						entities.GBKB01: makeTestMadaEvent(1, 4).FlushADC,
						entities.GBKB10: makeTestMadaEvent(1, 10).FlushADC,
						entities.GBKB11: makeTestMadaEvent(1, 13).FlushADC,
					},
					HitsMap: map[entities.BoardID]entities.RawHits{
						entities.GBKB00: wantHits(1),
						entities.GBKB01: wantHits(4),
						entities.GBKB10: wantHits(10),
						entities.GBKB11: wantHits(13),
					},
				},
				{
					Trigger: 2,
					ClockMap: map[entities.BoardID]entities.ClockCounter{
						entities.GBKB00: 2,
						entities.GBKB01: 5,
						entities.GBKB03: 8,
						entities.GBKB10: 11,
						entities.GBKB11: 14,
						entities.GBKB13: 17,
					},
					InputCh2Map: map[entities.BoardID]entities.InputCh2Counter{
						entities.GBKB00: 2,
						entities.GBKB01: 5,
						entities.GBKB03: 8,
						entities.GBKB10: 11,
						entities.GBKB11: 14,
						entities.GBKB13: 17,
					},
					VersionMap: map[entities.BoardID]entities.Version{
						entities.GBKB00: wantVersion(2),
						entities.GBKB01: wantVersion(5),
						entities.GBKB03: wantVersion(8),
						entities.GBKB10: wantVersion(11),
						entities.GBKB11: wantVersion(14),
						entities.GBKB13: wantVersion(17),
					},
					EncodingClockDepthMap: map[entities.BoardID]entities.EncodingClockDepth{
						entities.GBKB00: 2,
						entities.GBKB01: 5,
						entities.GBKB03: 8,
						entities.GBKB10: 11,
						entities.GBKB11: 14,
						entities.GBKB13: 17,
					},
					FlushAdcMap: map[entities.BoardID]entities.FlushAdc{
						entities.GBKB00: makeTestMadaEvent(2, 2).FlushADC,
						entities.GBKB01: makeTestMadaEvent(2, 5).FlushADC,
						entities.GBKB03: makeTestMadaEvent(2, 8).FlushADC,
						entities.GBKB10: makeTestMadaEvent(2, 11).FlushADC,
						entities.GBKB11: makeTestMadaEvent(2, 14).FlushADC,
						entities.GBKB13: makeTestMadaEvent(2, 17).FlushADC,
					},
					HitsMap: map[entities.BoardID]entities.RawHits{
						entities.GBKB00: wantHits(2),
						entities.GBKB01: wantHits(5),
						entities.GBKB03: wantHits(8),
						entities.GBKB10: wantHits(11),
						entities.GBKB11: wantHits(14),
						entities.GBKB13: wantHits(17),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawEvents := usecases.MergeMadaEvents(tt.input)
			if len(rawEvents) != len(tt.want) {
				t.Errorf("want %d, got %d", len(tt.want), len(rawEvents))
			}

			if diff := cmp.Diff(rawEvents, tt.want); diff != "" {
				t.Fatalf("MergeMadaEvents() mismatch (-want +got):\n%s", diff)
			}

			t.Logf("pass: %s", tt.name)
		})
	}
}
