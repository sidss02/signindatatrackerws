package controllers

import (
	"bytes"
	"fmt"
	healthAlive "github.mathworks.com/development/alive-go/pkg/alive"
	"github.mathworks.com/development/mito/pkg/config"
	"github.mathworks.com/development/mito/pkg/core"
	"github.mathworks.com/development/mito/pkg/mwhttp"
	"github.mathworks.com/development/signindatatrackerws/pkg/bootstrap"
	"github.mathworks.com/development/signindatatrackerws/pkg/collaborators"
	"github.mathworks.com/development/signindatatrackerws/pkg/domain"
	"github.mathworks.com/development/signindatatrackerws/pkg/utils"

	"go.uber.org/zap"
	"net/http"
	"strings"
	"text/template"
)

const (
	// DefaultPathDelimiter Adding these constants for better readability.
	DefaultPathDelimiter       = "/"
	HTMLResponseType           = "text/html"
	InvalidRequestBodyErrorMsg = "Request Body Empty"
	DefaultHTTPErrorCode       = 500
	InvalidKeyHTTPErrorCode    = 401
	DefaultHTTPSuccessCode     = 200
	SignInTrackerTable         = "signindatatracker"
)

var HealthControllerConstants = &ControllerMetaData{
	Name:            "healthEndpoint",
	Path:            []string{"/admin/health/v2"},
	LoggerName:      "health.controller",
	JsonContentType: "application/json",
	AllowedMethods:  []string{http.MethodGet, http.MethodPost},
}

func HealthControllerFactory(conf config.Config, router mwhttp.Router, registry core.Registry) *HealthController {
	controller := &HealthController{
		logger: zap.L().Named(HealthControllerConstants.Name),
		//repo:   adapter.SignInRepoFactory(SignInTrackerTable),
		signInDataService: collaborators.NewSignInTrackingService(),
	}
	registry.AddServiceProvider(HealthControllerConstants.Name, controller, core.PublicRoute)
	router.AddRoute(HealthControllerConstants.Path[0], HealthControllerConstants.Name)
	router.AddRoute("/admin/health/html", HealthControllerConstants.Name)
	return controller
}

type HealthController struct {
	logger            *zap.Logger
	signInDataService *collaborators.SignInTrackingService
}

func (hc HealthController) Receive(message core.Message, ctx core.Context) (core.Message, error) {
	var ar = new(domain.MonoResponse)
	packet, err := utils.HttpMsgExtractor(message, hc.logger, HealthControllerConstants.AllowedMethods, &ar)
	if err != nil && err.Error() != InvalidRequestBodyErrorMsg {
		return mwhttp.NewSimpleResponseText(DefaultHTTPErrorCode, err.Error()), nil
	}

	_, isHtml, err := hc.extractKeyAndResponseType(packet)
	if err != nil {
		return mwhttp.NewSimpleResponseText(InvalidKeyHTTPErrorCode, err.Error()), nil
	}

	health := &healthAlive.Health{
		Name:       "signindatatrackerws health status",
		Components: []*healthAlive.HealthComponent{hc.ProvideDbHealthComponent()},
	}
	health.Check()

	if isHtml {
		htmlContent := hc.RenderHtml(health)
		return mwhttp.NewSimpleResponseContent(DefaultHTTPSuccessCode, HTMLResponseType, htmlContent), nil
	}

	return utils.DispatchJsonResponse(health, hc.logger, http.StatusOK)
}

func (hc HealthController) extractKeyAndResponseType(packet *utils.HttpPacket) (key string, isHtml bool, err error) {
	// Replace 'PacketType' with the actual type of 'packet'.

	pathParams := strings.Split(packet.Request.Request.URL.Path, DefaultPathDelimiter)
	isHtml = pathParams[len(pathParams)-1] == "html"

	if isHtml {
		params := packet.QueryParams
		keys, exists := params[healthAlive.KeyValueVariableAlt]
		if exists && len(keys) > 0 {
			key = keys[0]
		}
	} else {
		key = packet.Request.Request.Header.Get(healthAlive.KeyValueVariable)
	}

	if key != healthAlive.DefaultKeyValue {
		err = fmt.Errorf("invalid/non-existent Monitor Key")
	}

	return
}

func (hc HealthController) ProvideDbHealthComponent() *healthAlive.HealthComponent {
	appContext := bootstrap.GetApplicationContext()
	return &healthAlive.HealthComponent{
		Name:        "Test Database Connectivity: dynamodb",
		Description: "Check if database is reachable",
		Essential:   true,
		Uri:         appContext.AppConfigData.Dynamo.EndPoint + "@" + appContext.AppConfigData.Dynamo.TableName,
		CheckHealthComponentFunc: func() (healthAlive.HealthComponentStatusCode, error) {
			/*tableNames, err := hc.signInDataService.PingDB()
			if err != nil {
				return healthAlive.Critical, err
			}
			// Here, check if your desired table (SignInTrackerTable) exists in the returned list of table names
			for _, tableName := range tableNames.TableNames {
				if tableName == SignInTrackerTable {
					return healthAlive.ComponentOk, nil
				}
			}
			return healthAlive.Critical, fmt.Errorf("Table %s not found in DynamoDB", SignInTrackerTable)
			*/
			return healthAlive.ComponentOk, nil
		},
	}
}

func (hc HealthController) RenderHtml(health *healthAlive.Health) string {
	t, err := template.ParseFS(healthAlive.Ht, "*.gohtml")
	if err != nil {
		hc.logger.Error("Error rendering html")
		return "<html><body>error</body>/<html>"
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, health)
	if err != nil {
		hc.logger.Error("Error rendering html")
		return "<html><body>error</body>/<html>"
	}
	return buf.String()
}
