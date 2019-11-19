package cronrange

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	cronParseOption = cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow
	cronParser      = cron.NewParser(cronParseOption)

	errZeroDuration = errors.New("duration should be positive")
)

// CronRange consists of cron expression along with time zone and duration info.
type CronRange struct {
	cronExpression string
	timeZone       string
	duration       time.Duration
	schedule       cron.Schedule
}

// TimeRange represents a time range between starting time and ending time.
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// New returns a CronRange instance with given config, time zone can be empty for local time zone.
//
// It returns an error if duration is not positive number, or cron expression is invalid, or time zone doesn't exist.
func New(cronExpr, timeZone string, durationMin uint64) (cr *CronRange, err error) {
	// Precondition check
	if durationMin == 0 {
		err = errZeroDuration
		return
	}

	// Clean up string parameters
	cronExpr, timeZone = strings.TrimSpace(cronExpr), strings.TrimSpace(timeZone)

	// Append time zone into cron spec if necessary
	cronSpec := cronExpr
	if strings.ToLower(timeZone) == "local" {
		timeZone = ""
	} else if len(timeZone) > 0 {
		cronSpec = fmt.Sprintf("CRON_TZ=%s %s", timeZone, cronExpr)
	}

	// Validate & retrieve crontab schedule
	var schedule cron.Schedule
	if schedule, err = cronParser.Parse(cronSpec); err != nil {
		return
	}

	cr = &CronRange{
		cronExpression: cronExpr,
		timeZone:       timeZone,
		duration:       time.Minute * time.Duration(durationMin),
		schedule:       schedule,
	}
	return
}

// Duration returns the duration of the CronRange.
func (cr *CronRange) Duration() time.Duration {
	cr.checkPrecondition()
	return cr.duration
}

// TimeZone returns the time zone string of the CronRange.
func (cr *CronRange) TimeZone() string {
	cr.checkPrecondition()
	return cr.timeZone
}

// CronExpression returns the Cron expression of the CronRange.
func (cr *CronRange) CronExpression() string {
	cr.checkPrecondition()
	return cr.cronExpression
}
