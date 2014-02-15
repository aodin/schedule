package schedule

import (
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
	cupid := time.Date(2014, time.Month(2), 14, 0, 0, 0, 0, time.UTC)
	ending := `filename.md`
	newEnding := TimestampFilename(cupid, ending)
	expectString(t, newEnding, `filename_2014-02-14.md`)
}
