package entconst

func HeaderSymbol() []byte {
	return []byte{0xeb, 0x90, 0x19, 0x64}
}

func IsValidHeaderSymbol(b [4]byte) bool {
	expected := HeaderSymbol()
	if b[0] == expected[0] && b[1] == expected[1] && b[2] == expected[2] && b[3] == expected[3] {
		return true
	} else {
		return false
	}
}
