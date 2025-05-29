package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/olljanat/cloud-vm-api/internal/api"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

func main() {
	router := api.NewRouter()
	log.Println("Starting server on :8080")
	fmt.Println("Registered Cloudmux providers are:", cloudprovider.GetRegistedProviderIds())
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
