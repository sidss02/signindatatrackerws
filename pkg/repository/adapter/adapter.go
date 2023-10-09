package adapter

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.mathworks.com/development/signindatatrackerws/pkg/bootstrap"
	"github.mathworks.com/development/signindatatrackerws/pkg/domain"
	"github.mathworks.com/development/signindatatrackerws/pkg/utils"
	"go.uber.org/zap"
)

type DynamoDBData map[string]types.AttributeValue

type SignInRepoInterface interface {
	SaveSignInTrackingInfo(request domain.SaveSignInInfo) (response domain.SaveSignInInfo, err error)
	FindUniqueSignInInfo(condition map[string]interface{}) (response *dynamodb.GetItemOutput, err error)
	FindSignInTrackingDetails(partitionKey string) (response []domain.SignInInfo, err error)
	GetSignInBetweenTimeStamps(request domain.RequestTimestampInput) ([]domain.SignInInfo, error)
	GetSignInForReferenceId(request domain.RequestReferenceIdInput) ([]domain.SignInInfo, error)
	PingDB() (*dynamodb.ListTablesOutput, error)
}

type SignInRepo struct {
	logger    *zap.Logger
	dbClient  bootstrap.DynamoDBClientInterface
	tableName string
}

func SignInRepoFactory(tableName string) *SignInRepo {
	dbClient, err := bootstrap.GetApplicationContext().GetDB()
	if err != nil {
		// Handle the error, maybe log it and exit
		log.Fatalf("Failed to initialize DynamoDB client: %v", err)
	}
	return &SignInRepo{logger: zap.L().Named("signindatatrackerws.signinRepo"), dbClient: dbClient, tableName: tableName}
}

func (repo *SignInRepo) SaveSignInTrackingInfo(request domain.SaveSignInInfo) (response domain.SaveSignInInfo, err error) {
	entityParsed, err := attributevalue.MarshalMap(request)
	if err != nil {
		return domain.SaveSignInInfo{}, err
	}

	err = putItem(repo.dbClient, repo.tableName, entityParsed)
	return request, err
}

// putItem inserts an item (key + attributes) in to a dynamodb table.
func putItem(c bootstrap.DynamoDBClientInterface, tableName string, item DynamoDBData) (err error) {
	_, err = c.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName), Item: item,
	})
	return err
}

func (repo *SignInRepo) FindUniqueSignInInfo(condition map[string]interface{}) (response *dynamodb.GetItemOutput, err error) {
	conditionParsed, err := attributevalue.MarshalMap(condition) //condition - what we want to match
	if err != nil {
		return nil, err
	}
	input := &dynamodb.GetItemInput{
		TableName: aws.String(repo.tableName),
		Key:       conditionParsed,
	}
	response, err = repo.dbClient.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (repo *SignInRepo) FindSignInTrackingDetails(partitionKeyValue string) (response []domain.SignInInfo, err error) {

	input := &dynamodb.QueryInput{
		TableName: aws.String(repo.tableName),
		KeyConditions: map[string]types.Condition{
			"uniqueId": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: partitionKeyValue},
				},
			},
		},
	}

	// Make the DynamoDB Query API call
	resp, err := repo.dbClient.Query(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	// Call unmarshalItems to unmarshal the result items
	items, err := utils.UnmarshalItems(resp.Items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *SignInRepo) GetSignInBetweenTimeStamps(request domain.RequestTimestampInput) ([]domain.SignInInfo, error) {
	// Define the filter expression and attribute values for the scan operation
	expr := "#uid = :uid_value AND #ts BETWEEN :start_time AND :end_time"
	// Set up the scan input
	input := &dynamodb.ScanInput{
		TableName:        aws.String(repo.tableName),
		FilterExpression: &expr,
		ExpressionAttributeNames: map[string]string{
			"#uid": "uniqueId",
			"#ts":  "timestamp",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid_value":  &types.AttributeValueMemberS{Value: request.UniqueID},
			":start_time": &types.AttributeValueMemberS{Value: changeUTCTimeEpoch(request.StartTime)},
			":end_time":   &types.AttributeValueMemberS{Value: changeUTCTimeEpoch(request.EndTime)},
		},
	}

	// Execute the scan operation
	resp, err := repo.dbClient.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	// Call unmarshalItems to unmarshal the result items
	items, err := utils.UnmarshalItems(resp.Items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *SignInRepo) GetSignInForReferenceId(request domain.RequestReferenceIdInput) ([]domain.SignInInfo, error) {
	// Define the filter expression and attribute values for the scan operation
	expr := "#uid = :uid_value AND #refId = :refId_value"

	// Set up the scan input
	input := &dynamodb.ScanInput{
		TableName:        aws.String(repo.tableName),
		FilterExpression: &expr,
		ExpressionAttributeNames: map[string]string{
			"#uid":   "uniqueId",
			"#refId": "referenceId",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid_value":   &types.AttributeValueMemberS{Value: request.UniqueID},
			":refId_value": &types.AttributeValueMemberS{Value: request.ReferenceId},
		},
	}

	// Execute the scan operation
	resp, err := repo.dbClient.Scan(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	// Call unmarshalItems to unmarshal the result items
	items, err := utils.UnmarshalItems(resp.Items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
func (repo *SignInRepo) PingDB() (*dynamodb.ListTablesOutput, error) {

	// Using ListTables as a way to check the connectivity
	// Adjust as needed based on your DynamoDB setup and permissions
	input := &dynamodb.ListTablesInput{
		Limit: aws.Int32(1), // Limiting to one table just to reduce the response size
	}
	result, err := repo.dbClient.ListTables(context.TODO(), input)
	if err != nil {
		return nil, errors.New("failed to get connection to db: " + err.Error())
	}

	return result, nil
}

func changeUTCTimeEpoch(timeStr string) string {
	// Parse as an epoch time (Unix timestamp)
	epochTime, err := strconv.ParseInt(timeStr, 10, 64)
	if err == nil {
		// If parsing as epoch time succeeds, return it as is
		return strconv.FormatInt(epochTime, 10)
	}
	// If parsing as epoch time fails, try to parse as a datetime
	parsedTime, err := time.Parse(time.DateTime, timeStr)
	if err != nil {
		fmt.Println("Error: Failed to Parse the timestamp", err)
		return ""
	}

	// Convert the time.Time object to Unix timestamp (epoch time)
	return strconv.FormatInt(parsedTime.Unix(), 10)
}
