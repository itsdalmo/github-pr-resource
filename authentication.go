package resource

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type InstallationResponse struct {
	Id int
}

type TokenResponse struct {
	Token string
}

func GenerateAccessToken(s *Source, now time.Time) (string, error) {
	if s.AccessToken != "" {
		return s.AccessToken, nil
	}

	decode, _ := pem.Decode([]byte(s.PrivateKey))
	key, _ := x509.ParsePKCS1PrivateKey(decode.Bytes)
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: key}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		panic(err)
	}

	cl := jwt.Claims{
		Issuer:   s.AppId,
		IssuedAt: jwt.NewNumericDate(now),
		Expiry:   jwt.NewNumericDate(now.Add(9 * time.Minute)),
	}
	signedJwt, err := jwt.Signed(sig).Claims(cl).CompactSerialize()
	if err != nil {
		panic(err)
	}

	var endpoint string
	if s.V3Endpoint != "" {
		endpoint = strings.TrimRight(s.V3Endpoint, "/")
	} else {
		endpoint = "https://api.github.com"
	}

	installationResponse := callApi("GET", endpoint+"/repos/"+s.Repository+"/installation", signedJwt)
	var ir InstallationResponse
	err = json.NewDecoder(installationResponse.Body).Decode(&ir)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("Error decoding installation response with status %d", installationResponse.StatusCode))
		_, _ = fmt.Fprintln(os.Stderr, installationResponse.Body)
		panic(err)
	}

	tokenResponse := callApi("POST", endpoint+"/app/installations/"+strconv.Itoa(ir.Id)+"/access_tokens", signedJwt)

	var tr TokenResponse
	err = json.NewDecoder(tokenResponse.Body).Decode(&tr)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("Error decoding token response with status %d", tokenResponse.StatusCode))
		_, _ = fmt.Fprintln(os.Stderr, tokenResponse.Body)
		panic(err)
	}

	return tr.Token, nil
}

func callApi(method string, endpoint string, signedJwt string) *http.Response {
	tokenRequest, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		panic(err)
	}
	tokenRequest.Header.Add("Authorization", "Bearer "+signedJwt)
	tokenRequest.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")
	client := &http.Client{}
	response, err := client.Do(tokenRequest)

	if err != nil {
		panic(err)
	}
	return response
}
