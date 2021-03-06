package buddyauth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
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
	http.Client
}

type HashRequest struct {
	Timestamp   uint64   `json:"timestamp"`
	Latitude    int64    `json:"Latitude64"`
	Longitude   int64    `json:"Longitude64"`
	Accuracy    int64    `json:"Accuracy64"`
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
	ErrNoAvailableKey  = errors.New("No hash key is available")
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
		Client: http.Client{
			Timeout: 10 * time.Second,
		},
	}

	return provider, nil
}

type BuddyKey struct {
	Key       string
	RPM       int
	Used      int
	NextReset time.Time
	Invalid   bool
	Expired   bool
	sync.RWMutex
}

func (b *BuddyKey) GetRPM() int {
	b.RLock()
	defer b.RUnlock()
	return b.RPM
}

func (b *BuddyKey) SetRPM(n int) {
	b.RLock()
	b.RPM = n
	b.RUnlock()
}

func (b *BuddyKey) GetUsed() int {
	b.RLock()
	defer b.RUnlock()
	return b.Used
}

func (b *BuddyKey) AddUsed(n int) {
	b.RLock()
	b.Used += n
	b.RUnlock()
}

func (b *BuddyKey) ResetUsed() {
	b.Lock()
	b.Used = 0
	b.Unlock()
}

func (b *BuddyKey) IsExpired() bool {
	b.RLock()
	defer b.RUnlock()
	return b.Expired
}

func (b *BuddyKey) SetExpired(v bool) {
	b.RLock()
	b.Expired = v
	b.RUnlock()
}

func (b *BuddyKey) IsInvalid() bool {
	b.RLock()
	defer b.RUnlock()
	return b.Invalid
}

func (b *BuddyKey) SetInvalid(v bool) {
	b.RLock()
	b.Invalid = v
	b.RUnlock()
}

func (b *BuddyKey) GetNextReset() time.Time {
	b.RLock()
	defer b.RUnlock()
	return b.NextReset
}

func (b *BuddyKey) SetNextReset(t time.Time) {
	b.RLock()
	b.NextReset = t
	b.RUnlock()
}

type byRPM []*BuddyKey

func (a byRPM) Len() int      { return len(a) }
func (a byRPM) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byRPM) Less(i, j int) bool {
	if a[i].GetRPM() > a[j].GetRPM() {
		return true
	} else if a[i].GetRPM() == a[j].GetRPM() {
		return a[i].GetUsed() > a[j].GetUsed()
	}
	return false
}

func (p *Provider) GetKeys() []interface{} {
	p.keysMutex.Lock()
	defer p.keysMutex.Unlock()

	keys := []interface{}{}
	for _, key := range p.keys {
		key.Lock()
		keys = append(keys, key)
		key.Unlock()
	}
	return keys
}

func (p *Provider) AddKey(key string) error {
	p.keysMutex.Lock()
	defer p.keysMutex.Unlock()

	index := -1
	for i, k := range p.keys {
		if k.Key == key {
			index = i
			break
		}
	}
	if index != -1 {
		return errors.New("Key already added")
	}

	p.keys = append(p.keys, &BuddyKey{
		Key:  key,
		RPM:  150, // Not defined yet
		Used: 0,
	})

	sort.Sort(byRPM(p.keys))

	return nil
}

func (p *Provider) DelKey(key string) error {
	p.keysMutex.Lock()
	defer p.keysMutex.Unlock()

	index := -1
	for i, k := range p.keys {
		if k.Key == key {
			index = i
			break
		}
	}
	if index == -1 {
		return errors.New("Key not found")
	}
	p.keys = append(p.keys[:index], p.keys[index+1:]...)

	sort.Sort(byRPM(p.keys))

	return nil
}

func (p *Provider) ApiURL() string {
	major := math.Floor(float64(p.version) / 100.0)
	minor := (float64(p.version)/100.0 - major) * 100
	v := fmt.Sprintf("0.%.f.%.f", major, minor)
	return fmt.Sprintf("http://pokehash.buddyauth.com/%s", versions[v])
}

func (p *Provider) GetAvailableKey() (*BuddyKey, error) {
	p.keysMutex.Lock()
	defer p.keysMutex.Unlock()

	var key *BuddyKey
	var found bool
	debug("Searching for available key")
	var errs []string
	for i := 0; i < len(p.keys); i++ {
		key = p.keys[i]

		if key.NextReset.Before(time.Now().UTC()) {
			debug("Resetting key: %s", key.Key)
			key.SetNextReset(time.Now().UTC().Add(1 * time.Minute))
			key.ResetUsed()
		}

		if !key.IsExpired() && !key.IsInvalid() && key.GetUsed() < key.GetRPM() {
			debug("Found valid key: %s", key.Key)
			found = true
			break
		}

		if key.IsExpired() {
			errs = append(errs, fmt.Sprintf("%s:%s", key.Key[0:6]+"..."+key.Key[len(key.Key)-6:len(key.Key)-1], "EXPIRED_KEY"))
		}

		if key.IsInvalid() {
			errs = append(errs, fmt.Sprintf("%s:%s", key.Key[0:6]+"..."+key.Key[len(key.Key)-6:len(key.Key)-1], "INVALID_KEY"))
		}

		if key.GetUsed() >= key.GetRPM() {
			errs = append(errs, fmt.Sprintf("%s:%s", key.Key[0:6]+"..."+key.Key[len(key.Key)-6:len(key.Key)-1], "EXHAUSTED_KEY"))
		}

		debug("Skipping invalid key: %s", key.Key)
	}
	if !found {
		debug("No valid key found")
		return nil, fmt.Errorf("%s: %v", ErrNoAvailableKey, errs)
	}

	key.AddUsed(1)
	return key, nil
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

	req.Close = true
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-AuthToken", key.Key)

	resp, err := p.Client.Do(req)
	if err != nil {
		return hresp, fmt.Errorf("Failed to do request: %s", err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("X-Maxrequestcount") != "" {
		// rateperiodend, _ := strconv.ParseInt(resp.Header.Get("X-Rateperiodend"), 10, 64)
		maxrequestcount, _ := strconv.ParseInt(resp.Header.Get("X-Maxrequestcount"), 10, 64)
		remaining, _ := strconv.ParseInt(resp.Header.Get("X-Raterequestsremaining"), 10, 64)
		authtokenexpiration, _ := strconv.ParseInt(resp.Header.Get("X-Authtokenexpiration"), 10, 64)
		debug("Updating key info: %s", key.Key)
		debug("Received header:", resp.Header)

		if time.Unix(authtokenexpiration, -1).Before(time.Now()) {
			key.SetExpired(true)
		}

		if key.GetRPM() != int(maxrequestcount) {
			key.SetRPM(int(maxrequestcount))
		}

		if key.GetUsed() != int(remaining) {
			used := key.GetRPM() - int(remaining)

			if used > key.GetUsed() {
				key.ResetUsed()
				key.AddUsed(used)
			}
		}
	}

	debug("Server response code: %s", resp.Status)
	switch resp.StatusCode {
	case http.StatusBadRequest, http.StatusNotFound:
		return hresp, ErrBadRequest
	case http.StatusUnauthorized:
		key.SetInvalid(true)
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

	reqBases := []string{}
	for _, b := range requests {
		reqBases = append(reqBases, base64.StdEncoding.EncodeToString(b))
	}

	hashReq := HashRequest{
		Timestamp:   timestamp,
		AuthTicket:  baseAuthTicket,
		SessionData: baseSessionData,
		Latitude:    int64(math.Float64bits(latitude)),
		Longitude:   int64(math.Float64bits(longitude)),
		Accuracy:    int64(math.Float64bits(accuracy)),
		Requests:    reqBases,
	}

	var err error
	var hashResp HashResponse
	var key *BuddyKey

	var success bool
	for i := 0; i < len(p.keys); i++ {
		key, err = p.GetAvailableKey()
		if err != nil {
			return 0, 0, []uint64{0}, err
		}
		debug("Found key: %s", key.Key)
		if Debug {
			d, _ := json.MarshalIndent(hashReq, "", "\t")
			debug("Sending hash request: %s", d)
		}
		hashResp, err = p.hashRequest(hashReq, key)
		if err == nil {
			success = true
			if Debug {
				d, _ := json.MarshalIndent(hashResp, "", "\t")
				debug("Valid response: %s", string(d))
			}
			break
		}
		key.AddUsed(-1)
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
