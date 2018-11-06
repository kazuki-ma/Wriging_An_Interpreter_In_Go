package main

import (
	"log"
	"os/user"
	"testing"
)

func TestGetCurrentUser(t *testing.T) {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	log.Printf("user = %s", user)

	t.Run("Name", func(t *testing.T) {
		if user == nil {
			t.Errorf("user must be NON-nil, but got nil")
		}

		t.Run("v2", func(t *testing.T) {

		})
	})
}
