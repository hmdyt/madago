package entities

import "github.com/hmdyt/madago/domain/entconst"

type Event struct {
	Header EventHeader
	Hits   []Hit
}

type EventHeader struct {
	Trigger            TriggerCounter
	Clock              ClockCounter
	InputCh2           InputCh2Counter
	FlushADC           EventFlushAdc
	Version            Version
	EncodingClockDepth EncodingClockDepth
}

type TriggerCounter uint32
type ClockCounter uint32
type InputCh2Counter uint32
type EventFlushAdc [4][entconst.Clock]uint16
type Version struct {
	Year  uint8
	Month uint8
	Sub   uint8
}
type EncodingClockDepth uint16

type Hit struct {
	Clock uint16
	IsHit [128]bool
}
