package bootstrap

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.mathworks.com/development/accesskeyfilter-go/pkg/accesskeyclient"
	"go.uber.org/zap"
)

const (
	LocalEnvironment = "local"
)

var appContext *ApplicationContext

func GetApplicationContext() *ApplicationContext {
	return appContext
}

func SetApplicationContext(ctx *ApplicationContext, log *zap.Logger) {
	if ctx == nil {
		log.Error("The context can't be empty")
		return
	}
	appContext = ctx
}

func BuildApplicationContext(log *zap.Logger, overridesLoc string) *ApplicationContext {

	configData := new(AppConfigData)
	configData.OverridesLocation = overridesLoc
	configData.BootstrapConfigData(log)
	cxt := new(ApplicationContext)
	cxt.AppConfigData = configData
	cxt.GetDB = func() (DynamoDBClientInterface, error) {
		client, err := cxt.getDb(log, *configData)
		return client, err
	}
	appContext = cxt
	return cxt
}

type ApplicationContextInterface interface {
	getDb(log *zap.Logger, data AppConfigData) (DynamoDBClientInterface, error)
}

type ApplicationContext struct {
	AppConfigData *AppConfigData
	GetDB         func() (DynamoDBClientInterface, error)
	dbClient      DynamoDBClientInterface
	akClient      *accesskeyclient.AccessKeyClient
	GetAkClient   func() (*accesskeyclient.AccessKeyClient, error)
}

func (cxt *ApplicationContext) getDb(log *zap.Logger, appConfig AppConfigData) (DynamoDBClientInterface, error) {
	logger := log.With(zap.String("DynamoDB", "signindatatracker.dynamo"))
	var cfg aws.Config
	var err error
	if appConfig.Dynamo.Env == LocalEnvironment {
		cfg, err = cxt.getLocalConfig(appConfig)
	} else {
		cfg, err = cxt.getAWSConfig(appConfig)
	}

	if err != nil {
		logger.Error("Failed to load AWS DynamoDB configuration", zap.Error(err))
		return nil, err
	}

	return NewDynamoDBClient(cfg), nil
}

func (cxt *ApplicationContext) getLocalConfig(appConfig AppConfigData) (aws.Config, error) {
	return config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(appConfig.Dynamo.Region),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: appConfig.Dynamo.EndPoint}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     "abcd",
				SecretAccessKey: "a1b2c3",
				SessionToken:    "",
				Source:          "Mock credentials for local instance",
			},
		}),
	)
}

func (cxt *ApplicationContext) getAWSConfig(appConfig AppConfigData) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return aws.Config{}, err
	}
	return cfg, nil
}
