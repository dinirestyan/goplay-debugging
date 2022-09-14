package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dinirestyan/goplay-debugging/models"
	"gotest.tools/assert"
)

func TestLogin(t *testing.T) {
	login := models.Login{
		Email:    "usertest@email.com",
		Password: "s!mp4n",
	}
	jsonValue, _ := json.Marshal(login)
	req, _ := http.NewRequest("POST", "/goplay/login", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	request, _ := json.Marshal(req)
	fmt.Printf(string(request))
	assert.Equal(t, http.StatusOK, w.Code)
}
