package usecases

import (
	"sort"

	"github.com/hmdyt/madago/domain/entities"
)

type MadaMergeCmd struct {
	MadaEventMap map[entities.BoardID][]*entities.MadaEvent
}

func MergeMadaEvents(cmd MadaMergeCmd) []*entities.RawEvent {
	var maxEvents int
	for _, madaEvents := range cmd.MadaEventMap {
		if maxEvents < len(madaEvents) {
			maxEvents = len(madaEvents)
		}
	}

	// trigger ID がキーのMapでmergeする
	rawEventMap := make(map[entities.TriggerCounter]*entities.RawEvent)
	for boardID, madaEvents := range cmd.MadaEventMap {
		for _, madaEvent := range madaEvents {
			if rawEvent, ok := rawEventMap[madaEvent.Trigger]; ok {
				rawEvent.ClockMap[boardID] = madaEvent.Clock
				rawEvent.InputCh2Map[boardID] = madaEvent.InputCh2
				rawEvent.VersionMap[boardID] = madaEvent.Version
				rawEvent.EncodingClockDepthMap[boardID] = madaEvent.EncodingClockDepth
				rawEvent.FlushAdcMap[boardID] = madaEvent.FlushADC
				rawEvent.HitsMap[boardID] = entities.RawHitsFromMadaHits(madaEvent.Hits)
			} else {
				rawEventMap[madaEvent.Trigger] = &entities.RawEvent{
					Trigger:               madaEvent.Trigger,
					ClockMap:              map[entities.BoardID]entities.ClockCounter{boardID: madaEvent.Clock},
					InputCh2Map:           map[entities.BoardID]entities.InputCh2Counter{boardID: madaEvent.InputCh2},
					VersionMap:            map[entities.BoardID]entities.Version{boardID: madaEvent.Version},
					EncodingClockDepthMap: map[entities.BoardID]entities.EncodingClockDepth{boardID: madaEvent.EncodingClockDepth},
					FlushAdcMap:           map[entities.BoardID]entities.FlushAdc{boardID: madaEvent.FlushADC},
					HitsMap:               map[entities.BoardID]entities.RawHits{boardID: entities.RawHitsFromMadaHits(madaEvent.Hits)},
				}
			}
		}
	}

	// map to slice
	rawEvents := make([]*entities.RawEvent, len(rawEventMap))
	for i, rawEvent := range rawEventMap {
		rawEvents[i] = rawEvent
	}

	// sort by trigger id
	sort.Slice(rawEvents, func(i, j int) bool {
		return rawEvents[i].Trigger < rawEvents[j].Trigger
	})

	return rawEvents
}
