package jobs

import (
	"strconv"
	"time"

	rtypes "github.com/pablonlr/poly-crown-relayer/types"
)

func WaitTime(params ...string) rtypes.TaskResult {
	durationSeconds, err := strconv.Atoi(params[0])
	if err != nil {
		return rtypes.TaskResult{
			Err: rtypes.GetError(rtypes.InvalidWaitTime, err),
		}
	}
	time.Sleep(time.Duration(durationSeconds) * time.Second)
	return rtypes.TaskResult{
		ResultValue: "Waited " + params[0] + " seconds",
	}
}
