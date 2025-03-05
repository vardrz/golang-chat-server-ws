package handlers

import (
	"bytes"
	"chat-server/models"
	"chat-server/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var (
	clients = make(map[uint]*websocket.Conn)
	mutex   = &sync.Mutex{}
)

type WebSocketHandler struct {
	db         *gorm.DB
	jwtService *services.JWTService
}

func NewWebSocketHandler(db *gorm.DB) *WebSocketHandler {
	return &WebSocketHandler{
		db:         db,
		jwtService: services.NewJWTService(),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatMessage struct {
	ToID    uint   `json:"to_id"`
	Content string `json:"content"`
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	token_catch := c.Query("token")
	if token_catch == "" {
		fmt.Println("error", "No token provided")
		c.JSON(401, gin.H{"error": "No token provided"})
		return
	}

	// Validate JWT
	tokenString := h.jwtService.ExtractToken(token_catch)
	token, err := h.jwtService.ValidateToken(tokenString)

	if err != nil {
		fmt.Println("JWT Parse Error:", err)
		fmt.Println("key:", h.jwtService.GetSecretKey())
		fmt.Println("token:", tokenString)
		c.JSON(401, gin.H{"error": "Invalid token"})
		return
	}

	if !token.Valid {
		fmt.Println("Token is invalid")
		c.JSON(401, gin.H{"error": "Invalid token"})
		return
	}

	userID, err := h.jwtService.ExtractUserID(token)
	if err != nil {
		fmt.Println("Error extracting user ID:", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// Rest of the WebSocket handling code remains the same
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("error", err)
		log.Println(err)
		return
	}
	// Register client
	mutex.Lock()
	clients[userID] = conn
	mutex.Unlock()
	// Load message history
	var messages []models.Message
	h.db.Where("from_id = ? OR to_id = ?", userID, userID).Order("created_at asc").Limit(50).Find(&messages)
	// Send message history to client
	conn.WriteJSON(gin.H{
		"type": "history",
		"messages": messages,
	})
	// Handle incoming messages
	for {
		var msg ChatMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			mutex.Lock()
			delete(clients, userID)
			mutex.Unlock()
			break
		}

		// Save message to database
		message := models.Message{
			FromID:  userID,
			ToID:    msg.ToID,
			Content: msg.Content,
		}
		h.db.Create(&message)

		var senderUser models.User
		var recipientUser models.User
        h.db.First(&senderUser, userID)
        h.db.First(&recipientUser, msg.ToID)

		fmt.Println("akan send notif ke: " + recipientUser.Nickname)
		fmt.Println("notif dari: " + senderUser.Nickname)

        // Send notification via HTTP request
        notifPayload := map[string]string{
            "topic":   recipientUser.Nickname,
            "title":   "Pesan Baru",
            "content": fmt.Sprintf("%s %s: %s", senderUser.First_name, senderUser.Last_name, msg.Content),
        }
        
        jsonData, err := json.Marshal(notifPayload)
        if err != nil {
            fmt.Println("Error marshaling notification payload:", err)
        } else {
            // Create HTTP client and request with authorization header
            client := &http.Client{}
            req, err := http.NewRequest("POST", "http://localhost:5000/api/notif", bytes.NewBuffer(jsonData))
            if err != nil {
                fmt.Println("Error creating notification request:", err)
            } else {
                // Add authorization header with the original token
                req.Header.Set("Content-Type", "application/json")
                req.Header.Set("Authorization", "Bearer "+token_catch)
                
                // Send the request
                resp, err := client.Do(req)
                if err != nil {
                    fmt.Println("Error sending notification:", err)
                } else {
                    defer resp.Body.Close()
                    fmt.Println("Notification sent, status:", resp.Status)
                }
            }
        }
		// Send message to recipient if online
		var fullMessage models.Message
        h.db.Preload("From").First(&fullMessage, message.ID)

		if recipient, ok := clients[msg.ToID]; ok {
			recipient.WriteJSON(gin.H{
				"type":    "message",
                "message": fullMessage,
			})
		}
	}
}