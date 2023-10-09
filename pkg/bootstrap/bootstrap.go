package bootstrap

import (
	"github.mathworks.com/development/opi-utils-go/pkg/configutils"
	"github.mathworks.com/development/signindatatrackerws/pkg/utils"
	"go.uber.org/zap"
)

type AppConfigDataInterface interface {
	BootstrapConfigData()
	loadOverrides() error
}

var DefaultOverridesLocation = "properties/overrides.properties"

type AppConfigData struct {
	Db                DatabaseConfig
	AccessKey         AccessKeyConfig
	Dynamo            DynamoConfig
	AppCallerId       string
	AppRunTime        string
	OverridesLocation string
	logger            *zap.Logger
}
type DatabaseConfig struct {
	DbHost string
	DbPort string
	Name   string
}
type AccessKeyConfig struct {
	AccessKeyPublic string
	AccessKeyHost   string
}
type DynamoConfig struct {
	EndPoint  string
	TableName string
	Region    string
	Env       string
}

func (appConfig *AppConfigData) BootstrapConfigData(logger *zap.Logger) {
	err := appConfig.loadOverrides()
	if err != nil {
		logger.Error("Error loading overrides.properties: " + err.Error())
	}
}
func (appConfig *AppConfigData) loadOverrides() error {
	loc := DefaultOverridesLocation
	if appConfig.OverridesLocation != "" {
		loc = appConfig.OverridesLocation
	}
	props, err := configutils.ReadPropFile(loc)
	if err != nil {
		return err
	}
	runTime := utils.GetValueFromMap(props, "app.runTime", "local")
	appConfig.AppRunTime = runTime
	appConfig.AppCallerId = "SIGNINDATATRACKINGWS"
	appConfig.AccessKey.AccessKeyPublic = utils.GetValueFromMap(props, "app.signindatatracker.ak.public", "")
	appConfig.AccessKey.AccessKeyHost = utils.GetValueFromMap(props, "app.signindatatracker.ak.host", "")
	appConfig.Dynamo.EndPoint = utils.GetValueFromMap(props, "app.signindatatracker.dynamo.endpoint", "")
	appConfig.Dynamo.TableName = utils.GetValueFromMap(props, "app.signindatatracker.dynamo.tablename", "")
	appConfig.Dynamo.Region = utils.GetValueFromMap(props, "app.signindatatracker.dynamo.region", "")
	appConfig.Dynamo.Env = utils.GetValueFromMap(props, "app.signindatatracker.dynamo.env", "")
	return nil
}
