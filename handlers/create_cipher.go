package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jamesjoshuahill/ciphers/repository"
)

type createCipherRequest struct {
	Data       string `json:"data"`
	ResourceID string `json:"resource_id"`
}

type createCipherResponse struct {
	ResourceID string `json:"resource_id"`
	Key        string `json:"key"`
}

type CreateCipher struct {
	Repository Repository
}

func (c *CreateCipher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reqBody := &createCipherRequest{}
	err := json.NewDecoder(r.Body).Decode(reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, errorResponseBody("decoding request body"))
		return
	}

	key := "key for client-cipher-id"

	err = c.Repository.Store(repository.Cipher{
		ResourceID: reqBody.ResourceID,
		Data:       reqBody.Data,
		Key:        key,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, errorResponseBody("storing cipher"))
		return
	}

	cipherRes := createCipherResponse{
		ResourceID: reqBody.ResourceID,
		Key:        key,
	}

	resBody, err := json.Marshal(cipherRes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, errorResponseBody("encoding response body"))
		return
	}

	w.Write(resBody)
}
