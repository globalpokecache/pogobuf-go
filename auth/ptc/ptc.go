package ptc

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/context/ctxhttp"
)

const authorizeURL = "https://sso.pokemon.com/sso/oauth2.0/accessToken"
const loginURL = "https://sso.pokemon.com/sso/login?service=https://sso.pokemon.com/sso/oauth2.0/callbackAuthorize"

const redirectURI = "https://www.nianticlabs.com/pokemongo/error"
const clientSecret = "w8ScCUXJQc6kXKw8FiOhd8Fixzht18Dq3PEVkUCP5ZPxtgyWsbTvWHFLm2wNY0JR"
const clientID = "mobile-app_pokemon-go"

const providerString = "ptc"

type loginRequest struct {
	Lt        string   `json:"lt"`
	Execution string   `json:"execution"`
	Errors    []string `json:"errors,omitempty"`
}

// Provider contains data about and manages the session with the Pokémon Trainer's Club
type Provider struct {
}

// NewProvider constructs a Pokémon Trainer's Club auth provider instance
func NewProvider() *Provider {
	return &Provider{}
}

// Login retrieves an access token from the Pokémon Trainer's Club
func (p *Provider) Login(ctx context.Context, username, password string) (string, error) {
	options := &cookiejar.Options{}
	jar, _ := cookiejar.New(options)
	httpClient := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("Use the last error")
		},
	}

	req1, _ := http.NewRequest("GET", loginURL, nil)
	req1.Header.Set("User-Agent", "pokemongo/1 CFNetwork/808.2.16 Darwin/16.3.0")

	resp1, err1 := ctxhttp.Do(ctx, httpClient, req1)
	if err1 != nil {
		return "", errors.New("Could not start login process, the website might be down")
	}

	defer resp1.Body.Close()
	body1, _ := ioutil.ReadAll(resp1.Body)
	var loginRespBody loginRequest
	json.Unmarshal(body1, &loginRespBody)
	resp1.Body.Close()

	loginForm := url.Values{}
	loginForm.Set("lt", loginRespBody.Lt)
	loginForm.Set("execution", loginRespBody.Execution)
	loginForm.Set("_eventId", "submit")
	loginForm.Set("username", username)
	loginForm.Set("password", password)

	loginFormData := strings.NewReader(loginForm.Encode())

	req2, _ := http.NewRequest("POST", loginURL, loginFormData)
	req2.Header.Set("User-Agent", "pokemongo/1 CFNetwork/808.2.16 Darwin/16.3.0")
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp2, err2 := ctxhttp.Do(ctx, httpClient, req2)
	if _, ok2 := err2.(*url.Error); !ok2 {
		if resp2.Body != nil {
			defer resp2.Body.Close()
			body2, _ := ioutil.ReadAll(resp2.Body)
			var respBody loginRequest
			json.Unmarshal(body2, &respBody)
			resp2.Body.Close()

			if len(respBody.Errors) > 0 {
				return "", errors.New(respBody.Errors[0])
			}
		}

		return "", errors.New("Could not request authorization")
	}

	if resp2.Header == nil {
		return "", errors.New("Could not request authorization")
	}
	location, _ := url.Parse(resp2.Header.Get("Location"))
	ticket := location.Query().Get("ticket")

	authorizeForm := url.Values{}
	authorizeForm.Set("client_id", clientID)
	authorizeForm.Set("redirect_uri", redirectURI)
	authorizeForm.Set("client_secret", clientSecret)
	authorizeForm.Set("grant_type", "refresh_token")
	authorizeForm.Set("code", ticket)

	authorizeFormData := strings.NewReader(authorizeForm.Encode())

	req3, _ := http.NewRequest("POST", authorizeURL, authorizeFormData)
	req3.Header.Set("User-Agent", "pokemongo/1 CFNetwork/808.2.16 Darwin/16.3.0")
	req3.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp3, err3 := ctxhttp.Do(ctx, httpClient, req3)
	if err3 != nil {
		return "", errors.New("Could not authorize code")
	}

	b, _ := ioutil.ReadAll(resp3.Body)
	query, _ := url.ParseQuery(string(b))

	return query.Get("access_token"), nil
}
