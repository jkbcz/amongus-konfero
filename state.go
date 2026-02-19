package amongus

import "time"

type PlayerState struct {
	FinishedStations []int
	TotalStations    int
	CodeMask         string
	IsVoting         bool
}

type StationState struct {
	CurrentCode      string
	CooldownUntil    time.Time `json:",format:unix"`
	CooldownDuration int
	IsVoting         bool
}

type ResultState struct {
	RemainingTasks int
	IsVoting       bool
}
