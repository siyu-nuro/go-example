package endpoint

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrAuthFailed             = errors.New("failed to authenticate request")
)

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func validateAdminAuth(token *jwt.Token) (interface{}, error) {
	return validateAuthForAudience(token, []string{"Nuro Admin"})
}

func validateRetailerAuth(token *jwt.Token) (interface{}, error) {
	return validateAuthForAudience(token, []string{"Nuro Retailer", "Nuro Admin"})
}

func validateAuthForAudience(token *jwt.Token, intendedAudiences []string) (interface{}, error) {
	// Validate audience
	var validAudience bool
	for _, audience := range intendedAudiences {
		if token.Claims.(jwt.MapClaims).VerifyAudience(audience, false) {
			validAudience = true
			break
		}
	}
	if !validAudience {
		fmt.Errorf("failed to validate auth - invalid audience")
		return token, ErrAuthFailed
	}
	// Validate issuer
	if !token.Claims.(jwt.MapClaims).VerifyIssuer("auth0 domain", false) {
		fmt.Errorf("failed to validate auth - invalid issuer")
		return token, ErrAuthFailed
	}

	cert, err := getPemCert(token)
	if err != nil {
		fmt.Errorf("failed to validate auth - cannot get cert from token. error %v", err)
		return nil, ErrAuthFailed
	}

	result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	return result, nil
}

func getPemCert(token *jwt.Token) (string, error) {
	var cert string

	var jwks = Jwks{}
	if err := json.Unmarshal([]byte("auth0-public-key"), &jwks); err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
			return cert, nil
		}
	}

	return cert, errors.New("unable to find appropriate key for cert")
}
