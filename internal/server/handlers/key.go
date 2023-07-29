package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vtno/zypher/internal/server/store"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type KeyHandler struct {
	store store.Store
}

type KeyPostRequest struct {
	Name string `json:"name" validate:"required"`
	Env  string `json:"env" validate:"required"`
	Key  string `json:"key" validate:"required"`
}

type KeyGetRequest struct {
	Name string `json:"name" validate:"required"`
	Env  string `json:"env" validate:"required"`
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
	validate := validator.New()
	if err := validate.Struct(kgr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
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
	logger := r.Context().Value("logger").(*zap.Logger)

	err := json.NewDecoder(r.Body).Decode(&kpr)
	if err != nil {
		logger.Error("error decoding as json", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	validate := validator.New()
	if err := validate.Struct(kpr); err != nil {
		logger.Error("error validating request body", zap.Error(err))
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
