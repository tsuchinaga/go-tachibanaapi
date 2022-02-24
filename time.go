package tachibana

import (
	"time"
)

type RequestTime struct {
	time.Time
}

func (t RequestTime) MarshalJSON() ([]byte, error) {
	return []byte(t.Time.Format(`"2006.01.02-15:04:05.000"`)), nil
}

func (t *RequestTime) UnmarshalJSON(b []byte) error {
	if b == nil || string(b) == `""` || string(b) == "null" {
		return nil
	}
	tt, err := time.Parse(`"2006.01.02-15:04:05.000"`, string(b))
	if err != nil {
		return err
	}
	*t = RequestTime{Time: time.Date(tt.Year(), tt.Month(), tt.Day(), tt.Hour(), tt.Minute(), tt.Second(), tt.Nanosecond(), time.Local)}
	return nil
}

type YmdHms struct {
	time.Time
}

func (t YmdHms) MarshalJSON() ([]byte, error) {
	return []byte(t.Time.Format(`"20060102150405"`)), nil
}

func (t *YmdHms) UnmarshalJSON(b []byte) error {
	if b == nil || string(b) == `""` || string(b) == "null" {
		return nil
	}
	tt, err := time.Parse(`"20060102150405"`, string(b))
	if err != nil {
		return err
	}
	*t = YmdHms{Time: time.Date(tt.Year(), tt.Month(), tt.Day(), tt.Hour(), tt.Minute(), tt.Second(), tt.Nanosecond(), time.Local)}
	return nil
}

type Ymd struct {
	time.Time
}

func (t Ymd) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte(`"0"`), nil
	}
	return []byte(t.Time.Format(`"20060102"`)), nil
}

func (t *Ymd) UnmarshalJSON(b []byte) error {
	if b == nil || string(b) == `""` || string(b) == "null" {
		return nil
	}
	tt, err := time.Parse(`"20060102"`, string(b))
	if err != nil {
		return err
	}
	*t = Ymd{Time: time.Date(tt.Year(), tt.Month(), tt.Day(), tt.Hour(), tt.Minute(), tt.Second(), tt.Nanosecond(), time.Local)}
	return nil
}
