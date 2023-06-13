package features

import (
	"sync"
)

var (
	schedulerMutex = &sync.Mutex{}
	jobLock        = &sync.Mutex{}
)
