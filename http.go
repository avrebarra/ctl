package ctl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type ConfigHandler struct {
	PathPrefix string
	Ctl        *Ctl
}

func MakeHandler(cfg ConfigHandler) http.Handler {
	svr := makerouter(cfg)

	r := mux.NewRouter()
	r.HandleFunc(svr.config.PathPrefix+"config", svr.HandleList).Methods("GET")
	r.HandleFunc(svr.config.PathPrefix+"config/{key}", svr.HandleGetConfig).Methods("GET")
	r.HandleFunc(svr.config.PathPrefix+"config/{key}", svr.HandleSetConfig).Methods("PATCH")

	return r
}

// ***

type router struct {
	config ConfigHandler
}

func makerouter(cfg ConfigHandler) *router {
	return &router{config: cfg}
}

func (e *router) HandleList(w http.ResponseWriter, r *http.Request) {
	list := e.config.Ctl.List()

	respdata := list

	w.Header().Set("Content-Type", "application/json")
	if err := respjson(w, respdata); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (e *router) HandleGetConfig(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	key := mux.Vars(r)["key"]
	value, err := e.config.Ctl.Get(key).String()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respdata := Response{
		Key:   key,
		Value: value,
	}

	w.Header().Set("Content-Type", "application/json")
	if err = respjson(w, respdata); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (e *router) HandleSetConfig(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Value string `json:"value"`
	}
	type Response struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	key := mux.Vars(r)["key"]

	reqdata := Request{}
	err := json.NewDecoder(r.Body).Decode(&reqdata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	value, err := e.config.Ctl.Set(key, reqdata.Value).String()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	respdata := Response{
		Key:   key,
		Value: value,
	}

	w.Header().Set("Content-Type", "application/json")
	if err = respjson(w, respdata); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// ***

func respjson(w io.Writer, obj interface{}) (err error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}

	_, err = fmt.Fprintf(w, "%s", data)
	if err != nil {
		return
	}

	return
}
