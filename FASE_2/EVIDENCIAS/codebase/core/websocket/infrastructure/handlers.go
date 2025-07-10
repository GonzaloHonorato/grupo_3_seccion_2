package infrastructure

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)


func (s *WebSocketService) HandleAdminWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("❌ Error al actualizar conexión admin WebSocket: %v", err)
		return
	}

	client := &WebSocketClient{
		Conn:     conn,
		Type:     AdminConnection,
		LastPing: time.Now(),
		Send:     make(chan []byte, 256),
	}

	s.register <- client

	
	go s.handleClientWrite(client)

	
	go s.handleClientRead(client)
}


func (s *WebSocketService) HandleUserNotificationWebSocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	if userID == "" {
		log.Printf("❌ UserID requerido para conexión de notificaciones")
		http.Error(w, "UserID requerido", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("❌ Error al actualizar conexión usuario WebSocket: %v", err)
		return
	}

	client := &WebSocketClient{
		Conn:     conn,
		Type:     UserNotificationConnection,
		UserID:   userID,
		LastPing: time.Now(),
		Send:     make(chan []byte, 256),
	}

	s.register <- client

	
	go s.handleClientWrite(client)

	
	go s.handleUserClientRead(client)
}


func (s *WebSocketService) handleClientWrite(client *WebSocketClient) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("❌ Error al escribir mensaje: %v", err)
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}


func (s *WebSocketService) handleClientRead(client *WebSocketClient) {
	defer func() {
		s.unregister <- client
		client.Conn.Close()
	}()

	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		client.LastPing = time.Now()
		return nil
	})

	for {
		_, _, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("❌ Error de conexión WebSocket admin: %v", err)
			}
			break
		}
		client.LastPing = time.Now()
	}
}


func (s *WebSocketService) handleUserClientRead(client *WebSocketClient) {
	defer func() {
		s.unregister <- client
		client.Conn.Close()
	}()

	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		client.LastPing = time.Now()
		return nil
	})

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("❌ Error de conexión WebSocket usuario %s: %v", client.UserID, err)
			}
			break
		}

		client.LastPing = time.Now()

		
		s.handleUserMessage(client, message)
	}
}


func (s *WebSocketService) handleUserMessage(client *WebSocketClient, message []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("❌ Error al parsear mensaje de usuario %s: %v", client.UserID, err)
		return
	}

	msgType, ok := msg["type"].(string)
	if !ok {
		log.Printf("⚠️ Mensaje sin tipo de usuario %s", client.UserID)
		return
	}

	switch msgType {
	case "pong":
		

	case "mark_as_read":
		
		
		if notificationID, exists := msg["notificationId"]; exists {
			log.Printf("✅ Usuario %s marcó notificación %v como leída", client.UserID, notificationID)
			
		}

	case "get_user_notifications":
		
		log.Printf("📥 Usuario %s solicitó sus notificaciones", client.UserID)
		

	case "auth":
		
	default:
		log.Printf("❓ Tipo de mensaje desconocido '%s' de usuario %s", msgType, client.UserID)
	}
}
