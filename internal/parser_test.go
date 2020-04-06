package internal

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		maxNewLine int32
	}
	tests := []struct {
		name    string
		args    args
		want    *Parser
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				maxNewLine: 1,
			},
			want:    &Parser{MaxNewLine: 1},
			wantErr: false,
		},
		{
			name: "should err with negative < 0 number",
			args: args{
				maxNewLine: -2,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.maxNewLine)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}
