package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jamesjoshuahill/secret/internal/aes"

	"github.com/gorilla/mux"
)

type GetSecretResponse struct {
	Data string `json:"data"`
}

type GetSecretRequest struct {
	Key string `json:"key"`
}

type GetSecret struct {
	Repository Repository
	Decrypt    DecryptFunc
}

func (g *GetSecret) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType != contentTypeJSON {
		writeError(w, http.StatusUnsupportedMediaType, "unsupported Content-Type")
		return
	}

	w.Header().Set("Content-Type", contentTypeJSON)

	vars := mux.Vars(r)
	id := vars["id"]

	body := &GetSecretRequest{}
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "decoding request body")
		return
	}

	secret, err := g.Repository.FindByID(id)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "wrong id or key")
		return
	}

	plainText, err := g.Decrypt(aes.Secret{
		Key:        body.Key,
		Nonce:      secret.Nonce,
		CipherText: secret.CipherText,
	})
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "wrong id or key")
		return
	}

	secretRes := &GetSecretResponse{
		Data: plainText,
	}

	err = json.NewEncoder(w).Encode(secretRes)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "encoding response body")
		return
	}
}
