package tachibana

import (
	"reflect"
	"testing"
)

func Test_NumberBool_Bool(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		enum  NumberBool
		want1 bool
	}{
		{name: "未指定 は false", enum: NumberBoolUnspecified, want1: false},
		{name: "False は false", enum: NumberBoolFalse, want1: false},
		{name: "True は true", enum: NumberBoolTrue, want1: true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got1 := test.enum.Bool()
			if !reflect.DeepEqual(test.want1, got1) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want1, got1)
			}
		})
	}
}
