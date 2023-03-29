package pattern

import (
	"fmt"
	"testing"
)

func TestFacadeWork(t *testing.T) {
	err := FacadeWork()
	if err != nil {
		fmt.Println(err)
	}
}
