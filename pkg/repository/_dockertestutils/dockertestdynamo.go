package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"time"
)

type DynamoDBData map[string]types.AttributeValue

func newclient(profile string) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("localhost"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: "http://dynamodb-local:8000"}, nil //For spin up using docker-compose.yml
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "abcd", SecretAccessKey: "a1b2c3", SessionToken: "",
				Source: "Mock credentials used above for local instance",
			},
		}),
	)
	if err != nil {
		return nil, err
	}

	c := dynamodb.NewFromConfig(cfg)
	return c, nil
}

// createTable creates a table in the client's dynamodb instance
// using the table parameters specified in input.
func createTable(c *dynamodb.Client,
	tableName string, input *dynamodb.CreateTableInput,
) error {
	var _ *types.TableDescription
	table, err := c.CreateTable(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to create table `%v` with error: %v\n", tableName, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(c)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName)}, 5*time.Minute)
		if err != nil {
			log.Printf("Failed to wait on create table `%v` with error: %v\n", tableName, err)
		}
		_ = table.TableDescription
	}
	fmt.Printf("Created table `%s", tableName)

	return err
}

func clearDB(c *dynamodb.Client) error {
	tables, err := listTables(c)
	if err != nil {
		return err
	}

	for _, t := range tables {
		_, err := c.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{TableName: aws.String(t)})
		if err != nil {
			return err
		}
	}

	return nil
}

// listTables returns a list of table names in the client's dynamodb instance.
func listTables(c *dynamodb.Client) ([]string, error) {
	tables, err := c.ListTables(
		context.TODO(),
		&dynamodb.ListTablesInput{},
	)
	if err != nil {
		return nil, err
	}

	return tables.TableNames, nil
}

func getSignInList() (signindata []DynamoDBData) {
	list := []struct {
		uniqueId    string
		timestamp   string
		callerId    string
		ipAddress   string
		userAgent   string
		sourceId    string
		region      string
		referenceId string
	}{
		{uniqueId: "MWA-1234445", timestamp: "1661285996251", callerId: "MWA", ipAddress: "192.168.0.0", userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36", sourceId: "", region: "us-east1", referenceId: ""},
		{uniqueId: "MWA-1234446", timestamp: "1661285996252", callerId: "MWA", ipAddress: "", userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36", sourceId: "web", region: "eu-west1", referenceId: ""},
		{uniqueId: "MWA-1234445", timestamp: "1661285996253", callerId: "MWA", ipAddress: "", userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36", sourceId: "desktop", region: "us-east1", referenceId: ""},
		{uniqueId: "MWA-1234445", timestamp: "1661285996254", callerId: "MWA", ipAddress: "", userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36", sourceId: "", region: "us-east1", referenceId: ""},
		{uniqueId: "MWA-1234446", timestamp: "1661285996255", callerId: "MWA", ipAddress: "192.168.0.1", userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36", sourceId: "pat", region: "eu-west1", referenceId: ""},
		{uniqueId: "MWA-1234450", timestamp: "1661285996256", callerId: "MWA", ipAddress: "", userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36", sourceId: "", region: "", referenceId: ""},
		{uniqueId: "PAT-1234451", timestamp: "1661285996257", callerId: "PAT", ipAddress: "", userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36", sourceId: "web", region: "us-east1", referenceId: "123456"},
		{uniqueId: "PAT-1234451", timestamp: "1661285996258", callerId: "PAT", ipAddress: "", userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36", sourceId: "", region: "", referenceId: "123456"},
		{uniqueId: "PAT-1234453", timestamp: "1661285996259", callerId: "PAT", ipAddress: "", userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36", sourceId: "web", region: "us-east1", referenceId: "123458"},
		{uniqueId: "PAT-1234454", timestamp: "1661285996261", callerId: "PAT", ipAddress: "", userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36", sourceId: "", region: "", referenceId: "123459"},
		{uniqueId: "PAT-1234453", timestamp: "1661285996262", callerId: "PAT", ipAddress: "", userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36", sourceId: "desktop", region: "us-east1", referenceId: ""},
		{uniqueId: "MWA-1234456", timestamp: "1661285996263", callerId: "MWA", ipAddress: "192.168.0.3", userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36", sourceId: "web", region: "us-east1", referenceId: ""},
		{uniqueId: "MWA-1234456", timestamp: "1661285996264", callerId: "MWA", ipAddress: "192.168.0.0", userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36", sourceId: "", region: "", referenceId: ""},
		{uniqueId: "MWA-1234456", timestamp: "1661285996265", callerId: "MWA", ipAddress: "", userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36", sourceId: "", region: "us-east1", referenceId: ""},
	}

	for _, m := range list {
		if m.referenceId == "" {
			m.referenceId = "NULL"
		}

		signindata = append(signindata, DynamoDBData{
			"uniqueId":    unsafeToAttrValue(m.uniqueId),
			"timestamp":   unsafeToAttrValue(m.timestamp),
			"callerId":    unsafeToAttrValue(m.callerId),
			"ipAddress":   unsafeToAttrValue(m.ipAddress),
			"userAgent":   unsafeToAttrValue(m.userAgent),
			"sourceId":    unsafeToAttrValue(m.sourceId),
			"region":      unsafeToAttrValue(m.region),
			"referenceId": unsafeToAttrValue(m.referenceId),
		})

	}

	return signindata
}

// SignInData represents our domain entity.
type SignInData struct {
	uniqueId    string `dynamodbav:"uniqueId"`
	timestamp   string `dynamodbav:"timestamp"`
	callerId    string `dynamodbav:"callerId"`
	ipAddress   string `dynamodbav:"ipAddress"`
	userAgent   string `dynamodbav:"userAgent"`
	sourceId    string `dynamodbav:"sourceId"`
	region      string `dynamodbav:"region"`
	referenceId string `dynamodbav:"referenceId"`
}

// ---------------------------------------------------- UTILS
func unsafeToAttrValue(in interface{}) types.AttributeValue {
	val, err := attributevalue.Marshal(in)
	if err != nil {
		log.Fatalf("could not marshal value `%v` with error: %v", in, err)
	}

	return val
}

// putItem inserts an item (key + attributes) in to a dynamodb table.
func putItem(c *dynamodb.Client, tableName string, item DynamoDBData) (err error) {
	_, err = c.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName), Item: item,
	})
	if err != nil {
		return err
	}

	return nil
}

func putItems(c *dynamodb.Client, tableName string, items []DynamoDBData) (err error) {
	// dynamodb batch limit is 25
	if len(items) > 25 {
		return fmt.Errorf("Max batch size is 25, attempted `%d`", len(items))
	}

	// create requests
	writeRequests := make([]types.WriteRequest, len(items))
	for i, item := range items {
		writeRequests[i] = types.WriteRequest{PutRequest: &types.PutRequest{Item: item}}
	}

	// write batch
	_, err = c.BatchWriteItem(
		context.TODO(),
		&dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{tableName: writeRequests},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// getItem returns an item if found based on the key provided.
// the key could be either a primary or composite key and values map.
func getItem(c *dynamodb.Client, tableName string, key DynamoDBData) (item DynamoDBData, err error) {
	resp, err := c.GetItem(context.TODO(), &dynamodb.GetItemInput{Key: key, TableName: aws.String(tableName)})
	if err != nil {
		return nil, err
	}

	return resp.Item, nil //
}
func main() {

	var c *dynamodb.Client

	var err error

	// singleton client
	if c == nil {
		var err error
		c, err = newclient("local-dynodb-admin") // named profile
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("**********************\nStarting Local Dynamo db...\n\n")

	// clear database tables
	err = clearDB(c)
	if err != nil {
		log.Fatal(err)
	}

	// example table name
	exampleTableName := "signindatatracker"

	// create table
	tableInput := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("uniqueId"),
				AttributeType: types.ScalarAttributeTypeS, // data type descriptor: S == string
			},
			{
				AttributeName: aws.String("timestamp"),
				AttributeType: types.ScalarAttributeTypeS, // data type descriptor: S == string
			},
			{
				AttributeName: aws.String("referenceId"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{ // key: uniqueId + timestamp
			{
				AttributeName: aws.String("uniqueId"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("timestamp"),
				KeyType:       types.KeyTypeRange,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("UniqueIdReferenceIdIndex"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("uniqueId"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("referenceId"),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(10),
					WriteCapacityUnits: aws.Int64(10),
				},
			},
		},

		TableName: aws.String(exampleTableName),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}
	err = createTable(c, exampleTableName, tableInput)
	if err != nil {
		log.Fatal(err)
	}

	// -----------------------------
	// list tables (should return single table, since we only created one here!)
	tables, err := listTables(c)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Tables: %v\n\n", tables)

	if err != nil {
		log.Printf("%v", err.Error())
	}

	// add items (signindata), single then in batch
	signindata := getSignInList()
	err = putItem(c, exampleTableName, signindata[0])
	if err != nil {
		log.Fatal(err)
	}
	err = putItems(c, exampleTableName, signindata[1:])
	if err != nil {
		log.Fatal(err)
	}

	// -----------------------------
	// get item
	uniqueId := "MWA-1234445"
	timestamp := "1661285996251"
	uniqueIdAttr, _ := attributevalue.Marshal(uniqueId)
	timestampAttr, _ := attributevalue.Marshal(timestamp)

	item, err := getItem(c, exampleTableName, DynamoDBData{"uniqueId": uniqueIdAttr, "timestamp": timestampAttr})
	if err != nil {
		log.Fatal(err)
	}

	var signin SignInData
	// unmarshal item
	err = attributevalue.UnmarshalMap(item, &signin)
	if err != nil {
		log.Fatal(err)
	}
}
