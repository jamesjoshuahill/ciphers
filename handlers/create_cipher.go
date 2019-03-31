package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jamesjoshuahill/ciphers/repository"
)

type createCipherRequest struct {
	Data string `json:"data"`
	ID   string `json:"id"`
}

type createCipherResponse struct {
	Key string `json:"key"`
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
		ID:         reqBody.ID,
		CipherText: cipherText,
	})
	if err != nil {
		writeError(w, http.StatusConflict, "cipher already exists")
		return
	}

	cipherRes := createCipherResponse{
		Key: key,
	}

	resBody, err := json.Marshal(cipherRes)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "encoding response body")
		return
	}

	w.Write(resBody)
}
