package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gonzalohonorato/servercorego/core/websocket/infrastructure"
	"github.com/gorilla/mux"
)

func WebSocketRoutes(router *mux.Router, wsService *infrastructure.WebSocketService) {
	
	router.HandleFunc("/ws/admin", wsService.HandleAdminWebSocket).Methods("GET")

	
	router.HandleFunc("/ws/notifications/{userId}", wsService.HandleUserNotificationWebSocket).Methods("GET")

	
	router.HandleFunc("/ws/stats", func(w http.ResponseWriter, r *http.Request) {
		stats := wsService.GetConnectionStats()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}).Methods("GET")

	
	router.HandleFunc("/ws/connected-users", func(w http.ResponseWriter, r *http.Request) {
		users := wsService.GetConnectedUsers()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"connectedUsers":   users,
			"totalConnections": len(users),
		})
	}).Methods("GET")
}
