package main

import (
	"net/http"

	"log"

	_ "git.nulana.com/bobrnor/battleship-server/db"
	"git.nulana.com/bobrnor/battleship-server/handlers"
)

func main() {
	mux := configMux()
	server := configServer(mux)
	startServer(server)
}

func configMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", handlers.AuthHandler())
	mux.HandleFunc("/longpoll", handlers.LongpollHandler())
	mux.HandleFunc("/game/search", handlers.SearchHandler())
	mux.HandleFunc("/game/start", handlers.StartHandler())
	mux.HandleFunc("/game/turn", handlers.TurnHandler())
	return mux
}

func configServer(mux *http.ServeMux) *http.Server {
	server := http.Server{
		Addr:    "0.0.0.0:80",
		Handler: mux,
	}
	return &server
}

func startServer(server *http.Server) {
	log.Printf("Battleship server started")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Listen and server failed %+v", err.Error())
	}
}
