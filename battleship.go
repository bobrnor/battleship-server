package main

import (
	"net/http"

	"go.uber.org/zap"

	"git.nulana.com/bobrnor/battleship/auth"
)

func main() {
	configLogger()
	mux := configMux()
	server := configServer(mux)
	start(server)
}

func configLogger() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
}

func configMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth", auth.Handler())
	return mux
}

func configServer(mux *http.ServeMux) *http.Server {
	server := http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: mux,
	}
	return &server
}

func start(server *http.Server) {
	zap.S().Infow("Battleship server started")
	if err := server.ListenAndServe(); err != nil {
		zap.S().Fatalw("Listen and server failed",
			"err", err,
		)
	}
}
