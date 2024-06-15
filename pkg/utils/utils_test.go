package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_getTimeBeginningOfTheDay(t *testing.T) {
	type args struct {
		timezone string
	}
	location, _ := time.LoadLocation("Asia/Jakarta")
	wantErr1 := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, location)
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr *time.Time
	}{
		// TODO: Add test cases.
		{
			name:    "Should Return start day of local in UTC timezone",
			args:    args{timezone: "UTC"},
			want:    time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).In(time.UTC),
			wantErr: nil,
		},
		{
			name:    "Should Return start day of local in UTC timezone when pass nil",
			want:    time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).In(time.UTC),
			wantErr: nil,
		},
		{
			name:    "this wantErr1 only set timezone . not the time at the timezone",
			args:    args{timezone: "UTC"},
			want:    time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).In(time.UTC),
			wantErr: &wantErr1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := GetTimeBeginningOfTheDay(&tt.args.timezone)
			fmt.Printf("== Debug Output ==\n")
			fmt.Printf("Got    : %s\n", got)
			fmt.Printf("Want   : %s\n", tt.want)
			fmt.Printf("WantErr: %s\n", tt.wantErr)
			if tt.wantErr == nil {
				assert.Equalf(t, &tt.want, got, "getTimeBeginingOfTheDay(%v)", tt.args.timezone)
			} else {
				assert.NotEqualf(t, tt.wantErr, got, "getTimeBeginingOfTheDay(%v)", tt.args.timezone)
			}
		})
	}
}
