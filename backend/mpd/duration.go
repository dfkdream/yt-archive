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
	s := fmt.Sprintf("PT%gS", time.Duration(i).Seconds())
	return []byte(s), nil
}
