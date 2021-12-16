package database

import (
	"api/internal/types"
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// SaveUser saves a user in the DB
func SaveUser(user types.User) error {
	// TODO: Handle error
	dynamo, _ := Client()

	data, err := attributevalue.MarshalMap(user)
	if err != nil {
		log.Printf("Error while marshalling: %v ", err)
		return err
	}

	_, err = dynamo.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      data,
		TableName: aws.String("Users"),
	})

	if err != nil {
		log.Printf("Error while saving data: %v ", err)
		return err
	}

	return nil
}

// GetUser find  a user in the database using userID
func GetUser(publicKey string) (types.User, error) {
	// TODO: Handle error
	dynamo, _ := Client()
	var user types.User

	key, err := attributevalue.MarshalMap(struct {
		PublicKey string
	}{
		publicKey,
	})

	if err != nil {
		log.Printf("Error while marshaling composite key: %v", err)
		return types.User{}, err
	}

	result, err := dynamo.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Users"),
		Key:       key,
	})

	if err != nil {
		log.Printf("Error while finding playlist: %v", err)
		return types.User{}, err
	}

	// 404 if we don't find the user
	if result.Item == nil {
		return types.User{}, errors.New("404")
	}

	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		log.Printf("Error unmarshalling: %v ", err)
		return types.User{}, err
	}

	return user, nil
}
