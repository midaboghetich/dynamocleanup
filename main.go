package main

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
	"strings"

	"fmt"
)

func main() {
	numArgs := len(os.Args)
	if numArgs < 2 {
		fmt.Println(">> No args passed in")
		os.Exit(0)
	}
	user := os.Args[1]

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// create the input configuration instance
	input := &dynamodb.ListTablesInput{}

	fmt.Printf("Tables:\n")

	for {
		// Get the list of tables
		result, err := svc.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {

				fmt.Println(err.Error())
			}
			return
		}

		for _, n := range result.TableNames {
			fmt.Printf(*n + "\n");
			if strings.Contains(*n, user) {
				fmt.Println("Deleting -> " + *n)
				deleteTable(svc, n)
			}
		}

		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}

}

func deleteTable(client *dynamodb.DynamoDB, tableName *string) {
	input := &dynamodb.DeleteTableInput{TableName: tableName}

	table, err := client.DeleteTable(input)

	if err != nil {

		fmt.Println(err.Error())
		return
	}

	fmt.Println(table.String())

}
