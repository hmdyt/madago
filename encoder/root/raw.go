package root

import (
	"log"
	"reflect"

	"github.com/cheggaaa/pb/v3"
	"github.com/hmdyt/madago/domain/entconst"
	"github.com/hmdyt/madago/domain/entities"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rdict"
	"go-hep.org/x/hep/groot/rtree"
)

type RawEncoder struct {
	file *groot.File
	pb   *pb.ProgressBar
}

// RawEvent6Boards 6 boards implement
type RawEvent6Boards struct {
	TriggerCounter uint32 `groot:"trigger_counter"`

	ClockCounter00 uint32 `groot:"clock_counter_00"`
	ClockCounter01 uint32 `groot:"clock_counter_01"`
	ClockCounter03 uint32 `groot:"clock_counter_03"`
	ClockCounter10 uint32 `groot:"clock_counter_10"`
	ClockCounter11 uint32 `groot:"clock_counter_11"`
	ClockCounter13 uint32 `groot:"clock_counter_13"`

	InputCh2Counter00 uint32 `groot:"input_ch2_counter_00"`
	InputCh2Counter01 uint32 `groot:"input_ch2_counter_01"`
	InputCh2Counter03 uint32 `groot:"input_ch2_counter_03"`
	InputCh2Counter10 uint32 `groot:"input_ch2_counter_10"`
	InputCh2Counter11 uint32 `groot:"input_ch2_counter_11"`
	InputCh2Counter13 uint32 `groot:"input_ch2_counter_13"`

	EncodingClockDepth00 uint16 `groot:"encoding_clock_depth_00"`
	EncodingClockDepth01 uint16 `groot:"encoding_clock_depth_01"`
	EncodingClockDepth03 uint16 `groot:"encoding_clock_depth_03"`
	EncodingClockDepth10 uint16 `groot:"encoding_clock_depth_10"`
	EncodingClockDepth11 uint16 `groot:"encoding_clock_depth_11"`
	EncodingClockDepth13 uint16 `groot:"encoding_clock_depth_13"`

	FlushADC00 [4][entconst.Clock]uint16 `groot:"fadc_00"`
	FlushADC01 [4][entconst.Clock]uint16 `groot:"fadc_01"`
	FlushADC03 [4][entconst.Clock]uint16 `groot:"fadc_03"`
	FlushADC10 [4][entconst.Clock]uint16 `groot:"fadc_10"`
	FlushADC11 [4][entconst.Clock]uint16 `groot:"fadc_11"`
	FlushADC13 [4][entconst.Clock]uint16 `groot:"fadc_13"`

	Hit00 [entconst.Clock][128]bool `groot:"hit_00"`
	Hit01 [entconst.Clock][128]bool `groot:"hit_01"`
	Hit03 [entconst.Clock][128]bool `groot:"hit_03"`
	Hit10 [entconst.Clock][128]bool `groot:"hit_10"`
	Hit11 [entconst.Clock][128]bool `groot:"hit_11"`
	Hit13 [entconst.Clock][128]bool `groot:"hit_13"`
}

func rawFromEntity(event *RawEvent6Boards, ent *entities.RawEvent) {
	event.TriggerCounter = uint32(ent.Trigger)
	log.Println("trigger counter: ", event.TriggerCounter)

	event.ClockCounter00 = uint32(getOr(ent.ClockMap, entities.GBKB00, 0))
	event.ClockCounter01 = uint32(getOr(ent.ClockMap, entities.GBKB01, 0))
	event.ClockCounter03 = uint32(getOr(ent.ClockMap, entities.GBKB03, 0))
	event.ClockCounter10 = uint32(getOr(ent.ClockMap, entities.GBKB10, 0))
	event.ClockCounter11 = uint32(getOr(ent.ClockMap, entities.GBKB11, 0))
	event.ClockCounter13 = uint32(getOr(ent.ClockMap, entities.GBKB13, 0))
	log.Println("clock counter: ", event.ClockCounter00, event.ClockCounter01, event.ClockCounter03, event.ClockCounter10, event.ClockCounter11, event.ClockCounter13)

	event.InputCh2Counter00 = uint32(getOr(ent.InputCh2Map, entities.GBKB00, 0))
	event.InputCh2Counter01 = uint32(getOr(ent.InputCh2Map, entities.GBKB01, 0))
	event.InputCh2Counter03 = uint32(getOr(ent.InputCh2Map, entities.GBKB03, 0))
	event.InputCh2Counter10 = uint32(getOr(ent.InputCh2Map, entities.GBKB10, 0))
	event.InputCh2Counter11 = uint32(getOr(ent.InputCh2Map, entities.GBKB11, 0))
	event.InputCh2Counter13 = uint32(getOr(ent.InputCh2Map, entities.GBKB13, 0))
	log.Println("input ch2 counter: ", event.InputCh2Counter00, event.InputCh2Counter01, event.InputCh2Counter03, event.InputCh2Counter10, event.InputCh2Counter11, event.InputCh2Counter13)

	event.EncodingClockDepth00 = uint16(getOr(ent.EncodingClockDepthMap, entities.GBKB00, 0))
	event.EncodingClockDepth01 = uint16(getOr(ent.EncodingClockDepthMap, entities.GBKB01, 0))
	event.EncodingClockDepth03 = uint16(getOr(ent.EncodingClockDepthMap, entities.GBKB03, 0))
	event.EncodingClockDepth10 = uint16(getOr(ent.EncodingClockDepthMap, entities.GBKB10, 0))
	event.EncodingClockDepth11 = uint16(getOr(ent.EncodingClockDepthMap, entities.GBKB11, 0))
	event.EncodingClockDepth13 = uint16(getOr(ent.EncodingClockDepthMap, entities.GBKB13, 0))
	log.Println("encoding clock depth: ", event.EncodingClockDepth00, event.EncodingClockDepth01, event.EncodingClockDepth03, event.EncodingClockDepth10, event.EncodingClockDepth11, event.EncodingClockDepth13)

	var defaultFADC entities.FlushAdc
	event.FlushADC00 = getOr(ent.FlushAdcMap, entities.GBKB00, defaultFADC)
	event.FlushADC01 = getOr(ent.FlushAdcMap, entities.GBKB01, defaultFADC)
	event.FlushADC03 = getOr(ent.FlushAdcMap, entities.GBKB03, defaultFADC)
	event.FlushADC10 = getOr(ent.FlushAdcMap, entities.GBKB10, defaultFADC)
	event.FlushADC11 = getOr(ent.FlushAdcMap, entities.GBKB11, defaultFADC)
	event.FlushADC13 = getOr(ent.FlushAdcMap, entities.GBKB13, defaultFADC)
	log.Println("FADC done")

	var defaultHit entities.RawHits
	event.Hit00 = getOr(ent.HitsMap, entities.GBKB00, defaultHit)
	event.Hit01 = getOr(ent.HitsMap, entities.GBKB01, defaultHit)
	event.Hit03 = getOr(ent.HitsMap, entities.GBKB03, defaultHit)
	event.Hit10 = getOr(ent.HitsMap, entities.GBKB10, defaultHit)
	event.Hit11 = getOr(ent.HitsMap, entities.GBKB11, defaultHit)
	event.Hit13 = getOr(ent.HitsMap, entities.GBKB13, defaultHit)
	log.Println("hit done")

}

func NewRawEncoder(f *groot.File, progressBar *pb.ProgressBar) *RawEncoder {
	rdict.StreamerInfos.Add(rdict.StreamerOf(
		rdict.StreamerInfos,
		reflect.TypeOf(RawEvent6Boards{}),
	))

	return &RawEncoder{
		file: f,
		pb:   progressBar,
	}
}

func (r *RawEncoder) Encode(rawEvents []*entities.RawEvent) error {
	// make tree
	var eventBuf RawEvent6Boards
	tree, err := rtree.NewWriter(r.file, "raw", rtree.WriteVarsFromStruct(&eventBuf))
	if err != nil {
		log.Printf("could not create tree writer: %+v", err)
		return err
	}

	// fill
	for _, ent := range rawEvents {
		rawFromEntity(&eventBuf, ent)
		if n, err := tree.Write(); err != nil {
			log.Printf("could not fill tree: %+v", err)
			return err
		} else {
			log.Println("bytes: ", n)
		}
		r.pb.Increment()
	}

	if err := tree.Close(); err != nil {
		return err
	}
	return nil
}

func getOr[T any](m map[entities.BoardID]T, k entities.BoardID, defaultValue T) T {
	if v, ok := m[k]; ok {
		return v
	} else {
		return defaultValue
	}
}
