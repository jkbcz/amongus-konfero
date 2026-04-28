package amongus

import "time"

type PlayerState struct {
	PlayerId         int
	FinishedStations []int
	TotalStations    int
	CodeMask         string
	VotingState      *PlayerVotingState
	IsDead           bool
}

type AdminState struct {
	Voting   *VotingState
	Settings GameSettings
	Players  []byte
}

type PlayerVotingState struct {
	Players      []byte
	TotalPlayers int
	MyVote       int
	Active       bool
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

	AlivePlayers int
	TotalPlayers int

	GameStart    time.Time     `json:",format:unix"`
	GameDuration time.Duration `json:",format:sec"`

	Voting *VotingState
}
