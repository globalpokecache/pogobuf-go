package buddyauth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Provider struct {
	version   int
	keysMutex sync.RWMutex
	keys      []*BuddyKey
}

type BuddyKey struct {
	Key       string
	RPM       int
	Used      int
	NextReset time.Time
	Invalid   bool
	Expired   bool
}

type HashRequest struct {
	Timestamp   uint64   `json:"timestamp"`
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
	Altitude    float64  `json:"altitude"`
	AuthTicket  string   `json:"authTicket"`
	SessionData string   `json:"sessionData"`
	Requests    []string `json:"requests"`
}

type HashResponse struct {
	LocationAuthHash uint32  `json:"locationAuthHash"`
	LocationHash     uint32  `json:"locationHash"`
	RequestHashes    []int64 `json:"RequestHashes"`
}

var (
	ErrFailedToRequest = errors.New("Failed to request hash server")
	ErrKeyPassedLimit  = errors.New("Key passed the limit")
	ErrBadRequest      = errors.New("Something wrong in the request")
	ErrInvalidKey      = errors.New("Invalid or expired key")
	ErrNoAvailableKey  = errors.New("No available key")
)

var (
	versions map[string]string
	Debug    bool
)

func NewProvider(apiVersion int) (*Provider, error) {
	resp, err := http.Get("https://pokehash.buddyauth.com/api/hash/versions")
	if err != nil {
		return nil, errors.New("Failed to load buddyauth versions")
	}

	err = json.NewDecoder(resp.Body).Decode(&versions)
	if err != nil {
		return nil, errors.New("Failed to load buddyauth versions")
	}

	provider := &Provider{
		version: apiVersion,
		keys:    []*BuddyKey{},
	}

	return provider, nil
}

type byRPM []*BuddyKey

func (a byRPM) Len() int           { return len(a) }
func (a byRPM) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byRPM) Less(i, j int) bool { return a[i].RPM < a[j].RPM }

func (p *Provider) AddKey(key string) error {
	p.keysMutex.Lock()
	defer p.keysMutex.Unlock()

	p.keys = append(p.keys, &BuddyKey{
		Key:  key,
		RPM:  150, // Not defined yet
		Used: 0,
	})

	sort.Sort(byRPM(p.keys))

	return nil
}

func (p *Provider) ApiURL() string {
	v := fmt.Sprintf("0.%.1f", float64(p.version)/100)
	return fmt.Sprintf("http://hashing.pogodev.io/%s", versions[v])
}

func (p *Provider) GetAvailableKey() (BuddyKey, error) {
	p.keysMutex.Lock()
	defer p.keysMutex.Unlock()

	var key *BuddyKey
	var found bool
	debug("Searching for available key")
	for i := 0; i < len(p.keys); i++ {
		key = p.keys[0]

		if key.NextReset.Before(time.Now()) {
			debug("Resetting key: %s", key.Key)
			key.NextReset = time.Now()
			key.Used = 0
		}

		if !key.Expired && !key.Invalid && key.Used < key.RPM {
			debug("Found valid key: %s", key.Key)
			found = true
			break
		}

		debug("Skipping invalid key: %s", key.Key)
	}
	if !found {
		debug("No valid key found")
		return BuddyKey{}, ErrNoAvailableKey
	}

	key.Used++
	return *key, nil
}

func (p *Provider) ReturnKey(k *BuddyKey) {
	p.keysMutex.Lock()
	defer p.keysMutex.Unlock()
	for i, key := range p.keys {
		if key.Key == k.Key {
			p.keys[i] = k
			break
		}
	}
}

func (p *Provider) hashRequest(hashReq HashRequest, key *BuddyKey) (HashResponse, error) {
	var hresp HashResponse

	requestBytes, err := json.Marshal(&hashReq)
	if err != nil {
		return hresp, fmt.Errorf("Failed to marshal hash request: %s", err)
	}

	data := bytes.NewReader(requestBytes)

	debug("Sending request to: %s", p.ApiURL())
	req, err := http.NewRequest("POST", p.ApiURL(), data)
	if err != nil {
		return hresp, fmt.Errorf("Failed to create request: %s", err)
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-AuthToken", key.Key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return hresp, fmt.Errorf("Failed to do request: %s", err)
	}

	if resp.Header.Get("x-maxrequestcount") != "" {
		rateperiodend, _ := strconv.ParseInt(resp.Header.Get("x-rateperiodend"), 10, 64)
		maxrequestcount, _ := strconv.ParseInt(resp.Header.Get("x-maxrequestcount"), 10, 64)
		raterequestsremaining, _ := strconv.ParseInt(resp.Header.Get("x-raterequestsremaining"), 10, 64)
		authtokenexpiration, _ := strconv.ParseInt(resp.Header.Get("x-authtokenexpiration"), 10, 64)
		debug("Updating key info: %s", key.Key)

		if time.Unix(authtokenexpiration, -1).Before(time.Now()) {
			key.Expired = true
		}

		key.NextReset = time.Unix(rateperiodend, -1)
		key.RPM = int(maxrequestcount)
		key.Used = int(maxrequestcount - raterequestsremaining)
	}

	debug("Server response code: %s", resp.Status)
	switch resp.StatusCode {
	case http.StatusBadRequest, http.StatusNotFound:
		return hresp, ErrBadRequest
	case http.StatusUnauthorized:
		key.Invalid = true
		return hresp, ErrInvalidKey
	case 429:
		debug("Key passed limit: %s", key.Key)
		return hresp, ErrKeyPassedLimit
	}

	err = json.NewDecoder(resp.Body).Decode(&hresp)
	if err != nil {
		return hresp, fmt.Errorf("Failed to decode hash server response: %s", err)
	}

	return hresp, nil
}

func (p *Provider) Hash(authTicket, sessionData []byte, latitude, longitude, accuracy float64, timestamp uint64, requests [][]byte) (uint32, uint32, []uint64, error) {
	baseAuthTicket := base64.StdEncoding.EncodeToString(authTicket)
	baseSessionData := base64.StdEncoding.EncodeToString(sessionData)

	var reqBases []string
	for _, b := range requests {
		reqBases = append(reqBases, base64.StdEncoding.EncodeToString(b))
	}

	hashReq := HashRequest{
		Timestamp:   timestamp,
		AuthTicket:  baseAuthTicket,
		SessionData: baseSessionData,
		Latitude:    latitude,
		Longitude:   longitude,
		Altitude:    accuracy,
		Requests:    reqBases,
	}

	var err error
	var hashResp HashResponse
	var key BuddyKey

	var success bool
	for i := 0; i < len(p.keys); i++ {
		key, err = p.GetAvailableKey()
		debug("Found key: %s", key.Key)
		if err != nil {
			return 0, 0, []uint64{0}, err
		}
		if Debug {
			d, _ := json.MarshalIndent(hashReq, "", "\t")
			debug("Sending hash request: %s", d)
		}
		hashResp, err = p.hashRequest(hashReq, &key)
		p.ReturnKey(&key)
		if err == nil {
			success = true
			if Debug {
				d, _ := json.MarshalIndent(hashResp, "", "\t")
				debug("Valid response: %s", string(d))
			}
			break
		}
		debug("Failed to hash request: %s", err)
		time.Sleep(1 * time.Second)
	}

	if !success {
		return 0, 0, []uint64{}, err
	}

	var reqHashes = make([]uint64, len(hashResp.RequestHashes))
	for i, hash := range hashResp.RequestHashes {
		reqHashes[i] = uint64(hash)
	}

	return hashResp.LocationAuthHash, hashResp.LocationHash, reqHashes, nil
}

func (p *Provider) SetDebug(d bool) {
	Debug = d
}
