package main

import (
	"github.mathworks.com/development/mito/pkg/config" // used for setting the configuration override (configfiles.AppDefault)
	"github.mathworks.com/development/mitoapp/pkg/debug"
	"github.mathworks.com/development/mitoapp/pkg/host"       // dependency injection and application hosting (lifecycle)
	"github.mathworks.com/development/mitoapp/pkg/webservice" // used for configuring all the mito dependencies
	"github.mathworks.com/development/signindatatrackerws/pkg/bootstrap"
	"github.mathworks.com/development/signindatatrackerws/pkg/controllers"
	"github.mathworks.com/development/signindatatrackerws/pkg/filters"
	"go.uber.org/zap"
)

type App struct {
	Services                    *webservice.Services
	AliveController             *controllers.AliveController
	PersistSignInDataController *controllers.PersistSignInDataController
	GetSignInDataController     *controllers.RetrieveSignInDataController
	HealthController            *controllers.HealthController
	Filters                     *filters.AKFilter
	DebugMessageClient          *debug.MessageClient
}

// wire all controllers instantiations here
var factories = []interface{}{
	controllers.AliveControllerFactory,
	controllers.PersistSignInControllerFactory,
	controllers.RetrieveSignInControllerFactory,
	controllers.HealthControllerFactory,
	filters.NewAKFilter,
}

var app *App

func main() {

	configuration := webservice.Defaults(config.AppDefault{
		"mito.http.contextroot":                    "/signindatatrackerws",
		"mito.config.env.filter":                   "^signinartifacts.*",
		"mito.http.headerstocontext":               `user-agent|userAgent,X-Forwarded-For|xForwardedFor,Accept-Language|acceptLanguage,X-MW-Caller-Id|xMWCallerId`,
		"mito.http.truststore.validatecertificate": "false",
		"signindatatracker.dynamo.endpoint":        "http://localhost:8000",
		"signindatatracker.dynamo.region":          "localhost",
		"signindatatracker.dynamo.tablename":       "signintracker",
		"signindatatracker.dynamo.ttlinyears":      "4",
		"signindatatracker.dynamo.env":             "local",
		"mito.debug":                               "false",
	})

	// instantiate an empty app struct to be filled in
	app = new(App)
	//Bootstrap App
	bootstrap.BuildApplicationContext(zap.L().Named("signindatatrackerws.app.context"), "")

	// recursively construct the app using the provided constructor functions
	host.Start("signindatatrackerws", app,
		configuration,
		webservice.Constructors,
		factories,

		debug.Constructors,
	)
}
