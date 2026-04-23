package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json/v2"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jkbcz/amongus"
	"github.com/rs/cors"
)

type UserHandler func(w http.ResponseWriter, r *http.Request, playerId int)

const (
	userPassLength = 20
)

var (
	password  = os.Getenv("PASSWORD")
	staticDir = os.Getenv("STATIC_DIR")
)

func main() {
	mux := http.NewServeMux()

	g := amongus.NewGame(amongus.GameSettings{
		RequiredCodes:    100,
		TotalStations:    10,
		TotalPlayers:     100,
		CodeMask:         "XXXX-XXXX-XXXX-XXXX",
		CooldownDuration: time.Second * 10,
		GameDuration:     time.Minute * 30,
	})

	mux.HandleFunc("POST /api/submit", UserAuthMiddleware(func(w http.ResponseWriter, r *http.Request, playerId int) {
		code := r.FormValue("code")
		err := g.SubmitCode(playerId, code)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, err.Error())
		}
	}))

	mux.HandleFunc("POST /api/vote", UserAuthMiddleware(func(w http.ResponseWriter, r *http.Request, playerId int) {
		targetPlayerId, _ := strconv.Atoi(r.FormValue("target_player_id"))
		err := g.Vote(playerId, targetPlayerId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, err.Error())
		}
	}))

	mux.HandleFunc("GET /api/player_state", UserAuthMiddleware(func(w http.ResponseWriter, r *http.Request, playerId int) {
		p := g.Players[playerId]
		result := amongus.PlayerState{
			PlayerId:         playerId,
			FinishedStations: p.FinishedStations,
			CodeMask:         g.Settings.CodeMask,
			TotalStations:    g.Settings.TotalStations,
			IsDead:           p.IsDead,
		}
		if g.Voting != nil {
			result.VotingState = &amongus.PlayerVotingState{
				MyVote:      g.Voting.Votes[playerId],
				PlayersDead: g.GetPlayerDeathBitField(),
			}
		}
		json.MarshalWrite(w, result)
	}))

	mux.HandleFunc("GET /api/station_state", AdminAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		station_id, _ := strconv.Atoi(r.FormValue("station_id"))
		if station_id < 1 || station_id >= len(g.Stations) {
			http.Error(w, "invalid station_id", http.StatusBadRequest)
			return
		}
		station := g.Stations[station_id]
		result := amongus.StationState{
			CurrentCode:      station.CurrentCode,
			CooldownUntil:    station.CooldownUntil,
			CooldownDuration: g.Settings.CooldownDuration,
			IsVoting:         g.Voting != nil,
		}
		json.MarshalWrite(w, result)
	}))

	mux.HandleFunc("GET /api/result_state", AdminAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		result := amongus.ViewState{
			SolvedTasks:  g.SolvedCodes,
			TotalTasks:   g.Settings.RequiredCodes,
			GameStart:    g.GameStart,
			GameDuration: g.Settings.GameDuration,
			IsVoting:     g.Voting != nil,
		}
		json.MarshalWrite(w, result)
	}))

	mux.HandleFunc("GET /api/admin_state", AdminAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		json.MarshalWrite(w, amongus.AdminState{
			IsVoting: g.Voting != nil,
			Settings: g.Settings,
			Players:  g.GetPlayerDeathBitField(),
		})
	}))

	mux.HandleFunc("POST /api/toggle_player", AdminAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		playerId, _ := strconv.Atoi(r.FormValue("player_id"))
		g.TogglePlayer(playerId)
	}))

	mux.HandleFunc("POST /api/settings", AdminAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		reset := r.FormValue("reset")
		if reset == "true" {
			g = amongus.NewGame(g.Settings)
		}

		isVoting := r.FormValue("is_voting")
		if isVoting != "" {
			if isVoting == "true" {
				g.StartVoting()
			} else {
				g.EndVoting()
			}
		}

		requiredCodes, err := strconv.Atoi(r.FormValue("required_codes"))
		if err == nil {
			g.Settings.RequiredCodes = requiredCodes
		}
		stations, err := strconv.Atoi(r.FormValue("total_stations"))
		if err == nil {
			g.Settings.TotalStations = stations
			g.InitializeStations()
		}
		players, err := strconv.Atoi(r.FormValue("total_players"))
		if err == nil {
			g.Settings.TotalPlayers = players
			g.Players = g.Players[:g.Settings.TotalPlayers+1]
		}
		codeMask := r.FormValue("code_mask")
		if codeMask != "" {
			g.Settings.CodeMask = codeMask
			g.InitializeStations()
		}
		cooldown, err := time.ParseDuration(r.FormValue("cooldown_duration"))
		if err == nil {
			g.Settings.CooldownDuration = cooldown
		}
		gameDuration, err := time.ParseDuration(r.FormValue("game_duration"))
		if err == nil {
			g.Settings.GameDuration = gameDuration
			g.GameStart = time.Now()
		}
	}))

	if staticDir != "" {
		mux.Handle("/assets/", http.FileServer(http.Dir(staticDir)))
		mux.Handle("/favicon.ico", http.FileServer(http.Dir(staticDir)))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
		})
	}

	if err := http.ListenAndServe(":8080", cors.AllowAll().Handler(mux)); err != nil {
		log.Fatal(err)
	}
}

func AdminAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := r.Header.Get("X-Pass")

		if subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func UserAuthMiddleware(next UserHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := r.Header.Get("X-Pass")
		if len(pass) != userPassLength {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		rawPlayerId := pass[:4]
		playerId, err := strconv.Atoi(pass[:4])
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		hash, err := hex.DecodeString(pass[4:])
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		hmac := hmac.New(sha256.New, []byte(password))
		hmac.Write([]byte(rawPlayerId))
		expectedHash := hmac.Sum(nil)
		if subtle.ConstantTimeCompare(hash, expectedHash[:len(hash)]) != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r, playerId)
	}
}
