package jwks

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
)

// ServeJWKS starts an HTTP server exposing the JWKS (JSON Web Key Set)
// This allows Envoy to fetch the public key to verify JWTs issued by the auth service.
func ServeJWKS(pubKey *rsa.PublicKey) {
    // Convert the RSA public key modulus (n) to base64url encoding
    n := base64.RawURLEncoding.EncodeToString(pubKey.N.Bytes())

    // Convert the RSA public key exponent (e) to base64url encoding
    e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pubKey.E)).Bytes())

    // Construct the JWKS structure with a single RSA key
    jwks := struct {
        Keys []interface{} `json:"keys"`
    }{
        Keys: []interface{}{
            map[string]string{
                "kty": "RSA",       // Key type
                "kid": "key1",      // Key ID (optional, used for key rotation)
                "use": "sig",       // Key usage: signature
                "alg": "RS256",     // Algorithm used for signing
                "n":   n,           // Public modulus
                "e":   e,           // Public exponent
            },
        },
    }

    // HTTP handler for JWKS endpoint
    http.HandleFunc("/.well-known/jwks.json", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(jwks) // Return JWKS as JSON
    })

    fmt.Println("JWKS server running on :5000")
    
    http.ListenAndServe(":5000", nil)
}