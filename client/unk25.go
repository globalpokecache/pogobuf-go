package client

import (
	"errors"
)

func getUnk25(version int) (int64, error) {
	switch version {
	case 5704:
		return -816976800928766045, nil
	case 5703:
		return -816976800928766045, nil
	case 5702:
		return -816976800928766045, nil
	case 5500:
		return -9156899491064153954, nil
	case 5300:
		return -8832040574896607694, nil
	case 5100:
		return -76506539888958491, nil
	case 4500:
		return -8408506833887075802, nil
	default:
		return 0, errors.New("Unsupported API version")
	}
}
