package controllers

import (
	"fmt"
	"github.mathworks.com/development/mito/pkg/config"
	"github.mathworks.com/development/mito/pkg/core"
	"github.mathworks.com/development/mito/pkg/mwhttp"
	"github.mathworks.com/development/signindatatrackerws/pkg/collaborators"
	"github.mathworks.com/development/signindatatrackerws/pkg/domain"
	"github.mathworks.com/development/signindatatrackerws/pkg/utils"
	"go.uber.org/zap"
	"net/http"
)

const (
	InvalidUniqueIdMsg = "Invalid UniqueId"
	ParamUniqueID      = "uniqueId"
	ParamReferenceID   = "referenceId"
	ParamTimestamp     = "timestamp"
	ParamStartTime     = "startTime"
	ParamEndTime       = "endTime"
)

var RetrieveSignInDataControllerConstants = &ControllerMetaData{
	Name:            "retrieveSignInData",
	Path:            []string{"/v1/getUniqueSignIn", "/v1/getSignInDetails", "/v1/signInPeriodDetails", "/v1/signInReferenceId"},
	LoggerName:      "retrieveSignInData.controller",
	JsonContentType: "application/json",
	AllowedMethods:  []string{http.MethodGet},
}

func RetrieveSignInControllerFactory(conf config.Config, router mwhttp.Router, registry core.Registry) *RetrieveSignInDataController {
	controller := &RetrieveSignInDataController{
		logger:            zap.L().Named(RetrieveSignInDataControllerConstants.Name),
		signInDataService: collaborators.NewSignInTrackingService(),
	}
	registry.AddServiceProvider(RetrieveSignInDataControllerConstants.Name, controller, core.PublicRoute)
	for _, path := range RetrieveSignInDataControllerConstants.Path {
		router.AddRoute(path, RetrieveSignInDataControllerConstants.Name)
	}
	return controller
}

type RetrieveSignInDataController struct {
	logger            *zap.Logger
	signInDataService *collaborators.SignInTrackingService
}

func (rsdc RetrieveSignInDataController) Receive(message core.Message, ctx core.Context) (core.Message, error) {

	var ar = new(domain.RequestInput)
	packet, err := utils.HttpMsgExtractor(message, rsdc.logger, RetrieveSignInDataControllerConstants.AllowedMethods, &ar)

	if err != nil {
		return packet.Response, nil
	}

	var pathToHandler = map[string]func(map[string][]string) (core.Message, error){
		"/v1/getUniqueSignIn":     rsdc.handleUniqueSignIn,
		"/v1/getSignInDetails":    rsdc.handleGetSignInDetails,
		"/v1/signInPeriodDetails": rsdc.handleSignInPeriodDetails,
		"/v1/signInReferenceId":   rsdc.handleSignInReferenceId,
	}

	handler, ok := pathToHandler[packet.Request.Request.URL.Path]
	if ok {
		return handler(packet.QueryParams)
	}
	return nil, fmt.Errorf("invalid Path: %s", packet.Request.Request.URL.Path)
}

func (rsdc RetrieveSignInDataController) handleUniqueSignIn(packet map[string][]string) (core.Message, error) {
	uniqueID, _, timestamp, _, err := extractQueryParams(packet)
	if err != nil {
		return mwhttp.NewSimpleResponseText(http.StatusBadRequest, err.Error()), nil
	}

	requestInput := domain.RequestInput{
		UniqueID:  uniqueID,
		Timestamp: timestamp,
	}

	pd, errResp, statusCode := rsdc.signInDataService.FindUniqueSignInInfo(requestInput)
	if errResp.ErrorCode != 0 {
		return utils.DispatchJsonResponse(errResp, rsdc.logger, statusCode)
	}

	return utils.DispatchJsonResponse(pd, rsdc.logger, http.StatusOK)
}

func (rsdc RetrieveSignInDataController) handleSignInPeriodDetails(packet map[string][]string) (core.Message, error) {
	uniqueID, _, startTime, endTime, err := extractQueryParams(packet)
	if err != nil {
		return mwhttp.NewSimpleResponseText(http.StatusBadRequest, err.Error()), nil
	}

	requestDetailsInput := domain.RequestTimestampInput{
		UniqueID:  uniqueID,
		StartTime: startTime,
		EndTime:   endTime,
	}

	pd, errResp, statusCode := rsdc.signInDataService.FindSignInPeriodDetails(requestDetailsInput)
	if errResp.ErrorCode != 0 {
		return utils.DispatchJsonResponse(errResp, rsdc.logger, statusCode)
	}

	return utils.DispatchJsonResponse(pd, rsdc.logger, http.StatusOK)
}
func (rsdc RetrieveSignInDataController) handleSignInReferenceId(packet map[string][]string) (core.Message, error) {
	uniqueID, referenceId, _, _, err := extractQueryParams(packet) // include referenceId here
	if err != nil {
		return mwhttp.NewSimpleResponseText(http.StatusBadRequest, err.Error()), nil
	}

	requestDetailsInput := domain.RequestReferenceIdInput{
		UniqueID:    uniqueID,
		ReferenceId: referenceId,
	}

	pd, errResp, statusCode := rsdc.signInDataService.FindSignInReferenceIds(requestDetailsInput)
	if errResp.ErrorCode != 0 {
		return utils.DispatchJsonResponse(errResp, rsdc.logger, statusCode)
	}

	return utils.DispatchJsonResponse(pd, rsdc.logger, http.StatusOK)
}

func (rsdc RetrieveSignInDataController) handleGetSignInDetails(packet map[string][]string) (core.Message, error) {
	uniqueID, _, _, _, err := extractQueryParams(packet)
	if err != nil {
		return mwhttp.NewSimpleResponseText(http.StatusBadRequest, err.Error()), nil
	}

	requestDetailsInput := domain.RequestDetailsInput{
		UniqueID: uniqueID,
	}

	pd, errResp, statusCode := rsdc.signInDataService.FindSignInTrackingDetails(requestDetailsInput)
	if errResp.ErrorCode != 0 {
		return utils.DispatchJsonResponse(errResp, rsdc.logger, statusCode)
	}

	return utils.DispatchJsonResponse(pd, rsdc.logger, http.StatusOK)
}

func extractQueryParams(queryParams map[string][]string) (uniqueID string, referenceId string, startTime string, endTime string, err error) {
	uniqueID = extractQueryParamHelper(queryParams, ParamUniqueID)
	if uniqueID == "" {
		return "", "", "", "", fmt.Errorf(InvalidUniqueIdMsg)
	}
	referenceId = extractQueryParamHelper(queryParams, ParamReferenceID)
	startTime = extractQueryParamHelper(queryParams, ParamTimestamp)
	if startTime == "" {
		startTime = extractQueryParamHelper(queryParams, ParamStartTime)
	}
	endTime = extractQueryParamHelper(queryParams, ParamEndTime)
	return
}

func extractQueryParamHelper(queryParams map[string][]string, param string) string {
	if val, ok := queryParams[param]; ok && len(val) > 0 {
		return val[0]
	}
	return ""
}
