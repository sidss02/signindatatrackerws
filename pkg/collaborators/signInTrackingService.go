package collaborators

import (
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.mathworks.com/development/signindatatrackerws/pkg/domain"
	"github.mathworks.com/development/signindatatrackerws/pkg/repository/adapter"
	"go.uber.org/zap"
)

type SignInTrackingServiceInterface interface {
	FindSignInTrackingInfo(input domain.RequestInput) (domain.RequestInput, error)
	SaveSignInData(input domain.SaveSignInInfo) (domain.SaveSignInInfo, error)
	FindSignInTrackingDetails(request domain.RequestDetailsInput) ([]domain.SignInInfo, domain.ErrorResponse, int)
	PingDB() error
}

type SignInTrackingService struct {
	logger *zap.Logger
	repo   adapter.SignInRepoInterface
	keys   []string
}

const SignInTrackerTable = "signindatatracker"

func NewSignInTrackingService() *SignInTrackingService {
	svc := &SignInTrackingService{

		logger: zap.L().Named("signindatatrackerws.signinTracking"),
		repo:   adapter.SignInRepoFactory(SignInTrackerTable),
	}
	return svc
}

func (ps *SignInTrackingService) SaveSignInData(request domain.SaveSignInInfo) (domain.SaveSignInInfo, domain.ErrorResponse, int) {

	// Check if uniqueId is empty
	if request.UniqueId == "" {
		errresp := domain.ErrorResponse{
			ErrorCode:    5300,
			ErrorMessage: "UniqueId cannot be empty",
		}
		return domain.SaveSignInInfo{}, errresp, http.StatusBadRequest
	}

	if request.ReferenceId == "" {
		request.ReferenceId = "NULL"
	}

	signInInfo := domain.SaveSignInInfo{
		UniqueId:    request.UniqueId,
		TimeStamp:   strconv.FormatInt(int64(time.Now().UnixMilli()), 10),
		CalledId:    request.CalledId,
		SourceId:    request.SourceId,
		IpAddress:   request.IpAddress,
		Region:      request.Region,
		UserAgent:   request.UserAgent,
		ReferenceId: request.ReferenceId,
		SsoOrgId:    request.SsoOrgId,
	}

	profiles, err := ps.repo.SaveSignInTrackingInfo(signInInfo)
	if err != nil {
		errresp := domain.ErrorResponse{
			ErrorCode:    5500,
			ErrorMessage: "Could not save SignInData",
			Error:        err.Error(),
		}
		return domain.SaveSignInInfo{}, errresp, http.StatusInternalServerError
	}

	return profiles, domain.ErrorResponse{}, http.StatusOK
}
func (ps *SignInTrackingService) FindUniqueSignInInfo(request domain.RequestInput) (domain.SignInInfo, domain.ErrorResponse, int) {
	// Define your condition and tableName
	condition := map[string]interface{}{
		"uniqueId":  request.UniqueID,  // Use UniqueID from RequestInput
		"timestamp": request.Timestamp, // Use ReferenceId from RequestInput
	}
	// Call the FindUniqueSignInInfo function
	profiles, err := ps.repo.FindUniqueSignInInfo(condition)
	if err != nil {
		errresp := domain.ErrorResponse{
			ErrorCode:    5500,
			ErrorMessage: "Could not execute findSignInTrackingInfo",
			Error:        err.Error(),
		}
		return domain.SignInInfo{}, errresp, http.StatusBadRequest
	}

	// Map the DynamoDB output to RequestInput
	resultOutput, err := mapDynamoDBOutputToRequestInput(profiles)

	if err != nil {
		// Handle the error appropriately, e.g., return an error response
		errresp := domain.ErrorResponse{
			ErrorCode:    5501,
			ErrorMessage: "Could not map DynamoDB output to output",
		}
		return domain.SignInInfo{}, errresp, http.StatusInternalServerError
	}

	return resultOutput, domain.ErrorResponse{}, http.StatusOK
}

func (ps *SignInTrackingService) FindSignInTrackingDetails(request domain.RequestDetailsInput) ([]domain.SignInInfo, domain.ErrorResponse, int) {

	// Call the FindSignInTrackingDetails function
	profiles, err := ps.repo.FindSignInTrackingDetails(request.UniqueID)
	if err != nil {
		errresp := domain.ErrorResponse{
			ErrorCode:    5500,
			ErrorMessage: "Could not execute findSignInTrackingInfo",
			Error:        err.Error(),
		}
		return []domain.SignInInfo{}, errresp, http.StatusBadRequest
	}

	return profiles, domain.ErrorResponse{}, http.StatusOK
}

func (ps *SignInTrackingService) FindSignInReferenceIds(request domain.RequestReferenceIdInput) ([]domain.SignInInfo, domain.ErrorResponse, int) {

	// Call the FindSignInTrackingDetails function
	profiles, err := ps.repo.GetSignInForReferenceId(request)
	if err != nil {
		errresp := domain.ErrorResponse{
			ErrorCode:    5500,
			ErrorMessage: "Could not execute findSignInTrackingInfo",
			Error:        err.Error(),
		}
		return []domain.SignInInfo{}, errresp, http.StatusBadRequest
	}

	return profiles, domain.ErrorResponse{}, http.StatusOK
}

func (ps *SignInTrackingService) FindSignInPeriodDetails(request domain.RequestTimestampInput) ([]domain.SignInInfo, domain.ErrorResponse, int) {

	// Call the FindSignInTrackingDetails function
	profiles, err := ps.repo.GetSignInBetweenTimeStamps(request)
	if err != nil {
		errresp := domain.ErrorResponse{
			ErrorCode:    5500,
			ErrorMessage: "Could not execute findSignInTrackingInfo",
			Error:        err.Error(),
		}
		return []domain.SignInInfo{}, errresp, http.StatusBadRequest
	}

	return profiles, domain.ErrorResponse{}, http.StatusOK
}

func (ps *SignInTrackingService) PingDB() (*dynamodb.ListTablesOutput, error) {

	table, err := ps.repo.PingDB()
	if err != nil {
		_ = domain.ErrorResponse{
			ErrorCode:    5700,
			ErrorMessage: "Could not connect to Dynamo DB Table",
			Error:        err.Error(),
		}
		return nil, err
	}
	return table, nil
}

func mapDynamoDBOutputToRequestInput(output *dynamodb.GetItemOutput) (domain.SignInInfo, error) {
	var requestInput domain.SignInInfo

	// Extract and map the relevant fields from the DynamoDB output to the RequestInput struct
	err := attributevalue.UnmarshalMap(output.Item, &requestInput)
	if err != nil {
		return domain.SignInInfo{}, err
	}

	return requestInput, nil
}
