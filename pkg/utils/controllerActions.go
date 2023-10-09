package utils

import (
	"encoding/json"
	"errors"
	"github.mathworks.com/development/mito/pkg/core"
	"github.mathworks.com/development/mito/pkg/mwhttp"
	"github.mathworks.com/development/opi-utils-go/pkg/stringutils"
	"go.uber.org/zap"
	"net/http"
	"reflect"
)

/**
1. This is a support utility file- it is STATELESS!! , please refrain from adding any state.
2. It is garbage in >> garbage out
3. It can be used with any controller, remains free of opinions
*/

type HttpPacket struct {
	Request     mwhttp.Request
	Method      string
	Response    mwhttp.SimpleResponse
	Body        interface{}
	QueryParams map[string][]string
}

func HttpMsgExtractor(msg core.Message, log *zap.Logger, allowedMethods []string, bodyStruct interface{}) (*HttpPacket, error) {
	httpReq, packet, err := CheckForHttpType(msg, log)
	if err != nil {
		return packet, err
	}
	isValid, method := validateSupportedHttpMethods(httpReq, allowedMethods)
	if !isValid {
		log.Warn("Unsupported Http Method invoked: ", zap.Any("method", method), zap.Any("route", httpReq.Request.URL.Path))
		return &HttpPacket{
			Request:  httpReq,
			Method:   method,
			Response: mwhttp.NewSimpleResponseText(http.StatusMethodNotAllowed, "Unsupported Method"),
		}, errors.New("unsupported Http Method invoked")
	}
	if method == http.MethodPost || method == http.MethodDelete {
		je := json.NewDecoder(httpReq.Request.Body).Decode(&bodyStruct)
		if je != nil {
			return &HttpPacket{
				Request:  httpReq,
				Method:   method,
				Response: mwhttp.NewSimpleResponseText(http.StatusBadRequest, "Failed to parse request"),
			}, errors.New("bad Request Body")
		}
	}
	return &HttpPacket{
		Request:     httpReq,
		Method:      method,
		Response:    mwhttp.SimpleResponse{},
		QueryParams: httpReq.Request.URL.Query(),
	}, nil
}

func CheckForHttpType(msg core.Message, log *zap.Logger) (mwhttp.Request, *HttpPacket, error) {
	httpReq, ok := msg.(mwhttp.Request)
	if !ok {
		log.Error("Unexpected request message type", zap.Any("msg", reflect.TypeOf(msg)))
		return mwhttp.Request{}, &HttpPacket{
			Request:  mwhttp.Request{},
			Method:   "N/A",
			Response: mwhttp.NewSimpleResponseText(http.StatusInternalServerError, "Error"),
		}, errors.New("failed to unpack http")
	}
	return httpReq, nil, nil
}

func DispatchJsonResponse(data interface{}, log *zap.Logger, statusCode int) (core.Message, error) {
	dtStr, de := json.Marshal(data)
	if de != nil {
		log.Error("Failed to marshal the Json: ", zap.Any("data", reflect.TypeOf(data)))
		return mwhttp.NewSimpleResponseText(http.StatusInternalServerError, "Sorry, something went wrong, please check back in a while"), nil
	}
	return mwhttp.NewSimpleResponseContent(statusCode, "application/json", string(dtStr)), nil
}

func validateSupportedHttpMethods(request mwhttp.Request, allowedMethods []string) (bool, string) {
	method := request.Request.Method
	return stringutils.ContainsAny(method, allowedMethods), method
}
