package entities

import "github.com/hmdyt/madago/domain/entconst"

type MadaEvent struct {
	Trigger            TriggerCounter
	Clock              ClockCounter
	InputCh2           InputCh2Counter
	Version            Version
	EncodingClockDepth EncodingClockDepth
	FlushADC           FlushAdc
	Hits               []MadaHit
}

type TriggerCounter uint32
type ClockCounter uint32
type InputCh2Counter uint32
type Version struct {
	Year  uint8
	Month uint8
	Sub   uint8
}
type EncodingClockDepth uint16
type FlushAdc [4][entconst.Clock]uint16

type MadaHit struct {
	Clock uint16
	IsHit [128]bool
}
