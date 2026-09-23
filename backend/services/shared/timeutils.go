package shared

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func NormalizeTimestamp(ts *timestamppb.Timestamp) time.Time {
	if ts != nil {
		return ts.AsTime().UTC().Truncate(time.Microsecond)
	}
	return time.Now().UTC().Truncate(time.Microsecond)
}

func NormalizeTime(t time.Time) time.Time {
	return t.UTC().Truncate(time.Microsecond)
}
