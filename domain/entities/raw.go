package entities

import "github.com/hmdyt/madago/domain/entconst"

type BoardID string

const (
	GBKB00 BoardID = "GBKB00"
	GBKB01 BoardID = "GBKB01"
	GBKB03 BoardID = "GBKB03"
	GBKB10 BoardID = "GBKB10"
	GBKB11 BoardID = "GBKB11"
	GBKB13 BoardID = "GBKB13"
)

func NewBoardID(boardID string) (BoardID, error) {
	if err := validateBoardID(BoardID(boardID)); err != nil {
		return "", err
	}
	return BoardID(boardID), nil
}

type RawEvent struct {
	Trigger               TriggerCounter
	ClockMap              map[BoardID]ClockCounter
	InputCh2Map           map[BoardID]InputCh2Counter
	VersionMap            map[BoardID]Version
	EncodingClockDepthMap map[BoardID]EncodingClockDepth
	FlushAdcMap           map[BoardID]FlushAdc
	HitsMap               map[BoardID]RawHits
}

type RawHits [entconst.Clock][128]bool

func RawHitsFromMadaHits(madaHits []MadaHit) RawHits {
	var rawHits RawHits
	for _, hit := range madaHits {
		rawHits[hit.Clock] = hit.IsHit
	}
	return rawHits
}

func validateBoardID(boardID BoardID) error {
	switch boardID {
	case GBKB00, GBKB01, GBKB03, GBKB10, GBKB11, GBKB13:
		return nil
	default:
		return InvalidBoardIDError{Got: string(boardID)}
	}
}
