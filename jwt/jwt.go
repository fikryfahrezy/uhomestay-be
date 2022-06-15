package jwt

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
)

func NewMiddleware(jwtKey []byte, jwtIssuerUrl string, jwtAudiences []string, customClaims validator.CustomClaims) func(next http.Handler) http.Handler {
	keyFunc := func(ctx context.Context) (interface{}, error) {
		// Our token must be signed using this data.
		return jwtKey, nil
	}

	// Set up the validator.
	jwtValidator, err := validator.New(
		keyFunc,
		validator.HS256,
		jwtIssuerUrl,
		jwtAudiences,
		validator.WithCustomClaims(func() validator.CustomClaims {
			return customClaims
		}),
	)
	if err != nil {
		log.Fatalf("Fail setup jwt validator: %s", err)
	}

	// Set up the middleware.
	jwtMidd := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithTokenExtractor(
			jwtmiddleware.MultiTokenExtractor(
				jwtmiddleware.AuthHeaderTokenExtractor,
				jwtmiddleware.CookieTokenExtractor("jwt"),
			),
		),
	).CheckJWT

	return jwtMidd
}

func MarshalClaims(r *http.Request) ([]byte, error) {
	claims := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)

	payload, err := json.Marshal(claims)
	if err != nil {
		return []byte{}, err
	}

	return payload, nil
}

func MarshalCustomClaims(r *http.Request) ([]byte, error) {
	claims := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)

	payload, err := json.Marshal(claims.CustomClaims)
	if err != nil {
		return []byte{}, err
	}

	return payload, nil
}

func DecodeCustomClaims(r *http.Request, in interface{}) error {
	payload, err := MarshalCustomClaims(r)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(payload, &in); err != nil {
		return err
	}

	return nil
}

type JwtPrivateClaim struct {
	Uid string `json:"uid"`
}

func (j *JwtPrivateClaim) Validate(ctx context.Context) error {
	return nil
}

type JwtPrivateAdminClaim struct {
	Uid     string `json:"uid"`
	IsAdmin bool   `json:"is_admin"`
}

func (j *JwtPrivateAdminClaim) Validate(ctx context.Context) error {
	if j.IsAdmin == false {
		return errors.New("jwt forbidden")
	}

	return nil
}

func Sign(ID, subject, issuer string, key []byte, audiences []string, notBefore, expiry, issuedAt time.Time, privateClaim interface{}) (string, error) {
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: key}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return "", err
	}

	cl := jwt.Claims{
		ID:        ID,
		Subject:   subject,
		Issuer:    issuer,
		Audience:  audiences,
		NotBefore: jwt.NewNumericDate(notBefore),
		Expiry:    jwt.NewNumericDate(expiry),
		IssuedAt:  jwt.NewNumericDate(issuedAt),
	}

	raw, err := jwt.Signed(sig).Claims(cl).Claims(privateClaim).CompactSerialize()
	if err != nil {
		return "", err
	}

	return raw, nil
}
