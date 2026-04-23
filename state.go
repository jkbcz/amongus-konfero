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
	CooldownUntil    time.Time     `json:",format:unix"`
	CooldownDuration time.Duration `json:",format:sec"`
	IsVoting         bool
}

type ViewState struct {
	SolvedTasks int
	TotalTasks  int

	GameStart    time.Time     `json:",format:unix"`
	GameDuration time.Duration `json:",format:sec"`
	IsVoting     bool
}
