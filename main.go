package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	configFile := flag.String("config", "", "config")
	flag.Parse()
	if *configFile == "" {
		log.Fatal("Config file not specified")
	}

	config, err := readConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.Database.Host,
			config.Database.Port,
			config.Database.User,
			config.Database.Password,
			config.Database.Name))
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	manager := Manager{rep: Repository{db: db}}
	r := mux.NewRouter()
	r.HandleFunc("/add", manager.Add).Methods("POST")
	r.HandleFunc("/status", manager.Status).Methods("GET")
	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:" + config.ServerPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sig
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}()

	log.Printf("Starting HTTP server on %s port", config.ServerPort)
	log.Fatal(srv.ListenAndServe())
}

func readConfig(fileName string) (conf Config, err error) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	return
}
