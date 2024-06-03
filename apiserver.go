package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type GameData struct {
	Instance interface{} `json:"instance"`
	Actions  interface{} `json:"actions"`
}

type Game struct {
	id      string
	created time.Time
	Data    GameData `json:"data"`
}

func NewGame(id string, data GameData) *Game {
	return &Game{
		id:      id,
		created: time.Now(),
		Data:    data,
	}
}

var games = make(map[string]Game)

type ApiServer struct{}

func NewApiServer() *ApiServer {
	return &ApiServer{}
}

func (s *ApiServer) handleGameRequest(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		s.getGame(w, req)
	case http.MethodPost:
	case http.MethodPatch:
		s.updateGame(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *ApiServer) updateGame(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/game/")
	decoder := json.NewDecoder(req.Body)

	var data GameData
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Printf("Failed to decode game")
	}

	game, ok := games[id]
	if ok {
		if data.Instance != nil {
			game.Data.Instance = data.Instance
		}
		if data.Actions != nil {
			game.Data.Actions = data.Actions
		}
		games[id] = game
	} else {
		games[id] = *NewGame(id, data)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *ApiServer) getGame(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/game/")

	w.Header().Set("Content-Type", "application/json")

	game, ok := games[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	j, err := json.Marshal(game.Data)
	if err != nil {
		fmt.Printf("failed to marshal data")
	}
	w.Write(j)
}

func (s *ApiServer) ListenAndServe(addr string) {
	http.HandleFunc("/game/{id}", s.handleGameRequest)

	http.ListenAndServe(":8090", nil)
}
