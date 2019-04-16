package handler_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/jamesjoshuahill/ciphers/aes"
	"github.com/jamesjoshuahill/ciphers/handler"
	"github.com/jamesjoshuahill/ciphers/handler/fake"
	"github.com/jamesjoshuahill/ciphers/inmemory"

	"github.com/gorilla/mux"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetSecret", func() {
	var (
		repo      *fake.Repo
		decrypter *fake.Decrypter
		res       *httptest.ResponseRecorder
		req       *http.Request
		router    *mux.Router
	)

	BeforeEach(func() {
		repo = new(fake.Repo)
		decrypter = new(fake.Decrypter)
		res = httptest.NewRecorder()
		router = mux.NewRouter()

		var err error
		req, err = http.NewRequest("GET", "/v1/secrets/client-secret-id", strings.NewReader(`{
			"key": "key for client-secret-id"
		}`))
		Expect(err).NotTo(HaveOccurred())
		req.Header.Set("Content-Type", "application/json")
	})

	It("retrieves the secret", func() {
		router.Handle("/v1/secrets/{id}", &handler.GetSecret{Repository: repo, Decrypter: decrypter})

		router.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusOK), res.Body.String())
		Expect(repo.FindByResourceIDCall.Received.ID).To(Equal("client-secret-id"))
	})

	It("decrypts the ciphertext", func() {
		repo.FindByResourceIDCall.Returns.Secret = inmemory.Secret{
			ID:         "client-secret-id",
			Nonce:      "some nonce",
			CipherText: "some cipher text",
		}
		router.Handle("/v1/secrets/{id}", &handler.GetSecret{Repository: repo, Decrypter: decrypter})

		router.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusOK), res.Body.String())
		Expect(decrypter.DecryptCall.Received.Secret).To(Equal(aes.Secret{
			Key:        "key for client-secret-id",
			Nonce:      "some nonce",
			CipherText: "some cipher text",
		}))
	})

	It("fails when the request content type is not JSON", func() {
		req.Header.Set("Content-Type", "text/plain")
		handler := handler.GetSecret{Repository: repo, Decrypter: decrypter}

		handler.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusUnsupportedMediaType))
		Expect(res.Body.String()).To(ContainSubstring("unsupported Content-Type"))
	})

	It("fails when the request body cannot be parsed", func() {
		req.Body = ioutil.NopCloser(strings.NewReader("not json"))
		router.Handle("/v1/secrets/{id}", &handler.GetSecret{Repository: repo, Decrypter: decrypter})

		router.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusBadRequest), res.Body.String())
		Expect(res.Body.String()).To(ContainSubstring("decoding request body"))
	})

	It("fails when the secret id is wrong", func() {
		repo.FindByResourceIDCall.Returns.Error = errors.New("fake error")
		router.Handle("/v1/secrets/{id}", &handler.GetSecret{Repository: repo, Decrypter: decrypter})

		router.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusUnprocessableEntity), res.Body.String())
		Expect(res.Body.String()).To(ContainSubstring("wrong id or key"))
	})

	It("fails when the secret key is wrong", func() {
		decrypter.DecryptCall.Returns.Error = errors.New("fake error")
		handler := handler.GetSecret{Repository: repo, Decrypter: decrypter}

		handler.ServeHTTP(res, req)

		Expect(res.Code).To(Equal(http.StatusUnprocessableEntity))
		Expect(res.Body.String()).To(ContainSubstring("wrong id or key"))
	})
})
