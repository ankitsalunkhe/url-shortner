package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/ankitsalunkhe/url-shortner/internal/api"
	"github.com/stretchr/testify/assert"
)

func Test_UrlShortner(t *testing.T) {
	t.Run("successfully create new short url", func(t *testing.T) {
		req := api.Request{
			Url: "https://github.com/ankitsalunkhe/url-shortner",
		}

		body, err := json.Marshal(req)
		assert.NoError(t, err)

		httpRequest, err := http.NewRequest(http.MethodPost, "http://url-shortner:8088/api/v1/url", bytes.NewBuffer(body))
		assert.NoError(t, err)

		httpRequest.Header.Add("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(httpRequest)
		assert.NoError(t, err)

		body, err = io.ReadAll(res.Body)
		log.Printf("body : %s", body)

		defer res.Body.Close()

		assert.Equal(t, res.StatusCode, http.StatusOK)
	})
}
