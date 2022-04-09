package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

type user struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func registerServiceWithConsul() {
	config := consulapi.DefaultConfig()
	config.Address = "consul:8500"
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	registration := new(consulapi.AgentServiceRegistration)

	registration.ID = "user-service"
	registration.Name = "user-service"
	address := hostname()
	registration.Address = address
	port, err := strconv.Atoi(port()[1:len(port())])
	if err != nil {
		log.Fatalln(err)
	}
	registration.Port = port
	registration.Check = new(consulapi.AgentServiceCheck)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck", address, port)
	registration.Check.Interval = "5s"
	registration.Check.Timeout = "3s"
	consul.Agent().ServiceRegister(registration)
}

func users(w http.ResponseWriter, r *http.Request) {
	users := []user{
		{
			ID:   1,
			Name: "Admin",
		},
		{
			ID:   2,
			Name: "Ucup",
		},
		{
			ID:   3,
			Name: "Mamang",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&users)
}

func main() {
	registerServiceWithConsul()
	http.HandleFunc("/healthcheck", healthcheck)
	http.HandleFunc("/users", users)
	fmt.Printf("user service is up on port: %s", port())
	http.ListenAndServe(port(), nil)
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `user service is good`)
}

func port() string {
	p := os.Getenv("user_SERVICE_PORT")
	if len(strings.TrimSpace(p)) == 0 {
		return ":8200"
	}
	return fmt.Sprintf(":%s", p)
}

func hostname() string {
	hn, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	return hn
}
