package features

import (
	"sync"
)

var (
	schedulerMutex = &sync.Mutex{}
	jobLock        = &sync.Mutex{}
)

func IsJobRunning() bool {
	jobLocked := jobLock.TryLock()
	if jobLocked {
		jobLock.Unlock()
	}
	return !jobLocked
}
