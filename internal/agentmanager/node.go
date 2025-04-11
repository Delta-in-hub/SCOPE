package agentmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"scope/internal/models"
	"scope/internal/utils"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var Token string

// RegisterNodeToCenter registers this node to the center node
func RegisterNodeToCenter(centerURL string) (string, error) {
	// Get machine ID
	machineID := getMachineID()

	// Get IP addresses
	ips := utils.GetMyIpAddrs()

	// Create agent info
	agentInfo := models.NodeInfo{
		ID:       machineID,
		IPs:      ips,
		LastSeen: time.Now(),
		Status:   "online",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(agentInfo)
	if err != nil {
		return "", fmt.Errorf("failed to marshal agent info: %w", err)
	}

	// 	r.Route("/api/v1/node", func(r chi.Router) {
	// 		r.Post("/up", handler.nodeHandler.NodeUp)
	// 		r.Post("/down", handler.nodeHandler.NodeDown)

	// Ensure the URL has the correct format
	if !strings.HasPrefix(centerURL, "http://") && !strings.HasPrefix(centerURL, "https://") {
		centerURL = "http://" + centerURL
	}

	// Ensure the URL has the register endpoint
	if !strings.HasSuffix(centerURL, "/api/v1/node/up") {
		if strings.HasSuffix(centerURL, "/") {
			centerURL = centerURL + "api/v1/node/up"
		} else {
			centerURL = centerURL + "/api/v1/node/up"
		}
	}

	// Send registration request
	resp, err := http.Post(centerURL, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return "", fmt.Errorf("failed to register with center: %w", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("center returned non-OK status: %s", resp.Status)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	Token = response["token"]
	return Token, nil
}

func SetupRouter() *chi.Mux {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/health"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Welcome message
		w.Write([]byte("Welcome to Scope Node Agent Manager\nSee More at https://github.com/Delta-in-hub/ebpf-golang\n"))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"timestamp": time.Now().String(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	})

	r.Post("/runEBPF", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println(data)
		if data["token"] != Token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// app := data["app"]
		// args := data["args"]
		w.WriteHeader(http.StatusOK)
	})

	return r
}
