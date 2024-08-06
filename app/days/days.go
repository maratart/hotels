package days

import "time"

func Between(from time.Time, to time.Time) []time.Time {
	if from.After(to) || from.Equal(to) || to.Before(from) {
		return nil
	}

	days := make([]time.Time, 0)
	for d := toDay(from); d.Before(toDay(to)); d = d.AddDate(0, 0, 1) {
		days = append(days, d)
	}

	return days
}

func toDay(timestamp time.Time) time.Time {
	return time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, time.UTC)
}

func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
