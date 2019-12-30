package usecase

import (
	"errors"
	"fmt"
	"time"
)

func (a *app) parseAndGetEndOfDay(date string) (time.Time, error) {
	end, err := time.Parse(timeParseLayout, date)
	if err != nil {
		return time.Now(), errors.New(fmt.Sprintf("Failed to parse date format : %s", date))
	}
	return end.Add(endOfDayDuration), nil // date must be with 23:59:59 time in the end
}

func (a *app) getStartOfPeriod(end time.Time, period string) (time.Time, error) {
	var start time.Time
	switch period {
	case "WEEK":
		start = end.AddDate(0, 0, -7)
	case "MONTH":
		start = end.AddDate(0, -1, 0)
	case "YEAR":
		start = end.AddDate(-1, 0, 0)
	default:
		return time.Now(), errors.New(fmt.Sprintf("Incorrect format : %s", period))
	}
	return start, nil
}
