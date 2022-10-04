package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/21satvik/dynamodb-go/config"
	"github.com/21satvik/dynamodb-go/internal/repository/adapter"
	"github.com/21satvik/dynamodb-go/internal/repository/instance"
	"github.com/21satvik/dynamodb-go/internal/routes"
	"github.com/21satvik/dynamodb-go/internal/rules"
	RulesProduct "github.com/21satvik/dynamodb-go/internal/rules/product"
	"github.com/21satvik/dynamodb-go/utils/logger"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	configs := config.GetConfig()
	connection := instance.GetConnection()
	repository := adapter.NewAdpter(connection)

	logger.INFO("Waiting for the server to start...", nil)

	errors := Migrate(connection)
	if len(errors) > 0 {
		for _, err := range errors {
			logger.PANIC("Error in migration: ", err)
		}
	}
	logger.PANIC("Error in migration: ", checkTables(connection))

	port := fmt.Sprintf(":%v", configs.Port)
	router := routes.NewRouter().SetRouters(repository)
	logger.INFO("Server is running on port: ", port)

	server := http.ListenAndServe(port, router)
	log.Fatal(server)
}

func Migrate(connection *dynamodb.DynamoDB) []error {
	var errors []error

	callMigrateAndAppendError(&errors, connection, &RulesProduct.Rules{})

	return errors
}

func callMigrateAndAppendError(errors *[]error, connection *dynamodb.DynamoDB, rules rules.Interface) {
	err := rules.Migrate(connection)
	if err != nil {
		*errors = append(*errors, err)
	}
}

func checkTables(connection *dynamodb.DynamoDB) error {
	response, err := connection.ListTables(&dynamodb.ListTablesInput{})

	if response != nil {
		if len(response.TableNames) == 0 {
			logger.INFO("Tables not found: ", nil)
		}

		for _, tableName := range response.TableNames {
			logger.INFO("Table Found: ", tableName)
		}
	}
	return err
}
