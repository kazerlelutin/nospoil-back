package utils_test

import (
	utils "k-space-go/utils"
	"testing"
)

func TestAdd(t *testing.T) {
	result := utils.Add(2, 3)
	// Test de la fonction Add
	if result != 5 {
		t.Errorf("Add(2, 3) failed, expected %d, got %d", 5, result)
	}
}
