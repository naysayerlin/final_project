package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("Repeat rule dont exist!")
	}
	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("wrong start date %q: %w", dstart, err)
	}
	splits := strings.Split(repeat, " ")
	switch splits[0] {
	case "d":
		return nextDateDay(now, date, splits)
	case "w":
		return nextDateWeek(now, date, splits)
	case "m":
		return nextDateMonth(now, date, splits)
	case "y":
		return nextDateYear(now, date), nil
	default:
		return "", fmt.Errorf("Repeat format is unsupported: %q", splits)
	}
}

func afterNow(date, currentDate time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := currentDate.Date()
	dt := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
	now := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)
	return dt.After(now)
}

func nextDateDay(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("incorrect rule 'd' - %q", strings.Join(parts, " "))
	}
	days, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("incorrect days interval: %q", parts[1])
	}
	if days > 400 || days < 1 {
		return "", fmt.Errorf("incorrect days interval (must be in 1-400 range) : %d", days)
	}
	for {
		date = date.AddDate(0, 0, days)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(dateFormat), nil
}

func nextDateWeek(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("incorrect rule 'w' - %q", strings.Join(parts, " "))
	}
	validDays := make(map[int]bool)
	for _, days := range strings.Split(parts[1], ",") {
		day, err := strconv.Atoi(days)
		if err != nil || day < 1 || day > 7 {
			return "", fmt.Errorf("day cannot be in weekdays: %q", days)
		}
		validDays[day] = true
	}
	for {
		date = date.AddDate(0, 0, 1)
		weekDay := int(date.Weekday())
		if weekDay == 0 {
			weekDay = 7
		}
		if validDays[weekDay] && afterNow(date, now) {
			break
		}
	}
	return date.Format(dateFormat), nil
}

func nextDateMonth(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", fmt.Errorf("incorrect rule 'm': %q", strings.Join(parts, ""))
	}
	validDays := make(map[int]bool)
	for _, days := range strings.Split(parts[1], ",") {
		day, err := strconv.Atoi(days)
		if err != nil || day == 0 || day > 31 || day < -2 {
			return "", fmt.Errorf("day cannot be in month: %q", days)
		}
		validDays[day] = true
	}
	var validMonths map[int]bool
	if len(parts) == 3 {
		validMonths = make(map[int]bool)
		for _, months := range strings.Split(parts[2], ",") {
			month, err := strconv.Atoi(months)
			if err != nil || month > 12 || month < 1 {
				return "", fmt.Errorf("month value doesnt match value of year months : %q", months)
			}
			validMonths[month] = true
		}
	}
	for {
		date = date.AddDate(0, 0, 1)
		if validMonths != nil && !validMonths[int(date.Month())] {
			continue
		}
		day := date.Day()
		lst := lastDayOfMonth(date)
		match := validDays[day] || (day == lst-1 && validDays[-2]) || (day == lst && validDays[-1])
		if match && afterNow(date, now) {
			break
		}
	}
	return date.Format(dateFormat), nil
}

func nextDateYear(now, date time.Time) string {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(dateFormat)
}

func lastDayOfMonth(date time.Time) int {
	firstOfNextMonth := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, date.Location())
	return firstOfNextMonth.AddDate(0, 0, -1).Day()
}
