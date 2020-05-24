package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

func unRegisterAndCloseConnection(c *Client) {
	c.hub.unregister <- c
	c.webSocketConnection.Close()
}

func setSocketPayloadReadConfig(c *Client) {
	c.webSocketConnection.SetReadLimit(maxMessageSize)
	c.webSocketConnection.SetReadDeadline(time.Now().Add(pongWait))
	c.webSocketConnection.SetPongHandler(func(string) error { c.webSocketConnection.SetReadDeadline(time.Now().Add(pongWait)); return nil })
}

func handleSocketPayloadEvents(client *Client, socketEventPayload SocketEventStruct) {
	type chatlistResponseStruct struct {
		Type     string      `json:"type"`
		Chatlist interface{} `json:"chatlist"`
	}
	switch socketEventPayload.EventName {

	case "join":
		userID := (socketEventPayload.EventPayload).(string)
		userDetails := GetUserByUserID(userID)
		if userDetails == (UserDetailsStruct{}) {
			log.Println("An invalid user with userID " + userID + " tried to connect to Chat Server.")
		} else {
			if userDetails.Online == "N" {
				log.Println("A logged out user with userID " + userID + " tried to connect to Chat Server.")
			} else {
				newUserOnlinePayload := SocketEventStruct{
					EventName: "chatlist-response",
					EventPayload: chatlistResponseStruct{
						Type: "new-user-joined",
						Chatlist: UserDetailsResponsePayloadStruct{
							Online:   userDetails.Online,
							UserID:   userDetails.ID,
							Username: userDetails.Username,
						},
					},
				}
				BroadcastSocketEventToAllClientExceptMe(client.hub, newUserOnlinePayload, userDetails.ID)

				allOnlineUsersPayload := SocketEventStruct{
					EventName: "chatlist-response",
					EventPayload: chatlistResponseStruct{
						Type:     "my-chat-list",
						Chatlist: GetAllOnlineUsers(userDetails.ID),
					},
				}
				EmitToSpecificClient(client.hub, allOnlineUsersPayload, userDetails.ID)
			}
		}
	case "disconnect":
		if socketEventPayload.EventPayload != nil {

			userID := (socketEventPayload.EventPayload).(string)
			userDetails := GetUserByUserID(userID)
			UpdateUserOnlineStatusByUserID(userID, "N")

			BroadcastSocketEventToAllClient(client.hub, SocketEventStruct{
				EventName: "chatlist-response",
				EventPayload: chatlistResponseStruct{
					Type: "user-disconnected",
					Chatlist: UserDetailsResponsePayloadStruct{
						Online:   "N",
						UserID:   userDetails.ID,
						Username: userDetails.Username,
					},
				},
			})
		}
	case "message":
		message := (socketEventPayload.EventPayload.(map[string]interface{})["message"]).(string)
		fromUserID := (socketEventPayload.EventPayload.(map[string]interface{})["fromUserID"]).(string)
		toUserID := (socketEventPayload.EventPayload.(map[string]interface{})["toUserID"]).(string)

		if message != "" && fromUserID != "" && toUserID != "" {

			messagePacket := MessagePayloadStruct{
				FromUserID: fromUserID,
				Message:    message,
				ToUserID:   toUserID,
			}
			StoreNewChatMessages(messagePacket)
			allOnlineUsersPayload := SocketEventStruct{
				EventName:    "message-response",
				EventPayload: messagePacket,
			}
			EmitToSpecificClient(client.hub, allOnlineUsersPayload, toUserID)

		}
	}
}

func (c *Client) readPump() {
	var socketEventPayload SocketEventStruct

	// Unregistering the client and closing the connection
	defer unRegisterAndCloseConnection(c)

	// Setting up the Payload configuration
	setSocketPayloadReadConfig(c)

	for {
		// ReadMessage is a helper method for getting a reader using NextReader and reading from that reader to a buffer.
		_, payload, err := c.webSocketConnection.ReadMessage()

		decoder := json.NewDecoder(bytes.NewReader(payload))
		decoderErr := decoder.Decode(&socketEventPayload)

		if decoderErr != nil {
			log.Printf("error: %v", decoderErr)
			break
		}

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error ===: %v", err)
			}
			break
		}

		//  Getting the proper Payload to send the client
		handleSocketPayloadEvents(c, socketEventPayload)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.webSocketConnection.Close()
	}()
	for {
		select {
		case payload, ok := <-c.send:

			reqBodyBytes := new(bytes.Buffer)
			json.NewEncoder(reqBodyBytes).Encode(payload)
			finalPayload := reqBodyBytes.Bytes()

			c.webSocketConnection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.webSocketConnection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.webSocketConnection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(finalPayload)

			n := len(c.send)
			for i := 0; i < n; i++ {
				json.NewEncoder(reqBodyBytes).Encode(<-c.send)
				w.Write(reqBodyBytes.Bytes())
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.webSocketConnection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.webSocketConnection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// CreateNewSocketUser creates a new socket user
func CreateNewSocketUser(hub *Hub, connection *websocket.Conn, userID string) {
	client := &Client{
		hub:                 hub,
		webSocketConnection: connection,
		send:                make(chan SocketEventStruct),
		userID:              userID,
	}

	go client.writePump()
	go client.readPump()

	client.hub.register <- client
}

// HandleUserRegisterEvent will handle the Join event for New socket users
func HandleUserRegisterEvent(hub *Hub, client *Client) {
	hub.clients[client] = true
	handleSocketPayloadEvents(client, SocketEventStruct{
		EventName:    "join",
		EventPayload: client.userID,
	})
}

// HandleUserDisconnectEvent will handle the Disconnect event for socket users
func HandleUserDisconnectEvent(hub *Hub, client *Client) {
	_, ok := hub.clients[client]
	if ok {
		delete(hub.clients, client)
		close(client.send)

		handleSocketPayloadEvents(client, SocketEventStruct{
			EventName:    "disconnect",
			EventPayload: client.userID,
		})
	}
}

// EmitToSpecificClient will emit the socket event to specific socket user
func EmitToSpecificClient(hub *Hub, payload SocketEventStruct, userID string) {
	for client := range hub.clients {
		if client.userID == userID {
			select {
			case client.send <- payload:
			default:
				close(client.send)
				delete(hub.clients, client)
			}
		}
	}
}

// BroadcastSocketEventToAllClient will emit the socket events to all socket users
func BroadcastSocketEventToAllClient(hub *Hub, payload SocketEventStruct) {
	for client := range hub.clients {
		select {
		case client.send <- payload:
		default:
			close(client.send)
			delete(hub.clients, client)
		}
	}
}

// BroadcastSocketEventToAllClientExceptMe will emit the socket events to all socket users,
// except the user who is emitting the event
func BroadcastSocketEventToAllClientExceptMe(hub *Hub, payload SocketEventStruct, myUserID string) {
	for client := range hub.clients {
		if client.userID != myUserID {
			select {
			case client.send <- payload:
			default:
				close(client.send)
				delete(hub.clients, client)
			}
		}
	}
}
