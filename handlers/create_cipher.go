package handlers

import (
	"encoding/json"
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
	Encrypter  Encrypter
}

func (c *CreateCipher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reqBody := &createCipherRequest{}
	err := json.NewDecoder(r.Body).Decode(reqBody)
	if err != nil {
		writeError(w, http.StatusBadRequest, "decoding request body")
		return
	}

	key, cipherText, err := c.Encrypter.Encrypt(reqBody.Data)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "encrypting data")
		return
	}

	err = c.Repository.Store(repository.Cipher{
		ResourceID: reqBody.ResourceID,
		CipherText: cipherText,
		Key:        key,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "storing cipher")
		return
	}

	cipherRes := createCipherResponse{
		ResourceID: reqBody.ResourceID,
		Key:        key,
	}

	resBody, err := json.Marshal(cipherRes)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "encoding response body")
		return
	}

	w.Write(resBody)
}
