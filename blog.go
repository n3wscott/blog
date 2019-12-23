package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type envConfig struct {
	FilePath string `envconfig:"FILE_PATH" default:"/var/run/ko/" required:"true"`
	Port     int    `envconfig:"PORT" default:"8080" required:"true"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}
	if !strings.HasSuffix(env.FilePath, "/") {
		env.FilePath = env.FilePath + "/"
	}

	m := http.NewServeMux()
	m.Handle("/", http.FileServer(http.Dir(env.FilePath)))

	log.Println("Listening on", env.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", env.Port), m))
}
