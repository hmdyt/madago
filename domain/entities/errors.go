package entities

import (
	"fmt"
	"strings"
)

type InvalidHeaderError struct {
	Got [4]byte
}

func (e InvalidHeaderError) Error() string {
	byteStrings := make([]string, 4)
	for i, b := range e.Got {
		byteStrings[i] = fmt.Sprintf("0x%02x", b)
	}
	return fmt.Sprintf("invalid header error got=[%s]", strings.Join(byteStrings, ", "))
}
