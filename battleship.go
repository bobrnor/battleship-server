package main

import (
	"net/http"

	"go.uber.org/zap"

	_ "git.nulana.com/bobrnor/battleship-server/db"
	"git.nulana.com/bobrnor/battleship-server/handlers"
)

func main() {
	configLogger()
	mux := configMux()
	server := configServer(mux)
	startServer(server)
}

func configLogger() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
}

func configMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", handlers.AuthHandler())
	mux.HandleFunc("/game/search", handlers.SearchHandler())
	//mux.HandleFunc("/game/confirm", handlers.ConfirmHandler())
	mux.HandleFunc("/game/start", handlers.StartHandler())
	mux.HandleFunc("/game/turn", handlers.TurnHandler())
	mux.HandleFunc("/game/longpoll", handlers.LongpollHandler())
	return mux
}

func configServer(mux *http.ServeMux) *http.Server {
	server := http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: mux,
	}
	return &server
}

func startServer(server *http.Server) {
	zap.S().Infow("Battleship server started")
	if err := server.ListenAndServe(); err != nil {
		zap.S().Fatalw("Listen and server failed",
			"err", err,
		)
	}
}
