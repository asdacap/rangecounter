package rangecounter

import (
	"github.com/pkg/errors"
	"time"
)

type DateRange int

const (
	Seconds DateRange = iota
	Minute
	Hour
)

func (drange DateRange) alignDate(at time.Time) (time.Time, error) {
	switch drange {
	case Seconds:
		return time.Date(at.Year(), at.Month(), at.Day(), at.Hour(), at.Minute(), at.Second(), 0, at.Location()), nil
	case Minute:
		return time.Date(at.Year(), at.Month(), at.Day(), at.Hour(), at.Minute(), 0, 0, at.Location()), nil
	case Hour:
		return time.Date(at.Year(), at.Month(), at.Day(), at.Hour(), 0, 0, 0, at.Location()), nil
	}
	return time.Time{}, errors.Errorf("unknown alignment: %v", drange)
}

func (drange DateRange) incrementDate(multiple int, at time.Time) (time.Time, error) {
	switch drange {
	case Seconds:
		return at.Add(time.Second * time.Duration(multiple)), nil
	case Minute:
		return at.Add(time.Minute * time.Duration(multiple)), nil
	case Hour:
		return at.Add(time.Hour * time.Duration(multiple)), nil
	}
	return time.Time{}, errors.Errorf("unknown alignment: %v", drange)
}

func (drange DateRange) String() string {
	switch drange {
	case Seconds:
		return "second"
	case Minute:
		return "minute"
	case Hour:
		return "hour"
	}
	return "unknown range"
}
