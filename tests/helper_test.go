package tests

func makeTestEventHeader() []byte {
	return []byte{0xeb, 0x90, 0x19, 0x64}
}

// makeTestCounters trigger counter = 10, clock = 11, input ch2 = 12
func makeTestCounters() []byte {
	return []byte{
		0x00, 0x00, 0x00, 0x0a, // trigger counter
		0x00, 0x00, 0x00, 0x0b, // clock counter
		0x00, 0x00, 0x00, 0x0c, // input ch2 counter
	}
}

// makeTestFADC: 全clockで同じ値, ch0 -> 13, ch1 -> 14, ch2 -> 15, ch3 -> 16
func makeTestFADC() []byte {
	ret := make([]byte, 0, 2*4*1024)
	for clock := 0; clock < 1024; clock++ {
		for ch := 0; ch < 4; ch++ {
			header := uint8(ch+4) << 4 // fixed value
			adcValue := uint8(ch + 13) // test random value
			ret = append(ret, header, adcValue)
		}
	}
	return ret
}

// makeTestLuckFADC: ch2だけ空っぽ
func makeTestLuckFADC() []byte {
	ret := make([]byte, 0, 2*4*1024)
	for clock := 0; clock < 1024; clock++ {
		for ch := 0; ch < 4; ch++ {
			if ch == 2 {
				continue
			}
			header := uint8(ch+4) << 4 // fixed value
			adcValue := uint8(ch + 13) // test random value
			ret = append(ret, header, adcValue)
		}
	}
	return ret
}

// makeTestVersionAndDepth year=23, month=3, sub=5, clock depth=50
func makeTestVersionAndDepth() []byte {
	ret := make([]byte, 0, 4)
	versionYear := uint8(23)
	versionMonth := uint8(3)
	versionSub := uint8(5) << 4
	encodingClockDepth := uint8(50)
	ret = append(ret, versionYear, versionMonth, versionSub, encodingClockDepth)
	return ret
}

// makeTestHit clock=50,51,52 だけ, hit = all true
func makeTestHit() []byte {
	ret := make([]byte, 0, 20*3)
	for clock := 50; clock < 53; clock++ {
		header := uint8(8 << 4)
		ret = append(ret, header, 0)
		ret = append(ret, uint8(clock>>8), uint8(clock&0xff))
		// 128ch 全部 1
		for i := 0; i < 16; i++ {
			ret = append(ret, 0xff)
		}
	}
	return ret
}

func makeTestEventFooter() []byte {
	return []byte{0x75, 0x50, 0x49, 0x43}
}
