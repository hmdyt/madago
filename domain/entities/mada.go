package entities

import "github.com/hmdyt/madago/domain/entconst"

type Event struct {
	Header EventHeader
}

type EventHeader struct {
	Counter            EventCounter
	Clock              ClockCounter
	FlushADC           EventFlushAdc
	Version            Version
	EncodingClockDepth EncodingClockDepth
}

type EventCounter uint32
type ClockCounter uint32
type EventFlushAdc [4][entconst.Clock]uint16
type Version struct {
	Year  uint8
	Month uint8
	Sub   uint8
}
type EncodingClockDepth uint16
