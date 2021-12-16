package database

import (
	"api/internal/types"
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// GetAllPlaylists find all the playlists in saved in the database
func GetAllPlaylists(vertical string) ([]types.Playlist, error) {
	// TODO: Handle error
	dynamo, _ := Client()
	var data []types.Playlist

	result, err := dynamo.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String("Playlists"),
		IndexName:              aws.String("vertical-gsi"),
		KeyConditionExpression: aws.String("Vertical = :Vertical"),
		ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
			":Vertical": &dynamoTypes.AttributeValueMemberS{Value: vertical},
		},
	})

	if err != nil {
		log.Printf("Error when running dynamo Query: %v", err)
		return nil, err
	}

	for _, i := range result.Items {
		var playlist types.Playlist
		err = attributevalue.UnmarshalMap(i, &playlist)
		if err != nil {
			log.Printf("Error unmarshalling: %v ", err)
			return nil, err
		}
		data = append(data, playlist)
	}

	return data, nil
}

// SavePlaylists saves one or multiple playlists in the database
func SavePlaylists(playlists []types.Playlist) error {
	// TODO: Handle error
	dynamo, _ := Client()

	// TODO: Batch operation instead of range loop
	for i, playlist := range playlists {

		// Assigning sort order
		playlist.Position = i

		// Convert Go types to DynamoDB Attribute Values
		data, err := attributevalue.MarshalMap(playlist)
		if err != nil {
			log.Printf("Error while marshalling: %v ", err)
			return err
		}

		_, err = dynamo.PutItem(context.TODO(), &dynamodb.PutItemInput{
			Item:      data,
			TableName: aws.String("Playlists"),
		})
		if err != nil {
			log.Printf("Error while saving data: %v ", err)
			return err
		}
	}

	return nil
}

// GetPlaylistsByProvider ...
func GetPlaylistsByProvider(provider string) ([]types.Playlist, error) {
	// TODO: Handle error
	dynamo, _ := Client()
	var data []types.Playlist

	result, err := dynamo.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String("Playlists"),
		IndexName:              aws.String("provider-gsi"),
		KeyConditionExpression: aws.String("Provider = :provider"),
		ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
			":provider": &dynamoTypes.AttributeValueMemberS{Value: provider},
		},
	})

	if err != nil {
		log.Printf("Error when running dynamo Query: %v", err)
		return nil, err
	}

	for _, i := range result.Items {
		var playlist types.Playlist
		err = attributevalue.UnmarshalMap(i, &playlist)
		if err != nil {
			log.Printf("Error unmarshalling: %v ", err)
			return nil, err
		}
		data = append(data, playlist)
	}

	return data, nil
}

// GetPlaylistByID queries the database to find playlist using
// the playlist id provided
func GetPlaylistByID(vertical string, ID string) (types.Playlist, error) {
	// TODO: Handle error
	dynamo, _ := Client()
	var playlist types.Playlist

	result, err := dynamo.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String("Playlists"),
		KeyConditionExpression: aws.String("Vertical = :Vertical and ID = :ID"),
		ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
			":Vertical": &dynamoTypes.AttributeValueMemberS{Value: vertical},
			":ID":       &dynamoTypes.AttributeValueMemberS{Value: ID},
		},
	})

	if err != nil {
		log.Printf("Error while finding playlist: %v", err)
		return types.Playlist{}, err
	}

	err = attributevalue.UnmarshalMap(result.Items[0], &playlist)
	if err != nil {
		log.Printf("Error unmarshalling: %v ", err)
		return types.Playlist{}, err
	}

	return playlist, nil
}
