package random

import (
	"fmt"
	"testing"
)

func TestRandom(t *testing.T) {
	fmt.Println(NewRandomString(7))
}
