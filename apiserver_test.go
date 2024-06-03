package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestApiServer_GetGame(t *testing.T) {
	apiServer := NewApiServer()

	t.Run("Should return StatusNotFound if no game is found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/game/1", nil)

		w := httptest.NewRecorder()

		apiServer.handleGameRequest(w, req)

		res := w.Result()

		defer res.Body.Close()

		if res.StatusCode != http.StatusNotFound {
			t.Errorf("Expected bad request code, got %d", res.StatusCode)
		}
	})

	t.Run("Should return game if exists", func(t *testing.T) {
		gameData := GameData{
			Instance: map[string]string{"foo": "bar"},
			Actions:  []string{"action1", "action2"},
		}
		games["1"] = *NewGame("1", gameData)

		req := httptest.NewRequest(http.MethodGet, "/game/1", nil)
		w := httptest.NewRecorder()

		apiServer.handleGameRequest(w, req)

		res := w.Result()

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status ok, got %d", res.StatusCode)
		}

		var got GameData
		json.NewDecoder(res.Body).Decode(&got)

		if got.Instance == nil {
			t.Errorf("expected response instance not to be nil")
		}

		gotStr := fmt.Sprintf("%v", got.Instance)
		wantStr := fmt.Sprintf("%v", got.Instance)

		if gotStr != wantStr {
			t.Errorf("Expected instance to equal %v got %v", got, gameData)
		}

	})
}

func TestApiServer_UpdateGame(t *testing.T) {
	apiServer := NewApiServer()

	t.Run("should create game if it doesn't exist", func(t *testing.T) {
		postBody := map[string]interface{}{
			"instance": "This is test data",
		}
		body, _ := json.Marshal(postBody)

		req := httptest.NewRequest(http.MethodPost, "/game/2", bytes.NewReader(body))
		w := httptest.NewRecorder()

		apiServer.updateGame(w, req)

		res := w.Result()

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status ok, got %d", res.StatusCode)
		}

		game := games["2"]

		if game.id != "2" {
			t.Errorf("Expected game id to be set")
		}

		oneSecondAgo := time.Now().Add(time.Duration(-1) * time.Minute)

		if game.created.Before(oneSecondAgo) {
			t.Errorf("Expected game created to be now, got %v", game.created)
		}

		if game.created.After(time.Now()) {
			t.Errorf("Game created time is in the future: %v", game.created)
		}

		wantGameData := GameData{
			Instance: "This is test data",
		}

		gotStr := fmt.Sprintf("%v", game.Data)
		wantStr := fmt.Sprintf("%v", wantGameData)

		if gotStr != wantStr {
			t.Errorf("Expected game data to have correct data, got %s want %s", gotStr, wantStr)
		}
	})

	t.Run("should update game", func(t *testing.T) {
		gameData := GameData{
			Instance: map[string]string{"foo": "bar"},
			Actions:  []string{"action1", "action2"},
		}
		existingGame := *NewGame("3", gameData)
		existingGame.created = time.Now().Add(time.Duration(-1) * time.Hour)
		games["3"] = existingGame

		postBody := map[string]interface{}{
			"test": "This is updated data",
		}
		body, _ := json.Marshal(postBody)

		req := httptest.NewRequest(http.MethodPost, "/game/3", bytes.NewReader(body))
		w := httptest.NewRecorder()

		apiServer.updateGame(w, req)

		res := w.Result()

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status ok, got %d", res.StatusCode)
		}

		createdGame := games["3"]

		gotStr := fmt.Sprintf("%v", createdGame.Data)
		wantStr := fmt.Sprintf("%v", gameData)

		if gotStr != wantStr {
			t.Errorf("Expected game data to have correct data, got %s want %s", gotStr, wantStr)
		}

		if !createdGame.created.Equal(existingGame.created) {
			t.Errorf("Expected game created date not to have been updated for existing game, got %v", createdGame.created)
		}
	})
}
