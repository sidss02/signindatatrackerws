package controllers

import (
	"net/http"

	"github.mathworks.com/development/mito/pkg/config"
	"github.mathworks.com/development/mito/pkg/core"
	"github.mathworks.com/development/mito/pkg/mwhttp"
	"github.mathworks.com/development/signindatatrackerws/pkg/collaborators"
	"github.mathworks.com/development/signindatatrackerws/pkg/domain"
	"github.mathworks.com/development/signindatatrackerws/pkg/utils"
	"go.uber.org/zap"
)

const (
	ControllerName         = "saveSignInData"
	ErrorCodeEmptyUniqueID = 5300
)

var PersistSignInControllerConstants = &ControllerMetaData{
	Name:            ControllerName,
	Path:            []string{"/v1/saveSignInData"},
	LoggerName:      "saveSignInData.controller",
	JsonContentType: "application/json",
	AllowedMethods:  []string{http.MethodPost},
}

func PersistSignInControllerFactory(conf config.Config, router mwhttp.Router, registry core.Registry) *PersistSignInDataController {
	controller := &PersistSignInDataController{
		logger:            zap.L().Named(PersistSignInControllerConstants.Name),
		signInDataService: collaborators.NewSignInTrackingService(),
	}
	registry.AddServiceProvider(PersistSignInControllerConstants.Name, controller, core.PublicRoute)
	router.AddRoute(PersistSignInControllerConstants.Path[0], PersistSignInControllerConstants.Name)
	return controller
}

type PersistSignInDataController struct {
	logger            *zap.Logger
	signInDataService *collaborators.SignInTrackingService
}

func (gl PersistSignInDataController) Receive(message core.Message, ctx core.Context) (core.Message, error) {

	var ar = new(domain.SaveSignInInfo)

	packet, err := utils.HttpMsgExtractor(message, gl.logger, PersistSignInControllerConstants.AllowedMethods, &ar)
	if err != nil {
		return packet.Response, nil
	}
	if packet.Method == http.MethodPost {
		//TODO: this is where the controller logic goes
		signindata, errResp, statusCode := gl.signInDataService.SaveSignInData(*ar)
		if (errResp != domain.ErrorResponse{}) {
			if errResp.ErrorCode == ErrorCodeEmptyUniqueID {
				gl.logger.Error("UniqueId was empty")
			}
			return utils.DispatchJsonResponse(errResp, gl.logger, statusCode)
		}
		return utils.DispatchJsonResponse(signindata, gl.logger, http.StatusCreated)
	}
	return utils.DispatchJsonResponse(ar, gl.logger, http.StatusNoContent)
}
