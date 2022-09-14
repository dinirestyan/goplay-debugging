package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock_models "github.com/dinirestyan/goplay-debugging/mock"
	"github.com/dinirestyan/goplay-debugging/utils"

	"github.com/dinirestyan/goplay-debugging/models"
	"github.com/gin-gonic/gin"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUsecaseMock := mock_models.NewMockGoplayRepo(ctrl)

	testCases := []struct {
		name,
		Request,
		wantResult string
		mock func()
	}{
		{
			name:       "Success",
			Request:    `{"email":"dummy@dummy.com","password":"dummy"}`,
			wantResult: `{"meta":{"status":false,"code":0,"message":"logged in"},"data":{"id":"00000000-0000-0000-0000-000000000000","created_at":"0001-01-01T00:00:00Z","name":"dummyName","email":"dummy@dummy.com","token":""}}`,
			mock: func() {
				resp := models.User{
					Email: "dummy@dummy.com",
					Name:  "dummyName",
				}
				userUsecaseMock.EXPECT().Login("dummy@dummy.com").Return(&resp, nil)
			},
		},
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery()) // Recovery middleware recovers from any panics and writes a 500 if there was one.
	db := utils.GetDBConnection()
	NewGoplayController(db, userUsecaseMock)

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/signin", strings.NewReader(test.Request))
			r.ServeHTTP(w, req)
			assert.Equal(t, test.wantResult, w.Body.String())
		})
	}
}

func TestCreateResponse(t *testing.T) {
	// Diabaikan func test ini, gue buat karena males ngetik json responsenya :)
	resp := models.User{
		Email: "dummy@dummy.com",
		Name:  "dummyName",
	}
	jsonRsp, _ := json.Marshal(utils.Response(http.StatusOK, "logged in", resp))
	log.Println(string(jsonRsp))
}
