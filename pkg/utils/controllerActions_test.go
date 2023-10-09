package utils

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.mathworks.com/development/mito/pkg/mwhttp"
	"github.mathworks.com/development/mito/pkg/mwhttptesttools"
	"go.uber.org/zap"
	"net/http"
	"testing"
)

const (
	testURL             = "/v1/getSignIn"
	testUsername        = "dummyUserName"
	testPassword        = "dummyPassword"
	testQueryParam      = "q"
	testQueryParamValue = "100"
)

type TestPostStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func setupTestRequestPayload() *bytes.Buffer {
	values := map[string]string{
		"username": testUsername,
		"password": testPassword,
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(values)
	return b
}

func TestCheckForHttpType(t *testing.T) {
	t.Run("Valid HTTP Type", func(t *testing.T) {
		req := mwhttptesttools.NewRequest(http.MethodPost, testURL, setupTestRequestPayload())
		_, _, err := CheckForHttpType(req, zap.L().Named("test-log-zap"))
		assert.NoError(t, err)
	})

	t.Run("Invalid HTTP Type", func(t *testing.T) {
		_, _, err := CheckForHttpType(map[string]string{}, zap.L().Named("test-log-zap"))
		assert.Error(t, err)
	})
}

func TestDispatchJsonResponse(t *testing.T) {
	t.Run("Valid JSON Response", func(t *testing.T) {
		data := map[string]string{"username": testUsername, "password": testPassword}
		m, err := DispatchJsonResponse(data, zap.L().Named("test-log-zap"), http.StatusOK)
		assert.NoError(t, err)
		sr := m.(mwhttp.SimpleResponse)
		assert.Equal(t, http.StatusOK, sr.Status)
	})

	t.Run("Invalid JSON Response", func(t *testing.T) {
		m, err := DispatchJsonResponse(make(chan int), zap.L().Named("test-log-zap"), http.StatusOK)
		assert.NoError(t, err)
		sr := m.(mwhttp.SimpleResponse)
		assert.Equal(t, http.StatusInternalServerError, sr.Status)
	})
}

func TestHttpMsgExtractor(t *testing.T) {
	t.Run("Valid POST Request", func(t *testing.T) {
		req := mwhttptesttools.NewRequest(http.MethodPost, testURL, setupTestRequestPayload())
		body := &TestPostStruct{}
		pack, err := HttpMsgExtractor(req, zap.L().Named("test-log-zap"), []string{http.MethodPost}, body)
		assert.NoError(t, err)
		assert.Equal(t, http.MethodPost, pack.Method)
		assert.Equal(t, testUsername, body.Username)
	})

	t.Run("Valid GET Request with Query Params", func(t *testing.T) {
		req := mwhttptesttools.NewRequest(http.MethodGet, testURL+"?"+testQueryParam+"="+testQueryParamValue, nil)
		body := &TestPostStruct{}
		pack, err := HttpMsgExtractor(req, zap.L().Named("test-log-zap"), []string{http.MethodPost, http.MethodGet}, body)
		assert.NoError(t, err)
		assert.Equal(t, http.MethodGet, pack.Method)
		val, exists := pack.QueryParams[testQueryParam]
		assert.True(t, exists)
		assert.Equal(t, testQueryParamValue, val[0])
	})

	t.Run("Unsupported HTTP Method", func(t *testing.T) {
		req := mwhttptesttools.NewRequest(http.MethodPost, testURL, nil)
		body := &TestPostStruct{}
		pack, err := HttpMsgExtractor(req, zap.L().Named("test-log-zap"), []string{http.MethodGet}, body)
		assert.Error(t, err)
		assert.Equal(t, http.MethodPost, pack.Method)
	})

	t.Run("Failed Payload Unmarshall", func(t *testing.T) {
		req := mwhttptesttools.NewRequest(http.MethodPost, testURL, setupTestRequestPayload())
		pack, err := HttpMsgExtractor(req, zap.L().Named("test-log-zap"), []string{http.MethodPost}, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.MethodPost, pack.Method)
		assert.Nil(t, pack.Body)
	})
}
