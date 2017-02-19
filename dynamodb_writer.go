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
	log.Infof("Starting dynamodb writer with configuration", config)
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

func (dynamo *DynamoDb) pipeThrough(inputChan chan *WriteEntry, output chan *WriteEntry) {
	go func() {
		for {
			entry := <-inputChan
			go func() {
				log.Debugf("Writing json entry to dynamo.", entry)
				var entryJson interface{}
				err := json.Unmarshal([]byte(entry.Json), &entryJson)
				if (err != nil) {
					log.Errorf("failed to parse json,", err)
					return
				}
				parsedJson := entryJson.(map[string]interface{})
				att, err := dynamodbattribute.MarshalMap(parsedJson)
				if (err != nil) {
					log.Errorf("failed to marshal json,", err)
					return
				}
				returnValue := dynamodb.ReturnValueNone
				_, err = dynamo.ApiClient.PutItem(&dynamodb.PutItemInput{
					TableName:&dynamo.Configuration.TableName,
					Item:att,
					ReturnValues:&returnValue,
				})
				if (err != nil) {
					log.Errorf("failed to put item to dynamodb,", err)
				} else {
					output <- entry
				}

			}()
		}
	}()
}