package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	hostname          string
	envSpecificConfig Configuration
)

type Configuration struct {
	Greetings []string `yaml:"greetings"`
}

func loadConfig(filename string) (Configuration, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Configuration{}, err
	}

	var c Configuration
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return Configuration{}, err
	}

	return c, nil
}

func main() {
	var err error

	deployEnv := os.Getenv("DEPLOY_ENV")
	log.Println("Initializing greetings api server for deployment environment", deployEnv)

	certPrivateKeyFile := os.Getenv("CERT_PRIVATE_KEY_FILE")
	log.Println(os.ExpandEnv("Will read TLS certificate private key from '${CERT_PRIVATE_KEY_FILE}'"))

	if certPrivateKeyFile != "" {

		certPrivateKey, err := ioutil.ReadFile(certPrivateKeyFile)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("chars in private key", len(certPrivateKey))
	}

	hostname, err = os.Hostname()
	if err != nil {
		log.Println("could not read hostname")
	}

	if os.Getenv("DEBUG") == "true" {
		log.Println("hostname", hostname)
	}

	envSpecificConfigFile := fmt.Sprintf("/config/config.%s.yml", deployEnv)
	log.Println("Loading env-specific configurations from", envSpecificConfigFile)
	envSpecificConfig, err = loadConfig(envSpecificConfigFile)
	for _, greeting := range envSpecificConfig.Greetings {
		log.Println(greeting)
	}

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/greeting", serveGreeting)

	log.Println("Initialization complete, starting http service on", hostname)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func serveIndex(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")

	fmt.Fprintln(resp, "Welcome to the Greetings API Server!")
	fmt.Fprintln(resp, fmt.Sprintf("Container with id %s responded at %s", hostname, time.Now().UTC()))
}

func serveGreeting(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// select a random greeting
	greeting := "Hello World!"
	fmt.Fprintln(resp, greeting)
}
