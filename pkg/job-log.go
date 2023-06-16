package pkg

import (
	"sync"
	"time"
)

type JobLog struct {
	JobName string `json:"jobName" example:"upload movie"`
	Date    string `json:"date" example:"2021-01-01 12:00:00"`
	Message string `json:"message" example:"Uploading movie test.mp4"`
}

func AppendJobLog(message string) {
	jobLogsLock.Lock()
	defer jobLogsLock.Unlock()
	jobLogs = append(jobLogs, JobLog{
		JobName: jobName,
		Message: message,
		Date:    time.Now().Format("2006-01-02 15:04:05"),
	})
}

func GetJobLogs() []JobLog {
	jobLogsLock.Lock()
	defer jobLogsLock.Unlock()
	return jobLogs
}

func PopJobLogs() []JobLog {
	jobLogsLock.Lock()
	defer jobLogsLock.Unlock()
	logs := jobLogs
	jobLogs = make([]JobLog, 0)
	return logs
}

func ClearJobLogs(newJobName string) {
	jobLogsLock.Lock()
	defer jobLogsLock.Unlock()
	jobLogs = make([]JobLog, 0)
	jobName = newJobName
}

func GetJobName() string {
	jobLogsLock.Lock()
	defer jobLogsLock.Unlock()
	return jobName
}

var (
	jobName     = ""
	jobLogsLock = &sync.Mutex{}
	jobLogs     = make([]JobLog, 0)
)
