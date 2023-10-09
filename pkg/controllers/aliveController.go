package controllers

import (
	"net/http"
	"strings"

	"github.mathworks.com/development/mito/pkg/config"
	"github.mathworks.com/development/mito/pkg/core"
	"github.mathworks.com/development/mito/pkg/mwhttp"
	"github.mathworks.com/development/opi-utils-go/pkg/configutils"
	"github.mathworks.com/development/signindatatrackerws/pkg/domain"
	"github.mathworks.com/development/signindatatrackerws/pkg/utils"
	"go.uber.org/zap"
)

var AliveControllerConstants = &ControllerMetaData{
	Name:            "aliveEndpoint",
	Path:            []string{"/v1/admin/alive"},
	LoggerName:      "alive.controller",
	JsonContentType: "application/json",
	AllowedMethods:  []string{http.MethodGet},
}

func AliveControllerFactory(conf config.Config, router mwhttp.Router, registry core.Registry) *AliveController {
	prepAlive()
	controller := &AliveController{
		logger: zap.L().Named(AliveControllerConstants.Name),
	}
	registry.AddServiceProvider(AliveControllerConstants.Name, controller, core.PublicRoute)
	router.AddRoute(AliveControllerConstants.Path[0], AliveControllerConstants.Name)
	router.AddRoute("/v1/admin/alive/{echo}", AliveControllerConstants.Name)
	return controller
}

var Alive *domain.Alive

func prepAlive() {
	props, _ := configutils.ReadPropFile("configfiles/alive.properties")
	alive := new(domain.Alive)
	alive.Version = utils.GetValueFromMap(props, "version", "no-clue")
	alive.ArtifactId = utils.GetValueFromMap(props, "artifactId", "no-clue")
	alive.GroupId = utils.GetValueFromMap(props, "groupId", "no-clue")
	alive.BuildTimestamp = utils.GetValueFromMap(props, "buildTimestamp", "no-clue")
	alive.Echo = "Hello"
	Alive = alive
}

type AliveController struct {
	logger *zap.Logger
}

func (a AliveController) Receive(message core.Message, ctx core.Context) (core.Message, error) {
	var ar = new(domain.MonoResponse)
	packet, err := utils.HttpMsgExtractor(message, a.logger, AliveControllerConstants.AllowedMethods, &ar)
	if err != nil {
		return packet.Response, nil
	}
	pathParams := strings.Split(packet.Request.Request.URL.Path, "/")
	if pathParams[len(pathParams)-1] == "alive" {
		Alive.Echo = "Hello"
	} else {
		Alive.Echo = pathParams[len(pathParams)-1]
	}
	return utils.DispatchJsonResponse(Alive, a.logger, http.StatusOK)
}
