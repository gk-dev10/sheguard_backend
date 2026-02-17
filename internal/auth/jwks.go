package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/MicahParks/keyfunc"
)

var JWKS *keyfunc.JWKS

// apiKeyTransport injects the Supabase `apikey` header into outgoing requests.
type apiKeyTransport struct {
	apiKey string
	rt     http.RoundTripper
}

func (t *apiKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.apiKey != "" {
		req.Header.Set("apikey", t.apiKey)
	}
	return t.rt.RoundTrip(req)
}

func InitJWKS() error {
	jwksURL := os.Getenv("SUPABASE_URL") + "/auth/v1/keys"
	// fmt.Println("JWKS URL:", jwksURL)

	apiKey := os.Getenv("SUPABASE_ANON_KEY")

	client := &http.Client{
		Transport: &apiKeyTransport{apiKey: apiKey, rt: http.DefaultTransport},
		Timeout:   10 * time.Second,
	}
	j, err := keyfunc.Get(jwksURL, keyfunc.Options{Client: client})
	if err != nil {
		altURL := os.Getenv("SUPABASE_URL") + "/auth/v1/.well-known/jwks.json"
		// fmt.Println("JWKS fallback URL:", altURL)
		j, err = keyfunc.Get(altURL, keyfunc.Options{Client: client})
		if err != nil {
			return err
		}
	}

	JWKS = j
	return nil
}
