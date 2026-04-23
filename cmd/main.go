package main

import (
	"crypto/subtle"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jkbcz/amongus"
	"github.com/rs/cors"
)

var (
	password  = os.Getenv("PASSWORD")
	staticDir = os.Getenv("STATIC_DIR")
)

func main() {
	mux := http.NewServeMux()
	g := amongus.NewGame(100, "XXXX-XXXX-XXXX-XXXX", 10, time.Second*10, time.Minute*30)

	mux.HandleFunc("POST /api/submit", func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		playerId := r.FormValue("player_id")
		err := g.SubmitCode(playerId, code)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, err.Error())
		}
	})

	mux.HandleFunc("GET /api/player_state", func(w http.ResponseWriter, r *http.Request) {
		playerId := r.FormValue("player_id")
		p, _ := g.Players[playerId]
		result := amongus.PlayerState{
			FinishedStations: p.FinishedStations,
			CodeMask:         g.CodeMask,
			TotalStations:    g.TotalStations,
			IsVoting:         g.IsVoting,
		}
		json.MarshalWrite(w, result)
	})

	mux.HandleFunc("GET /api/station_state", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		station_id, _ := strconv.Atoi(r.FormValue("station_id"))
		if station_id < 1 || station_id >= len(g.Stations) {
			http.Error(w, "invalid station_id", http.StatusBadRequest)
			return
		}
		station := g.Stations[station_id]
		result := amongus.StationState{
			CurrentCode:      station.CurrentCode,
			CooldownUntil:    station.CooldownUntil,
			CooldownDuration: g.CooldownDuration,
			IsVoting:         g.IsVoting,
		}
		json.MarshalWrite(w, result)
	}))

	mux.HandleFunc("GET /api/result_state", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		result := amongus.ViewState{
			SolvedTasks:  g.SolvedCodes,
			TotalTasks:   g.RequiredCodes,
			GameStart:    g.GameStart,
			GameDuration: g.GameDuration,
			IsVoting:     g.IsVoting,
		}
		json.MarshalWrite(w, result)
	}))

	mux.HandleFunc("GET /api/admin_state", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		json.MarshalWrite(w, g)
	}))

	mux.HandleFunc("POST /api/settings", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		reset := r.FormValue("reset")
		if reset == "true" {
			g = amongus.NewGame(g.RequiredCodes, g.CodeMask, g.TotalStations, g.CooldownDuration, g.GameDuration)
		}

		isVoting := r.FormValue("is_voting")
		if isVoting != "" {
			g.IsVoting = isVoting == "true"
			if g.IsVoting {
				g.GameDuration -= time.Since(g.GameStart)
			}
			g.GameStart = time.Now()
		}

		requiredCodes, err := strconv.Atoi(r.FormValue("required_codes"))
		if err == nil {
			g.RequiredCodes = requiredCodes
		}
		stations, err := strconv.Atoi(r.FormValue("total_stations"))
		if err == nil {
			g.TotalStations = stations
			g.InitializeStations()
		}
		codeMask := r.FormValue("code_mask")
		if codeMask != "" {
			g.CodeMask = codeMask
			g.InitializeStations()
		}
		cooldown, err := time.ParseDuration(r.FormValue("cooldown_duration"))
		if err == nil {
			g.CooldownDuration = cooldown
		}
		gameDuration, err := time.ParseDuration(r.FormValue("game_duration"))
		if err == nil {
			g.GameDuration = gameDuration
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

	http.ListenAndServe(":8080", cors.AllowAll().Handler(mux))
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := r.Header.Get("X-Pass")

		if subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
