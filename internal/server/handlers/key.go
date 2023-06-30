package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vtno/zypher/internal/store"
)

type KeyHandler struct {
	store store.Store
}

type KeyPostRequest struct {
	Name string `json:"name"`
	Env  string `json:"env"`
	Key  string `json:"key"`
}

type KeyGetRequest struct {
	Name string `json:"name"`
	Env  string `json:"env"`
}

type KeyGetResponse struct {
	Key string `json:"key"`
}

func NewKeyHandler(store store.Store) *KeyHandler {
	return &KeyHandler{
		store: store,
	}
}

func (kh *KeyHandler) Get(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	kgr := &KeyGetRequest{
		Name: params.Get("name"),
		Env:  params.Get("env"),
	}
	lookupKey := fmt.Sprintf("%s#%s", kgr.Name, kgr.Env)
	v, err := kh.store.Get(lookupKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if v == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	response := &KeyGetResponse{
		Key: v,
	}
	res, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(res)
}

func (kh *KeyHandler) Post(w http.ResponseWriter, r *http.Request) {
	var kpr KeyPostRequest
	if err := json.NewDecoder(r.Body).Decode(&kpr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	lookupKey := fmt.Sprintf("%s#%s", kpr.Name, kpr.Env)
	if err := kh.store.Set(lookupKey, kpr.Key); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
