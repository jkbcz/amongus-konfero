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
	Stations    []Station         `json:"-"`
	Players     map[string]Player `json:"-"`

	IsVoting bool

	RequiredCodes    int
	TotalStations    int
	CodeMask         string
	CooldownDuration time.Duration `json:",format:sec"`
}

type Player struct {
	FinishedStations []int
}

type Station struct {
	CurrentCode   string
	CooldownUntil time.Time
}

func NewGame(requiredCodes int, codeMask string, totalStations int, cooldownDuration time.Duration) *Game {
	g := &Game{
		mu: &sync.Mutex{},

		RequiredCodes:    requiredCodes,
		CodeMask:         codeMask,
		TotalStations:    totalStations,
		CooldownDuration: cooldownDuration,

		Players: map[string]Player{},
	}
	g.InitializeStations()
	return g
}

var chars = "ABCDEF01234567890"

func (g *Game) SubmitCode(playerId string, code string) error {

	g.mu.Lock()
	defer g.mu.Unlock()

	if g.IsVoting {
		return errors.New("voting is underway")
	}

	player, _ := g.Players[playerId]

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
		g.Stations[i].CurrentCode = generateCode(g.CodeMask)
		g.Stations[i].CooldownUntil = time.Now().Add(g.CooldownDuration)
		player.FinishedStations = append(player.FinishedStations, i)
		g.Players[playerId] = player

		return nil
	}

	return errors.New("incorrect code")
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
	g.Stations = make([]Station, g.TotalStations+1)
	for i := range g.TotalStations + 1 {
		g.Stations[i].CurrentCode = generateCode(g.CodeMask)
		g.Stations[i].CooldownUntil = time.Now()
	}
}
