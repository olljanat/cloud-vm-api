package api

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

type VMListResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func ListVMsHandler(w http.ResponseWriter, r *http.Request) {
	envName := r.URL.Query().Get("environment")
	host, env, status, err := getHost(r, envName)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	vms, err := host.GetIVMs()
	if err != nil {
		http.Error(w, "Failed to list VMs: "+err.Error(), http.StatusInternalServerError)
		return
	}
	resp := []VMListResponse{}
	for _, vm := range vms {
		id := vm.GetGlobalId()
		if env.Cloud == "Azure" {
			parts := strings.Split(id, "/")
			slices.Reverse(parts)
			id = parts[0]
		}

		resp = append(resp, VMListResponse{
			ID:     id,
			Name:   vm.GetName(),
			Status: vm.GetStatus(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
