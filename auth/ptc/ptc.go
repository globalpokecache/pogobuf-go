package ptc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/globalpokecache/pogobuf-go"
	"golang.org/x/net/context/ctxhttp"
)

const authorizeURL = "https://sso.pokemon.com/sso/oauth2.0/accessToken"
const loginURL = "https://sso.pokemon.com/sso/login?locale=en&service=https://sso.pokemon.com/sso/oauth2.0/callbackAuthorize"

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
	username, password string
	debug              bool
}

// NewProvider constructs a Pokémon Trainer's Club auth provider instance
func NewProvider(username, password string) *Provider {
	return &Provider{username, password, false}
}

func (p *Provider) Type() string {
	return "ptc"
}

func (p *Provider) SetDebug(d bool) {
	p.debug = d
}

func (p *Provider) GetUsername() string {
	return p.username
}

// Login retrieves an access token from the Pokémon Trainer's Club
func (p *Provider) Login(ctx context.Context) (string, error) {
	options := &cookiejar.Options{}
	jar, _ := cookiejar.New(options)
	httpClient := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("Use the last error")
		},
	}

	req1, _ := http.NewRequest("GET", loginURL, nil)
	req1.Header.Set("User-Agent", "niantic")

	p.DebugMsg("Requesting: %s", loginURL)
	resp1, err1 := ctxhttp.Do(ctx, httpClient, req1)
	if err1 != nil {
		return "", errors.New("Failed to connect to login servers... Probably server issues")
	}

	defer resp1.Body.Close()
	body1, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		return "", errors.New("Failed to connect to login servers... Probably server issues")
	}
	var loginRespBody loginRequest
	err = json.Unmarshal(body1, &loginRespBody)
	if err != nil {
		return "", errors.New("Failed to connect to login servers... Probably server issues")
	}
	p.Debug(body1)
	resp1.Body.Close()

	loginForm := url.Values{}
	loginForm.Set("lt", loginRespBody.Lt)
	loginForm.Set("execution", loginRespBody.Execution)
	loginForm.Set("_eventId", "submit")
	loginForm.Set("username", p.username)
	loginForm.Set("password", p.password)

	loginFormData := strings.NewReader(loginForm.Encode())

	req2, _ := http.NewRequest("POST", loginURL, loginFormData)
	req2.Header.Set("User-Agent", "niantic")
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	p.DebugMsg("Requesting: %s", loginURL)
	resp2, err2 := ctxhttp.Do(ctx, httpClient, req2)
	if _, ok2 := err2.(*url.Error); !ok2 {
		if resp2 != nil && resp2.Body != nil {
			defer resp2.Body.Close()
			body2, _ := ioutil.ReadAll(resp2.Body)
			var respBody loginRequest
			json.Unmarshal(body2, &respBody)
			if strings.Contains(respBody.Errors[0], "unexpected error") {
				return "", pogobuf.ErrAccountBanned
			}
			p.Debug(respBody)
			return "", errors.New("Failed to connect to login servers... Probably server issues")
		}

		return "", errors.New("Could not request authorization")
	}

	if resp2 == nil || resp2.Header == nil {
		return "", errors.New("Could not request authorization")
	}
	location, _ := url.Parse(resp2.Header.Get("Location"))
	ticket := location.Query().Get("ticket")
	p.Debug(location)

	authorizeForm := url.Values{}
	authorizeForm.Set("client_id", clientID)
	authorizeForm.Set("redirect_uri", redirectURI)
	authorizeForm.Set("client_secret", clientSecret)
	authorizeForm.Set("grant_type", "refresh_token")
	authorizeForm.Set("code", ticket)

	authorizeFormData := strings.NewReader(authorizeForm.Encode())

	req3, _ := http.NewRequest("POST", authorizeURL, authorizeFormData)
	req3.Header.Set("User-Agent", "niantic")
	req3.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	p.DebugMsg("Requesting: %s", authorizeURL)
	p.Debug(authorizeForm)
	resp3, err3 := ctxhttp.Do(ctx, httpClient, req3)
	if err3 != nil {
		return "", errors.New("Failed to connect to login servers...")
	}

	if resp3 == nil || err3 != nil {
		return "", errors.New("Failed to connect to login servers...")
	}

	if resp3.StatusCode != http.StatusOK {
		return "", errors.New("Failed to connect to login servers... Probably server issues")
	}

	b, _ := ioutil.ReadAll(resp3.Body)
	query, err := url.ParseQuery(string(b))

	if err != nil {
		return "", errors.New("Failed to connect to login servers... Probably server issues")
	}

	p.Debug(query)

	return query.Get("access_token"), nil
}

func (p *Provider) Debug(m interface{}) {
	if p.debug && m != nil {
		b, _ := json.MarshalIndent(m, "", "\t")
		p.DebugMsg("%s", string(b))
	}
}

func (p *Provider) DebugMsg(format string, a ...interface{}) {
	if p.debug {
		fmt.Printf(fmt.Sprintf("(PTC) %s\n", format), a...)
	}
}
