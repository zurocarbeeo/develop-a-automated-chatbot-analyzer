package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ChatbotAnalyzer struct {
	ChatbotID   string `json:"chatbot_id"`
	Conversation []struct {
		UserInput  string `json:"user_input"`
		BotResponse string `json:"bot_response"`
	} `json:"conversation"`
}

type API struct {
_analyzer *ChatbotAnalyzer
}

func NewAPI() *API {
	return &API{_analyzer: &ChatbotAnalyzer{}}
}

func (a *API) AnalyzeConversation(w http.ResponseWriter, r *http.Request) {
	var conversation struct {
		UserInput  string `json:"user_input"`
		BotResponse string `json:"bot_response"`
	}
	err := json.NewDecoder(r.Body).Decode(&conversation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	a._analyzer.Conversation = append(a._analyzer.Conversation, conversation)
	w.WriteHeader(http.StatusCreated)
}

func (a *API) GetConversationAnalysis(w http.ResponseWriter, r *http.Request) {
	analysis := a.analyzeConversation(a._analyzer.Conversation)
	json.NewEncoder(w).Encode(analysis)
}

func (a *API) analyzeConversation(conversation []struct {
	UserInput  string `json:"user_input"`
	BotResponse string `json:"bot_response"`
}) map[string]int {
	analysis := make(map[string]int)
	for _, msg := range conversation {
		if msg.BotResponse != "" {
			analysis["successful_responses"]++
		} else {
			analysis["failed_responses"]++
		}
	}
	return analysis
}

func main() {
	r := mux.NewRouter()
	api := NewAPI()
	r.HandleFunc("/analyze", api.AnalyzeConversation).Methods("POST")
	r.HandleFunc("/analysis", api.GetConversationAnalysis).Methods("GET")
	fmt.Println("Starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}