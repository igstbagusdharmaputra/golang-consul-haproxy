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

type product struct {
	ID    uint64  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func registerServiceWithConsul() {
	config := consulapi.DefaultConfig()
	config.Address = "consul:8500"
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	registration := new(consulapi.AgentServiceRegistration)

	registration.ID = "product-service"
	registration.Name = "product-service"
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

func Products(w http.ResponseWriter, r *http.Request) {
	products := []product{
		{
			ID:    1,
			Name:  "Baju",
			Price: 50000.00,
		},
		{
			ID:    2,
			Name:  "Celana",
			Price: 25000.00,
		},
		{
			ID:    3,
			Name:  "Sandal",
			Price: 1500000.00,
		},
		{
			ID:    4,
			Name:  "Es Teh Tarik",
			Price: 50000.00,
		},
		{
			ID:    5,
			Name:  "Lenovo Laptop",
			Price: 20000000.00,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&products)
}

func main() {
	registerServiceWithConsul()
	http.HandleFunc("/healthcheck", healthcheck)
	http.HandleFunc("/products", Products)
	fmt.Printf("product service is up on port: %s", port())
	http.ListenAndServe(port(), nil)
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `product service is good`)
}

func port() string {
	p := os.Getenv("PRODUCT_SERVICE_PORT")
	if len(strings.TrimSpace(p)) == 0 {
		return ":8100"
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
