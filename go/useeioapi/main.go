package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

func main() {
	args := GetArgs()

	log.Println("load data from folder:", args.DataDir)
	// TODO: check if data folder exists

	r := mux.NewRouter()

	// TODO: serve static files only when the folder exists
	log.Println("Create server with static files from:", args.StaticDir)
	fs := http.FileServer(http.Dir(args.StaticDir))
	r.Handle("/", NoCache(fs))

	// handle CORS preflight requests
	r.PathPrefix("/api").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			WriteAccessOptions(w)
		}).Methods("OPTIONS")

	// model routes
	r.HandleFunc("/api/models",
		HandleGetModels(args.DataDir)).Methods("GET")
	r.HandleFunc("/api/{model}/demands",
		HandleGetDemands(args.DataDir)).Methods("GET")
	r.HandleFunc("/api/{model}/demands/{id}",
		HandleGetDemand(args.DataDir)).Methods("GET")
	r.HandleFunc("/api/{model}/sectors",
		HandleGetSectors(args.DataDir)).Methods("GET")
	r.HandleFunc("/api/{model}/sectors/{id:.+}",
		HandleGetSector(args.DataDir)).Methods("GET")
	r.HandleFunc("/api/{model}/indicators",
		HandleGetIndicators(args.DataDir)).Methods("GET")
	r.HandleFunc("/api/{model}/indicators/{id}",
		HandleGetIndicator(args.DataDir)).Methods("GET")
	r.HandleFunc("/api/{model}/matrix/{matrix}",
		HandleGetMatrix(args.DataDir)).Methods("GET")
	r.HandleFunc("/api/{model}/calculate",
		HandleCalculate(args.DataDir)).Methods("POST")

	// register shutdown hook
	log.Println("Register shutdown routines")
	ossignals := make(chan os.Signal)
	signal.Notify(ossignals, syscall.SIGTERM)
	signal.Notify(ossignals, syscall.SIGINT)
	go func() {
		<-ossignals
		log.Println("Shutdown server")
		os.Exit(0)
	}()

	log.Println("Starting server at port:", args.Port)
	http.ListenAndServe(":"+args.Port, r)
}
