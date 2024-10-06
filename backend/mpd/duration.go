package mpd

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ISO8601Duration time.Duration

func (i *ISO8601Duration) UnmarshalText(text []byte) error {
	s := string(text)
	s = strings.TrimPrefix(s, "PT")
	s = strings.TrimSuffix(s, "S")

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	*i = ISO8601Duration(f * float64(time.Second))

	return nil
}

func (i ISO8601Duration) MarshalText() ([]byte, error) {
	strSec := strconv.FormatFloat(time.Duration(i).Seconds(), 'f', -1, 64)
	s := fmt.Sprintf("PT%sS", strSec)
	return []byte(s), nil
}
