package root

import (
	"log"
	"reflect"

	"github.com/cheggaaa/pb/v3"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rdict"
	"go-hep.org/x/hep/groot/rtree"

	"github.com/hmdyt/madago/domain/entconst"
	"github.com/hmdyt/madago/domain/entities"
)

type Encoder struct {
	file *groot.File
	pb   *pb.ProgressBar
}

type Event struct {
	TriggerCounter     uint32                    `groot:"trigger_counter"`
	ClockCounter       uint32                    `groot:"clock_counter"`
	InputCh2Counter    uint32                    `groot:"input_ch2_counter"`
	FlushADC           [4][entconst.Clock]uint16 `groot:"fadc"`
	EncodingClockDepth uint16                    `groot:"encoding_clock_depth"`
	Hit                [entconst.Clock][128]bool `groot:"hit"`
}

func eventFromEntity(event *Event, ent *entities.Event) {
	event.TriggerCounter = uint32(ent.Header.Trigger)
	event.ClockCounter = uint32(ent.Header.Clock)
	event.InputCh2Counter = uint32(ent.Header.InputCh2)
	event.FlushADC = ent.Header.FlushADC
	event.EncodingClockDepth = uint16(ent.Header.EncodingClockDepth)
	for _, hit := range ent.Hits {
		event.Hit[hit.Clock] = hit.IsHit
	}
}

func NewEncoder(f *groot.File, progressBar *pb.ProgressBar) (Encoder, error) {
	rdict.StreamerInfos.Add(rdict.StreamerOf(
		rdict.StreamerInfos,
		reflect.TypeOf(Event{}),
	))

	return Encoder{file: f, pb: progressBar}, nil
}

func (r *Encoder) Write(events []*entities.Event) error {
	// make tree
	var eventBuf Event
	tree, err := rtree.NewWriter(r.file, "tree", rtree.WriteVarsFromStruct(&eventBuf))
	if err != nil {
		log.Printf("could not create tree writer: %+v", err)
		return err
	}

	// fill
	for i, e := range events {
		eventFromEntity(&eventBuf, e)
		if _, err := tree.Write(); err != nil {
			log.Printf("could not write event %d: %+v", i, err)
			return err
		}
		r.pb.Increment()
	}

	if err := tree.Close(); err != nil {
		return err
	}
	return nil
}
