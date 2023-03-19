package entconst

// constants

const (
	Clock = 1024
)

func EventHeaderSymbol() []byte {
	return []byte{0xeb, 0x90, 0x19, 0x64}
}

func EventFooterSymbol() []byte {
	return []byte{0x75, 0x50, 0x49, 0x43}
}

func FlushAdcHeaderSymbol() []byte {
	return []byte{0x04, 0x05, 0x06, 0x07}
}

const HitHeaderSymbol = uint8(8)

// helper methods

func IsEventHeaderSymbol(b [4]byte) bool {
	expected := EventHeaderSymbol()
	if b[0] == expected[0] && b[1] == expected[1] && b[2] == expected[2] && b[3] == expected[3] {
		return true
	} else {
		return false
	}
}

func IsEventFooterSymbol(b []byte) bool {
	if len(b) != 4 {
		return false
	}

	expected := EventFooterSymbol()
	if b[0] == expected[0] && b[1] == expected[1] && b[2] == expected[2] && b[3] == expected[3] {
		return true
	} else {
		return false
	}
}

func IsAdcHeaderSymbol(ch int, s uint16) bool {
	return s == uint16(FlushAdcHeaderSymbol()[ch])
}

func IsHitHeaderSymbol(b uint8) bool {
	return b == HitHeaderSymbol
}
