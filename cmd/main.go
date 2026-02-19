package main

import (
	"crypto/subtle"
	"encoding/json/v2"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/jkbcz/amongus"
	"github.com/rs/cors"
)

const (
	password = "secret-password"
)

func main() {
	mux := http.NewServeMux()
	g := amongus.NewGame(100, "XXXX-XXXX-XXXX-XXXX", 10, time.Second*10)

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
			CooldownDuration: int(math.Round(g.CooldownDuration.Seconds())),
			IsVoting:         g.IsVoting,
		}
		json.MarshalWrite(w, result)
	}))

	mux.HandleFunc("GET /api/result_state", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		result := amongus.ResultState{
			RemainingTasks: g.RequiredCodes - g.SolvedCodes,
			IsVoting:       g.IsVoting,
		}
		json.MarshalWrite(w, result)
	}))

	mux.HandleFunc("GET /api/admin_state", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		json.MarshalWrite(w, g)
	}))

	mux.HandleFunc("POST /api/settings", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		reset := r.FormValue("reset")
		if reset == "true" {
			g = amongus.NewGame(g.RequiredCodes, g.CodeMask, g.TotalStations, g.CooldownDuration)
		}

		isVoting := r.FormValue("is_voting")
		if isVoting != "" {
			g.IsVoting = isVoting == "true"
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
	}))

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
