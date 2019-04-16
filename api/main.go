package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var (
	hostname  string
	deployEnv string
	config    Configuration
	r         *rand.Rand
)

func main() {
	var err error

	deployEnv = os.Getenv("DEPLOY_ENV")
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

	commonConfigFile := "/config/config.common.yml"
	log.Println("Loading env-specific configurations from", commonConfigFile)
	commonConfig, err := LoadConfig(commonConfigFile)

	envSpecificConfigFile := fmt.Sprintf("/config/config.%s.yml", deployEnv)
	log.Println("Loading env-specific configurations from", envSpecificConfigFile)
	envSpecificConfig, err := LoadConfig(envSpecificConfigFile)

	config = Configuration{
		Greetings: append(commonConfig.Greetings, envSpecificConfig.Greetings...),
	}

	log.Printf("Greetings: %v", config.Greetings)

	r = rand.New(rand.NewSource(42))

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
	fmt.Fprintln(resp, fmt.Sprintf("DEPLOY_ENV: %s", deployEnv))
}

func SelectRandom(strings []string, r *rand.Rand) string {
	return strings[r.Intn(len(strings))]
}

func serveGreeting(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	greeting := SelectRandom(config.Greetings, r)
	fmt.Fprintln(resp, greeting)
}
