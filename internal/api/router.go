package api

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/vm", CreateVMHandler).Methods("POST")
	r.HandleFunc("/vm", ListVMsHandler).Methods("GET")
	r.HandleFunc("/vm/{id}/start", StartVMHandler).Methods("GET")
	r.HandleFunc("/vm/{id}/stop", StopVMHandler).Methods("GET")
	r.HandleFunc("/vm/{id}", DeleteVMHandler).Methods("DELETE")
	return r
}
