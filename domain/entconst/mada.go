package entconst

const (
	Clock = 1024
)

func HeaderSymbol() []byte {
	return []byte{0xeb, 0x90, 0x19, 0x64}
}

func FlushAdcHeaderSymbol() []byte {
	return []byte{0x04, 0x05, 0x06, 0x07}
}

func IsValidHeaderSymbol(b [4]byte) bool {
	expected := HeaderSymbol()
	if b[0] == expected[0] && b[1] == expected[1] && b[2] == expected[2] && b[3] == expected[3] {
		return true
	} else {
		return false
	}
}

func IsValidAdcHeaderSymbol(ch int, s uint16) bool {
	return s == uint16(FlushAdcHeaderSymbol()[ch])
}
