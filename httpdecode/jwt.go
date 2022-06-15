package httpdecode

import (
	"encoding/json"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

func Jwt(r *http.Request, in interface{}) error {
	claims := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)

	payload, err := json.Marshal(claims.CustomClaims)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(payload, &in); err != nil {
		return err
	}

	return nil
}
