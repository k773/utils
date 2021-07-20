package utils

import (
	"github.com/k773/utils"
	"sync"
	"time"
)

type DelayedProcess struct {
	Name      string
	Interval  time.Duration
	Callback  func() error
	SyncGroup int
	IsRunning bool
	LastRun   time.Time
}

// DelayedProcessesPerformer : Interval may be 0 -> action will only be performed once
func DelayedProcessesPerformer(delayedProcesses ...DelayedProcess) {
	logger := utils.Logger{
		LogLevel:   utils.LevelDebug,
		LoggerName: "delayedProcessesPerformer",
	}
	var s sync.Mutex
	syncGroupToProcesses := make(map[int]map[int]DelayedProcess)
	for _, process := range delayedProcesses {
		if syncGroupToProcesses[process.SyncGroup] == nil {
			syncGroupToProcesses[process.SyncGroup] = map[int]DelayedProcess{}
		}
		syncGroupToProcesses[process.SyncGroup][len(syncGroupToProcesses[process.SyncGroup])] = process
	}

	runAll := func(firstRun bool) {
		s.Lock()
		defer s.Unlock()

	a:
		for group, processes := range syncGroupToProcesses {
			// All sync group processes must be finished to re-run them
			for _, process := range processes {
				if process.IsRunning {
					continue a
				}
			}

			go func(group int, processes map[int]DelayedProcess) {
				for i, process := range processes {
					if (process.Interval != 0 || firstRun) && !process.IsRunning && time.Now().After(process.LastRun.Add(process.Interval)) {
						s.Lock()
						process.IsRunning = true
						syncGroupToProcesses[group][i] = process
						s.Unlock()

						//fmt.Println(process)
						if err := process.Callback(); err != nil {
							logger.Error("Error while performing process '", process.Name, "':", err)
						}
						s.Lock()
						process.IsRunning = false
						process.LastRun = time.Now()
						syncGroupToProcesses[group][i] = process
						s.Unlock()
					}
				}
			}(group, processes)
		}
	}

	runAll(true)
	for range time.NewTicker(time.Second).C {
		runAll(false)
	}
}
