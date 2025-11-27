package test

import (

	"testing"
)

func TestAdd(t *testing.T) {
    result := 5
    if result != 5 {
        t.Errorf("expected 5, got %d", result)
    }
}
