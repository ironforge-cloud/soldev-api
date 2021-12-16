package database

import (
	"api/internal/types"
	"context"
	"log"

	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// SaveTweets ...
func SaveTweets(tweets []types.TwitterListResponseData, replace bool) error {
	// TODO: Handle error
	dynamo, _ := Client()

	// TODO: Batch operation instead of range loop
	for _, tweet := range tweets {

		// Convert Go types to DynamoDB Attribute Values
		data, err := attributevalue.MarshalMap(tweet)
		if err != nil {
			log.Printf("Error while marshalling: %v ", err)
			return err
		}

		if replace {
			_, err = dynamo.PutItem(context.TODO(), &dynamodb.PutItemInput{
				Item:      data,
				TableName: aws.String("Social"),
			})
		} else {
			_, err = dynamo.PutItem(context.TODO(), &dynamodb.PutItemInput{
				Item:                data,
				TableName:           aws.String("Social"),
				ConditionExpression: aws.String("attribute_not_exists(ID)"),
			})
		}

		if err != nil {
			log.Printf("Error while saving data: %v ", err)
			return err
		}
	}

	return nil
}

// QueryTweets ...
func QueryTweets(listID string) ([]types.TwitterListResponseData, error) {
	// TODO: Handle error
	dynamo, _ := Client()
	var data []types.TwitterListResponseData

	result, err := dynamo.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String("Social"),
		KeyConditionExpression: aws.String("PK = :PK"),
		ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
			":PK": &dynamoTypes.AttributeValueMemberS{Value: listID},
		},
		ScanIndexForward: aws.Bool(false),
	})

	if err != nil {
		log.Printf("Error when running dynamo Query: %v", err)
		return nil, err
	}

	for _, i := range result.Items {
		var tweet types.TwitterListResponseData
		err = attributevalue.UnmarshalMap(i, &tweet)
		if err != nil {
			log.Printf("Error unmarshalling: %v ", err)
			return nil, err
		}
		data = append(data, tweet)
	}

	return data, nil
}

// QueryTweetsUsingTweetID using Tweet ID
func QueryTweetsUsingTweetID(tweetId string) ([]types.TwitterListResponseData, error) {
	// TODO: Handle error
	dynamo, _ := Client()
	var data []types.TwitterListResponseData

	result, err := dynamo.Query(context.TODO(), &dynamodb.QueryInput{
		IndexName:              aws.String("tweet-gsi"),
		TableName:              aws.String("Social"),
		KeyConditionExpression: aws.String("ID = :ID"),
		ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
			":ID": &dynamoTypes.AttributeValueMemberS{Value: tweetId},
		},
	})

	if err != nil {
		log.Printf("Error when running dynamo Query: %v", err)
		return nil, err
	}

	for _, i := range result.Items {
		var tweet types.TwitterListResponseData
		err = attributevalue.UnmarshalMap(i, &tweet)
		if err != nil {
			log.Printf("Error unmarshalling: %v ", err)
			return nil, err
		}
		data = append(data, tweet)
	}

	return data, nil
}

// QueryPinnedTweets
func QueryPinnedTweets() ([]types.TwitterListResponseData, error) {
	// TODO: Handle error
	dynamo, _ := Client()
	var data []types.TwitterListResponseData

	result, err := dynamo.Query(context.TODO(), &dynamodb.QueryInput{
		IndexName:              aws.String("pinned-tweet-gsi"),
		TableName:              aws.String("Social"),
		KeyConditionExpression: aws.String("Pinned = :Pinned"),
		ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
			":Pinned": &dynamoTypes.AttributeValueMemberN{Value: "1"},
		},
		ScanIndexForward: aws.Bool(false),
	})

	if err != nil {
		log.Printf("Error when running dynamo Query: %v", err)
		return nil, err
	}

	for _, i := range result.Items {
		var tweet types.TwitterListResponseData
		err = attributevalue.UnmarshalMap(i, &tweet)
		if err != nil {
			log.Printf("Error unmarshalling: %v ", err)
			return nil, err
		}
		data = append(data, tweet)
	}

	return data, nil
}
