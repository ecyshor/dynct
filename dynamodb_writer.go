package main

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/aws/session"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDb struct {
	Configuration *Configuration
	ApiClient     dynamodbiface.DynamoDBAPI
}

func NewDynamo(config *Configuration) *DynamoDb {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		panic(err)
	}

	client := &DynamoDb{
		Configuration:config,
		ApiClient:dynamodb.New(sess),
	}
	return client
}

func (d *DynamoDb) pipeThrough(inputChan chan *WriteEntry, output chan *WriteEntry) {
	go func() {
		for {
			entry := <-inputChan
			go func() {
				var entryJson interface{}
				err := json.Unmarshal([]byte(entry.Json), &entryJson)
				if (err != nil) {
					fmt.Println("failed to parse json,", err)
					return
				}
				parsedJson := entryJson.(map[string]interface{})
				att, err := dynamodbattribute.MarshalMap(parsedJson)
				if (err != nil) {
					fmt.Println("failed to marshal json,", err)
					return
				}
				returnValue := dynamodb.ReturnValueNone
				_, err = d.ApiClient.PutItem(&dynamodb.PutItemInput{
					TableName:&d.Configuration.TableName,
					Item:att,
					ReturnValues:&returnValue,
				})
				if (err != nil) {
					fmt.Println("failed to put item to dynamodb,", err)
					return
				} else {
					output <- entry
				}

			}()
		}
	}()
}