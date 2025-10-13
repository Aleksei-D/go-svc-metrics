package middleware

import (
	"bytes"
	"crypto/rsa"
	"go-svc-metrics/internal/utils/crypto"
	"io"
	"net/http"
)

type cryptoRSAWriter struct {
	w          http.ResponseWriter
	privateKey *rsa.PrivateKey
}

func newCryptoRSAWriter(w http.ResponseWriter, privateKey *rsa.PrivateKey) *cryptoRSAWriter {
	return &cryptoRSAWriter{
		w:          w,
		privateKey: privateKey,
	}
}

func (c *cryptoRSAWriter) Header() http.Header {
	return c.w.Header()
}

func (c *cryptoRSAWriter) Write(p []byte) (int, error) {
	cryptData, err := crypto.DecryptRSAData(c.privateKey, p)
	if err != nil {
		return 0, err
	}

	return c.w.Write(cryptData)
}

func (c *cryptoRSAWriter) WriteHeader(statusCode int) {
	c.w.WriteHeader(statusCode)
}

type CryptoRSAMiddleware struct {
	PrivateKey *rsa.PrivateKey
}

func (c *CryptoRSAMiddleware) GetCryptoRSAMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		if c.PrivateKey != nil {
			cw := newCryptoRSAWriter(w, c.PrivateKey)
			ow = cw

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body.Close()

			newBody, err := crypto.DecryptRSAData(c.PrivateKey, bodyBytes)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(newBody))
		}

		next.ServeHTTP(ow, r)
	})
}
