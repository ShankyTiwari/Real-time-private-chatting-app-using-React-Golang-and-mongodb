package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	config "private-chat/config"
	utils "private-chat/utils"
)

func main() {

	godotenv.Load()

	fmt.Println(
		fmt.Sprintf("%s%s%s%s", "Server will start at http://", os.Getenv("HOST"), ":", os.Getenv("PORT")),
	)

	config.ConnectDatabase()

	route := mux.NewRouter()

	AddApproutes(route)

	serverPath := ":" + os.Getenv("PORT")

	cors := utils.GetCorsConfig()

	http.ListenAndServe(serverPath, cors.Handler(route))
}
