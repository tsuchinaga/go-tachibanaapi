package tachibana

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func Test_RequestTime_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		time RequestTime
		want []byte
	}{
		{name: "正常な日付をパースできる",
			time: RequestTime{Time: time.Date(2022, 2, 10, 9, 30, 15, 123000000, time.Local)},
			want: []byte(`"2022.02.10-09:30:15.123"`)},
		{name: "time.Timeのゼロ値はゼロ値になる",
			time: RequestTime{},
			want: []byte(`"0001.01.01-00:00:00.000"`)},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := json.Marshal(test.time)
			if !reflect.DeepEqual(test.want, got) || err != nil {
				t.Errorf("%s error\nwant: %s\ngot: %s, %+v\n", t.Name(), test.want, got, err)
			}
		})
	}
}

func Test_RequestTime_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		src      []byte
		want1    RequestTime
		hasError bool
	}{
		{name: "正常系のパース",
			src:      []byte(`"2022.02.10-09:30:15.123"`),
			want1:    RequestTime{Time: time.Date(2022, 2, 10, 9, 30, 15, 123000000, time.Local)},
			hasError: false},
		{name: "nullはゼロ値にする",
			src:      []byte(`null`),
			want1:    RequestTime{},
			hasError: false},
		{name: "空文字はゼロ値にする",
			src:      []byte(`""`),
			want1:    RequestTime{},
			hasError: false},
		{name: "違う形式だとエラー",
			src:      []byte(`"2022-02-24"`),
			want1:    RequestTime{},
			hasError: true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := RequestTime{}
			err := json.Unmarshal(test.src, &got)
			if !reflect.DeepEqual(test.want1, got) || (err != nil) != test.hasError {
				t.Errorf("%s error\nwant: %v, %v\ngot: %v, %v\n", t.Name(), test.want1, test.hasError, got, err)
			}
		})
	}
}

func Test_YmdHms_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		time YmdHms
		want []byte
	}{
		{name: "正常な日付をパースできる",
			time: YmdHms{Time: time.Date(2022, 2, 10, 9, 30, 15, 123000000, time.Local)},
			want: []byte(`"20220210093015"`)},
		{name: "time.Timeのゼロ値はゼロ値になる",
			time: YmdHms{},
			want: []byte(`"00010101000000"`)},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := json.Marshal(test.time)
			if !reflect.DeepEqual(test.want, got) || err != nil {
				t.Errorf("%s error\nwant: %s, %+v\ngot: %s\n", t.Name(), test.want, err, got)
			}
		})
	}
}

func Test_YmdHms_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		src      []byte
		want1    YmdHms
		hasError bool
	}{
		{name: "正常系のパース",
			src:      []byte(`"20220210093015.123000000"`),
			want1:    YmdHms{Time: time.Date(2022, 2, 10, 9, 30, 15, 123000000, time.Local)},
			hasError: false},
		{name: "nullはゼロ値にする",
			src:      []byte(`null`),
			want1:    YmdHms{},
			hasError: false},
		{name: "空文字はゼロ値にする",
			src:      []byte(`""`),
			want1:    YmdHms{},
			hasError: false},
		{name: "違う形式だとエラー",
			src:      []byte(`"2022-02-24"`),
			want1:    YmdHms{},
			hasError: true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := YmdHms{}
			err := json.Unmarshal(test.src, &got)
			if !reflect.DeepEqual(test.want1, got) || (err != nil) != test.hasError {
				t.Errorf("%s error\nwant: %v, %v\ngot: %v, %v\n", t.Name(), test.want1, test.hasError, got, err)
			}
		})
	}
}

func Test_Ymd_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		time Ymd
		want []byte
	}{
		{name: "正常な日付をパースできる",
			time: Ymd{Time: time.Date(2022, 2, 10, 9, 30, 15, 123000000, time.Local)},
			want: []byte(`"20220210"`)},
		{name: "time.Timeのゼロ値は0になる",
			time: Ymd{},
			want: []byte(`"0"`)},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := json.Marshal(test.time)
			if !reflect.DeepEqual(test.want, got) || err != nil {
				t.Errorf("%s error\nwant: %s, %+v\ngot: %s\n", t.Name(), test.want, err, got)
			}
		})
	}
}

func Test_Ymd_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		src      []byte
		want1    Ymd
		hasError bool
	}{
		{name: "正常系のパース",
			src:      []byte(`"20220210"`),
			want1:    Ymd{Time: time.Date(2022, 2, 10, 0, 0, 0, 0, time.Local)},
			hasError: false},
		{name: "nullはゼロ値にする",
			src:      []byte(`null`),
			want1:    Ymd{},
			hasError: false},
		{name: "空文字はゼロ値にする",
			src:      []byte(`""`),
			want1:    Ymd{},
			hasError: false},
		{name: "違う形式だとエラー",
			src:      []byte(`"2022-02-24"`),
			want1:    Ymd{},
			hasError: true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := Ymd{}
			err := json.Unmarshal(test.src, &got)
			if !reflect.DeepEqual(test.want1, got) || (err != nil) != test.hasError {
				t.Errorf("%s error\nwant: %v, %v\ngot: %v, %v\n", t.Name(), test.want1, test.hasError, got, err)
			}
		})
	}
}
