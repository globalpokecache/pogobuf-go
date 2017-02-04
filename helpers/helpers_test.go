package helpers

import (
	"fmt"
	"testing"
)

func TestGetCellsFromRadius(t *testing.T) {
	fmt.Println(GetCellsFromRadius(-16.7034135, -49.2385342, 210, 17))
}
