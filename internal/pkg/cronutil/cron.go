package cronutil

import (
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	parserStandard = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	parserWithSeconds = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
)

func ParseCron(expr string) (cron.Schedule, error) {
	expr = strings.TrimSpace(expr)

	fields := strings.Fields(expr)

	var schedule cron.Schedule
	var err error

	if len(fields) == 6 {
		schedule, err = parserWithSeconds.Parse(expr)
	} else if len(fields) == 5 {
		schedule, err = parserStandard.Parse(expr)
	} else if len(fields) == 7 {
		schedule, err = cron.NewParser(
			cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
		).Parse(expr)
	} else {
		schedule, err = parserStandard.Parse(expr)
	}

	if err != nil {
		return nil, err
	}

	return schedule, nil
}

func NextTime(expr string, from time.Time) (time.Time, error) {
	schedule, err := ParseCron(expr)
	if err != nil {
		return time.Time{}, err
	}
	return schedule.Next(from), nil
}

func ValidateCron(expr string) error {
	_, err := ParseCron(expr)
	return err
}
