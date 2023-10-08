package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseMDYYYYToDate(t *testing.T) {

	location, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Printf("Error getting asia/kolkata location")
		panic("unable to load asia kolkata location")
	}

	type args struct {
		dateString string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "Test ParseMDYYYYToDate for 9/23/2023",
			args: args{dateString: "9/23/2023"},
			want: time.Date(2023, time.September, 23, 0, 0, 0, 0, location),
		},
		{
			name: "Test ParseMDYYYYToDate for 12/9/2023",
			args: args{dateString: "12/9/2023"},
			want: time.Date(2023, time.December, 9, 0, 0, 0, 0, location),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseMDYYYYToDate(tt.args.dateString); !assert.Equal(t, got, tt.want) {
				t.Errorf("ParseMDYYYYToDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
