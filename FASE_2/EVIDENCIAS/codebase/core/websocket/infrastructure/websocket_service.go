package infrastructure

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}


type ConnectionType string

const (
	AdminConnection            ConnectionType = "admin"
	UserNotificationConnection ConnectionType = "user_notification"
)


type WebSocketClient struct {
	Conn     *websocket.Conn
	Type     ConnectionType
	UserID   string 
	LastPing time.Time
	Send     chan []byte
}

type WebSocketService struct {
	
	adminClients map[*WebSocketClient]bool

	
	userClients map[string]*WebSocketClient 

	
	broadcast     chan []byte
	userBroadcast chan UserMessage
	register      chan *WebSocketClient
	unregister    chan *WebSocketClient

	mutex sync.RWMutex
}


type UserMessage struct {
	UserID  string      `json:"userId"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}


type WebSocketMessage struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp string      `json:"timestamp"`
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		adminClients:  make(map[*WebSocketClient]bool),
		userClients:   make(map[string]*WebSocketClient),
		broadcast:     make(chan []byte),
		userBroadcast: make(chan UserMessage),
		register:      make(chan *WebSocketClient),
		unregister:    make(chan *WebSocketClient),
	}
}

func (s *WebSocketService) Start() {
	log.Println("🚀 WebSocket Service iniciado")

	
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case client := <-s.register:
			s.handleRegister(client)

		case client := <-s.unregister:
			s.handleUnregister(client)

		case message := <-s.broadcast:
			s.handleAdminBroadcast(message)

		case userMsg := <-s.userBroadcast:
			s.handleUserBroadcast(userMsg)

		case <-ticker.C:
			s.pingClients()
		}
	}
}

func (s *WebSocketService) handleRegister(client *WebSocketClient) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	switch client.Type {
	case AdminConnection:
		s.adminClients[client] = true
		log.Printf("👨‍💼 Administrador conectado. Total admins: %d", len(s.adminClients))

	case UserNotificationConnection:
		
		if existingClient, exists := s.userClients[client.UserID]; exists {
			log.Printf("🔄 Usuario %s reconectándose, cerrando conexión anterior", client.UserID)
			existingClient.Conn.Close()
		}

		s.userClients[client.UserID] = client
		log.Printf("👤 Usuario %s conectado para notificaciones. Total usuarios: %d",
			client.UserID, len(s.userClients))

		
		welcomeMsg := WebSocketMessage{
			Type: "welcome",
			Payload: map[string]interface{}{
				"message": "Conectado al sistema de notificaciones",
				"userId":  client.UserID,
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		s.sendToClient(client, welcomeMsg)
	}

	client.LastPing = time.Now()
}

func (s *WebSocketService) handleUnregister(client *WebSocketClient) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	switch client.Type {
	case AdminConnection:
		if _, ok := s.adminClients[client]; ok {
			delete(s.adminClients, client)
			close(client.Send)
			log.Printf("👨‍💼 Administrador desconectado. Total admins: %d", len(s.adminClients))
		}

	case UserNotificationConnection:
		if existingClient, exists := s.userClients[client.UserID]; exists && existingClient == client {
			delete(s.userClients, client.UserID)
			close(client.Send)
			log.Printf("👤 Usuario %s desconectado. Total usuarios: %d",
				client.UserID, len(s.userClients))
		}
	}
}

func (s *WebSocketService) handleAdminBroadcast(message []byte) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for client := range s.adminClients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(s.adminClients, client)
		}
	}
}

func (s *WebSocketService) handleUserBroadcast(userMsg UserMessage) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if client, exists := s.userClients[userMsg.UserID]; exists {
		message := WebSocketMessage{
			Type:      userMsg.Type,
			Payload:   userMsg.Payload,
			Timestamp: time.Now().Format(time.RFC3339),
		}

		s.sendToClient(client, message)
	} else {
		log.Printf("⚠️ Intento de envío a usuario %s que no está conectado", userMsg.UserID)
	}
}

func (s *WebSocketService) sendToClient(client *WebSocketClient, message interface{}) {
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("❌ Error al serializar mensaje: %v", err)
		return
	}

	select {
	case client.Send <- jsonData:
	default:
		log.Printf("⚠️ Canal de envío bloqueado para cliente")
	}
}

func (s *WebSocketService) pingClients() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	now := time.Now()
	pingMessage := WebSocketMessage{
		Type:      "ping",
		Payload:   map[string]string{"timestamp": now.Format(time.RFC3339)},
		Timestamp: now.Format(time.RFC3339),
	}

	
	for client := range s.adminClients {
		if now.Sub(client.LastPing) > 60*time.Second {
			log.Printf("⏰ Cliente admin sin respuesta, desconectando")
			s.unregister <- client
		} else {
			s.sendToClient(client, pingMessage)
		}
	}

	
	for userID, client := range s.userClients {
		if now.Sub(client.LastPing) > 60*time.Second {
			log.Printf("⏰ Usuario %s sin respuesta, desconectando", userID)
			s.unregister <- client
		} else {
			s.sendToClient(client, pingMessage)
		}
	}
}




func (s *WebSocketService) BroadcastParkingUsage(data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("❌ Error al serializar ParkingUsage: %v", err)
		return
	}

	log.Printf("📡 Broadcasting ParkingUsage a %d administrador(es)", len(s.adminClients))
	s.broadcast <- jsonData
}


func (s *WebSocketService) NotifyUser(userID string, notificationType string, payload interface{}) {
	log.Printf("🔔 Enviando notificación tipo '%s' a usuario %s", notificationType, userID)

	userMsg := UserMessage{
		UserID:  userID,
		Type:    notificationType,
		Payload: payload,
	}

	s.userBroadcast <- userMsg
}


func (s *WebSocketService) NotifyMultipleUsers(userIDs []string, notificationType string, payload interface{}) {
	log.Printf("🔔 Enviando notificación tipo '%s' a %d usuario(s)", notificationType, len(userIDs))

	for _, userID := range userIDs {
		s.NotifyUser(userID, notificationType, payload)
	}
}


func (s *WebSocketService) GetConnectedUsers() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	users := make([]string, 0, len(s.userClients))
	for userID := range s.userClients {
		users = append(users, userID)
	}
	return users
}


func (s *WebSocketService) GetConnectionStats() map[string]int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return map[string]int{
		"admin_clients": len(s.adminClients),
		"user_clients":  len(s.userClients),
		"total_clients": len(s.adminClients) + len(s.userClients),
	}
}
