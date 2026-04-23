package amongus

import (
	"errors"
	"math/rand"
	"slices"
	"strings"
	"sync"
	"time"
)

type Game struct {
	mu *sync.Mutex

	SolvedCodes int
	Stations    []Station `json:"-"`
	Players     []Player  `json:"-"`

	Voting *VotingState

	Settings  GameSettings
	GameStart time.Time
}

type GameSettings struct {
	RequiredCodes    int
	TotalStations    int
	TotalPlayers     int
	CodeMask         string
	CooldownDuration time.Duration `json:",format:sec"`
	GameDuration     time.Duration `json:",format:sec"`
}

type VotingState struct {
	Votes []int
}

type Player struct {
	IsDead           bool
	FinishedStations []int
}

type Station struct {
	CurrentCode   string
	CooldownUntil time.Time
}

func NewGame(settings GameSettings) *Game {
	g := &Game{
		mu: &sync.Mutex{},

		Settings:  settings,
		GameStart: time.Now(),
		Players:   make([]Player, settings.TotalPlayers),
	}
	g.InitializeStations()
	return g
}

var chars = "ABCDEF01234567890"

func (g *Game) SubmitCode(playerId int, code string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Voting != nil {
		return errors.New("voting is underway")
	}

	player := g.Players[playerId]

	for i, station := range g.Stations {
		if station.CurrentCode != code {
			continue
		}
		if slices.Contains(player.FinishedStations, i) {
			return errors.New("already finished this station")
		}
		if station.CooldownUntil.After(time.Now()) {
			return errors.New("station is under cooldown")
		}

		g.SolvedCodes++
		g.Stations[i].CurrentCode = generateCode(g.Settings.CodeMask)
		g.Stations[i].CooldownUntil = time.Now().Add(g.Settings.CooldownDuration)
		player.FinishedStations = append(player.FinishedStations, i)
		g.Players[playerId] = player

		return nil
	}

	return errors.New("incorrect code")
}

func (g *Game) Vote(playerId int, targetPlayerId int) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Voting == nil {
		return errors.New("voting isn't active")
	}

	player := g.Players[playerId]
	if player.IsDead {
		return errors.New("the dead cannot vote")
	}

	if g.Voting.Votes[playerId] != 0 {
		return errors.New("already voted")
	}

	g.Voting.Votes[playerId] = targetPlayerId

	return nil
}

func (g *Game) StartVoting() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.Settings.GameDuration -= time.Since(g.GameStart)
	g.Voting = &VotingState{
		Votes: make([]int, g.Settings.TotalPlayers+1),
	}
	g.GameStart = time.Now()
	return nil
}

func (g *Game) EndVoting() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.Voting = nil
	g.GameStart = time.Now()
	return nil
}

func (g *Game) GetPlayerDeathBitField() []byte {
	result := make([]byte, g.Settings.TotalPlayers/8+1)
	for i, p := range g.Players {
		byteIdx := i / 8
		bytePos := i % 8
		if p.IsDead {
			result[byteIdx] = result[byteIdx] | 1<<bytePos
		}
	}
	return result
}

func (g *Game) TogglePlayer(playerId int) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.Players[playerId].IsDead = !g.Players[playerId].IsDead
	return nil
}

func generateCode(mask string) string {
	code := mask
	for strings.Contains(code, "X") {
		c := chars[rand.Intn(len(chars))]
		code = strings.Replace(code, "X", string(c), 1)
	}
	return code
}

func (g *Game) InitializeStations() {
	g.Stations = make([]Station, g.Settings.TotalStations+1)
	for i := range g.Settings.TotalStations + 1 {
		g.Stations[i].CurrentCode = generateCode(g.Settings.CodeMask)
		g.Stations[i].CooldownUntil = time.Now()
	}
}
