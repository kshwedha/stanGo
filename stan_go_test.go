package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// cmd to run specific testcase
// go test -timeout 30s -run ^TestDeleteOp$ stan_go

func TestVerifyUpdate(t *testing.T) {
	user := CreateUserData()
	assert.GreaterOrEqual(t, len(user), 1)
}

func TestDeleteOp(t *testing.T) {
	Data = make(map[string]User)
	user := CreateUserData()
	for _, data := range user {
		Data[data.Username] = data
	}
	doDeleteOperation()
	assert.Less(t, len(Data), len(user))
}

func TestUpdateUser(t *testing.T) {
	Data = make(map[string]User)
	user := CreateUserData()
	for _, data := range user {
		Data[data.Username] = data
	}
	randomUser := Data[fetchUserRecord()]
	email := randomUser.Email
	randomUser.update()
	assert.NotEqual(t, randomUser.Email, email)
}
