package filters

import (
	"github.mathworks.com/development/accesskeyfilter-go/pkg/accesskeyfilter"
	"github.mathworks.com/development/mito/pkg/config"
	"github.mathworks.com/development/mito/pkg/core"
	"github.mathworks.com/development/signindatatrackerws/pkg/bootstrap"
	"github.mathworks.com/development/signindatatrackerws/pkg/domain"
	"github.mathworks.com/development/signindatatrackerws/pkg/utils"
	"go.uber.org/zap"
)

type AKFilter struct {
	logger *zap.Logger
	Filter *accesskeyfilter.AccessKeyFilter
}

func (ak AKFilter) Receive(message core.Message, ctx core.Context) (core.Message, error) {
	return ak.validate(message, ctx)
}

func (ak *AKFilter) validate(message core.Message, ctx core.Context) (core.Message, error) {
	httpReq, packet, err := utils.CheckForHttpType(message, ak.logger)
	if err != nil {
		return packet, err
	}
	isValidToken, akv := ak.Filter.VerifyToken(httpReq.Request)
	if isValidToken {
		return ctx.Send(message)
	} else {
		requestid := httpReq.Request.Header.Get("mathworks-requestid")
		errresp := domain.ErrorResponse{
			ErrorCode:    4405,
			ErrorMessage: akv.Message,
			RequestID:    requestid,
		}
		return utils.DispatchJsonResponse(errresp, ak.logger, akv.Code)
	}
}

func initialize(filter *AKFilter, publicKey string, logger *zap.Logger) {
	keyFunc := func() string {
		return publicKey
	}
	// 1. create akfilterconfig object
	accesskeyfilter.DefaultExclusions = append(accesskeyfilter.DefaultExclusions, "/admin/health/v2", "/admin/health/html")
	cfg := new(accesskeyfilter.AKFilterConfig)
	// 2. initialize it with a bunch of params
	cfg.InitConfig(keyFunc, false, make([]string, 0), make(map[string]string, 0), logger)
	filter.Filter = &accesskeyfilter.AccessKeyFilter{
		FilterConfig: cfg,
	}

}

func NewAKFilter(config config.Config, registry core.Registry) *AKFilter {
	logger := zap.L().Named("signindatatrackerws.akfilter")
	r := &AKFilter{logger: logger}
	initialize(r, bootstrap.GetApplicationContext().AppConfigData.AccessKey.AccessKeyPublic, logger)
	registry.AddServiceProvider("auth/accesskeys", r)
	registry.AddRouteFilter("auth/accesskeys", "http/default", "auth/accesskeys", 10)
	return r
}
