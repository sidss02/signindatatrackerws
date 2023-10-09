package utils

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.mathworks.com/development/signindatatrackerws/pkg/domain"
)

func GetValueFromMap(mp map[string]interface{}, key string, def string) string {
	v, e := mp[key]
	if e {
		return v.(string)
	}
	return def
}

func UnmarshalItems(items []map[string]types.AttributeValue) ([]domain.SignInInfo, error) {

	//Unmarshal the result items into a slice of domain.SignInInfo structs
	var result []domain.SignInInfo
	for _, item := range items {
		var singleItem domain.SignInInfo
		err := attributevalue.UnmarshalMap(item, &singleItem)
		if err != nil {
			return nil, err
		}
		result = append(result, singleItem)
	}
	return result, nil
}
