package days

import (
	"reflect"
	"testing"
	"time"
)

func TestDaysBetween(t *testing.T) {
	tests := []struct {
		name string
		from time.Time
		to   time.Time
		want []time.Time
	}{
		{
			name: "Normal range",
			from: time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			to:   time.Date(2023, 8, 5, 0, 0, 0, 0, time.UTC),
			want: []time.Time{
				time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 8, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 8, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 8, 4, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Same date",
			from: time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			to:   time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			want: nil,
		},
		{
			name: "From date after to date",
			from: time.Date(2023, 8, 5, 0, 0, 0, 0, time.UTC),
			to:   time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			want: nil,
		},
		{
			name: "One day range",
			from: time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			to:   time.Date(2023, 8, 2, 0, 0, 0, 0, time.UTC),
			want: []time.Time{
				time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Leap year",
			from: time.Date(2020, 2, 28, 0, 0, 0, 0, time.UTC),
			to:   time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC),
			want: []time.Time{
				time.Date(2020, 2, 28, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "Crossing year boundary",
			from: time.Date(2023, 12, 30, 0, 0, 0, 0, time.UTC),
			to:   time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			want: []time.Time{
				time.Date(2023, 12, 30, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Between(tt.from, tt.to); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Between() = %v, want %v", got, tt.want)
			}
		})
	}
}
