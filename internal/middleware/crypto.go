package middleware

import (
	"bytes"
	"encoding/hex"
	"go-svc-metrics/internal/config"
	"go-svc-metrics/internal/utils/crypto"
	"io"
	"net/http"
)

type cryptoWriter struct {
	w    http.ResponseWriter
	key  string
	body []byte
}

func newCryptoWriter(w http.ResponseWriter, key string) *cryptoWriter {
	return &cryptoWriter{
		w:   w,
		key: key,
	}
}

func (c *cryptoWriter) Header() http.Header {
	return c.w.Header()
}

func (c *cryptoWriter) Write(p []byte) (int, error) {
	cryptData, err := crypto.EncryptData(c.key, p)
	if err != nil {
		return 0, err
	}

	return c.w.Write(cryptData)
}

func (c *cryptoWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		hash := crypto.GetHash(c.key, c.body)
		c.w.Header().Set("HashSHA256", hex.EncodeToString(hash))
	}
	c.w.WriteHeader(statusCode)
}

type CryptoMiddleware struct {
	Config *config.Config
}

func (c *CryptoMiddleware) GetCryptoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		HashSHA256Header := r.Header.Get("HashSHA256")
		if HashSHA256Header != "" {
			cw := newCryptoWriter(w, *c.Config.Key)
			ow = cw

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body.Close()
			messageMAC, err := hex.DecodeString(HashSHA256Header)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			newBody, err := crypto.DecryptData(*c.Config.Key, bodyBytes)
			if ok := crypto.ValidMAC(newBody, messageMAC, *c.Config.Key); !ok {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(newBody))
		}

		next.ServeHTTP(ow, r)
	})
}
