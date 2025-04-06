package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mcdotjs/chirpy/internal/auth"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	id := uuid.New()
	token, err := auth.MakeJWT(id, "secret", time.Minute)
	if err != nil {

		t.Errorf("makejwt %s", err)
	}
	fmt.Println("TOKEN: ", token)
	data, err := auth.ValidateJWT(token, "secret")
	if err != nil {
		t.Errorf("validate %s", err)
	}

	fmt.Println("DDDDD: ", data)
	t.Log(data)
}
