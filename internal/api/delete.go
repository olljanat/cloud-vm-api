package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

func DeleteVMHandler(w http.ResponseWriter, r *http.Request) {
	envName := r.URL.Query().Get("environment")
	host, env, status, err := getHost(r, envName)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	vars := mux.Vars(r)
	vmID := vars["id"]
	vm, err := getVMByID(host, env, vmID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	ctx := context.Background()
	if err := vm.DeleteVM(ctx); err != nil {
		http.Error(w, "Failed to delete VM: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
