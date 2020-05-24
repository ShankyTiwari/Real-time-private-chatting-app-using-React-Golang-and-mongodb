package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	handlers "private-chat/handlers"
)

// AddApproutes will add the routes for the application
func AddApproutes(route *mux.Router) {

	log.Println("Loadeding Routes...")

	hub := handlers.NewHub()
	go hub.Run()

	route.HandleFunc("/", handlers.RenderHome)

	route.HandleFunc("/isUsernameAvailable/{username}", handlers.IsUsernameAvailable)

	route.HandleFunc("/login", handlers.Login).Methods("POST")

	route.HandleFunc("/registration", handlers.Registertation).Methods("POST")

	route.HandleFunc("/userSessionCheck/{userID}", handlers.UserSessionCheck)

	route.HandleFunc("/getConversation/{toUserID}/{fromUserID}", handlers.GetMessagesHandler)

	route.HandleFunc("/ws/{userID}", func(responseWriter http.ResponseWriter, request *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		// Reading username from request parameter
		userID := mux.Vars(request)["userID"]

		// Upgrading the HTTP connection socket connection
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		connection, err := upgrader.Upgrade(responseWriter, request, nil)
		if err != nil {
			log.Println(err)
			return
		}

		handlers.CreateNewSocketUser(hub, connection, userID)

	})

	log.Println("Routes are Loaded.")
}
