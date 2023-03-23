package entities

import (
	"fmt"
	"strings"
)

type InvalidHeaderError struct {
	Got      []byte
	Expected []byte
}

func (e InvalidHeaderError) Error() string {
	byteStrings := make([]string, 4)
	for i, b := range e.Got {
		byteStrings[i] = fmt.Sprintf("0x%02x", b)
	}
	return fmt.Sprintf("invalid header error got=[%s]", strings.Join(byteStrings, ", "))
}

type InvalidFlushAdcHeaderError struct {
	Got      uint16
	Expected byte
}

func (e InvalidFlushAdcHeaderError) Error() string {
	return fmt.Sprintf("invalid flush adc header error got=%d, expected=%d", e.Got, e.Expected)
}
