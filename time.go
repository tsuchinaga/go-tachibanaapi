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

type Ym struct {
	time.Time
}

func (t Ym) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(t.Time.Format(`"200601"`)), nil
}

func (t *Ym) UnmarshalJSON(b []byte) error {
	str := string(b)
	if b == nil || str == `""` || str == "null" || str == `"0"` || str == `"000000"` {
		return nil
	}

	tt, err := time.Parse(`"200601"`, string(b))
	if err != nil {
		return err
	}
	*t = Ym{Time: time.Date(tt.Year(), tt.Month(), 1, 0, 0, 0, 0, time.Local)}
	return nil
}

type Ymd struct {
	time.Time
	isNoChange bool
	isToday    bool
}

func (t Ymd) MarshalJSON() ([]byte, error) {
	if t.isNoChange {
		return []byte(`"*"`), nil
	}
	if t.isToday {
		return []byte(`"0"`), nil
	}
	if t.Time.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(t.Time.Format(`"20060102"`)), nil
}

func (t *Ymd) UnmarshalJSON(b []byte) error {
	str := string(b)
	if b == nil || str == `""` || str == "null" || str == `"0"` || str == `"00000000"` {
		return nil
	}

	tt, err := time.Parse(`"20060102"`, string(b))
	if err != nil {
		return err
	}
	*t = Ymd{Time: time.Date(tt.Year(), tt.Month(), tt.Day(), 0, 0, 0, 0, time.Local)}
	return nil
}

type YmdHms struct {
	time.Time
}

func (t YmdHms) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(t.Time.Format(`"20060102150405"`)), nil
}

func (t *YmdHms) UnmarshalJSON(b []byte) error {
	str := string(b)
	if b == nil || str == `""` || str == "null" || str == `"00000000000000"` {
		return nil
	}
	tt, err := time.Parse(`"20060102150405"`, str)
	if err != nil {
		return err
	}
	*t = YmdHms{Time: time.Date(tt.Year(), tt.Month(), tt.Day(), tt.Hour(), tt.Minute(), tt.Second(), 0, time.Local)}
	return nil
}

type YmdHm struct {
	time.Time
}

func (t YmdHm) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(t.Time.Format(`"200601021504"`)), nil
}

func (t *YmdHm) UnmarshalJSON(b []byte) error {
	str := string(b)
	if b == nil || str == `""` || str == "null" || str == `"000000000000"` {
		return nil
	}
	tt, err := time.Parse(`"200601021504"`, str)
	if err != nil {
		return err
	}
	*t = YmdHm{Time: time.Date(tt.Year(), tt.Month(), tt.Day(), tt.Hour(), tt.Minute(), 0, 0, time.Local)}
	return nil
}

type Hm struct {
	time.Time
}

func (t Hm) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(t.Time.Format(`"15:04"`)), nil
}

func (t *Hm) UnmarshalJSON(b []byte) error {
	str := string(b)
	if b == nil || str == `""` || str == "null" {
		return nil
	}
	tt, err := time.Parse(`"15:04"`, str)
	if err != nil {
		return err
	}
	*t = Hm{Time: time.Date(0, 1, 1, tt.Hour(), tt.Minute(), 0, 0, time.Local)}
	return nil
}

type Hms struct {
	time.Time
}

func (t Hms) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(t.Time.Format(`"150405"`)), nil
}

func (t *Hms) UnmarshalJSON(b []byte) error {
	str := string(b)
	if b == nil || str == `""` || str == "null" {
		return nil
	}
	tt, err := time.Parse(`"150405"`, str)
	if err != nil {
		return err
	}
	*t = Hms{Time: time.Date(0, 1, 1, tt.Hour(), tt.Minute(), tt.Second(), 0, time.Local)}
	return nil
}
