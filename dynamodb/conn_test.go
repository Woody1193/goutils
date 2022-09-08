package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Database Connection Tests", func() {

	// Ensure that the AWS config is created before each test; this could be set as a global variable
	var cfg aws.Config
	BeforeAll(func() {
		cfg = TestDynamoDBConfig(context.Background(), "us-east-1", 9000)
	})

	// Create our test table definition that we'll use for all module tests
	testTable := dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("sort_key"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("sort_key"),
				KeyType:       types.KeyTypeRange,
			},
		},
		TableName:   aws.String("TEST_TABLE"),
		BillingMode: types.BillingModeProvisioned,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableClass: types.TableClassStandard,
	}

	// Esnure that the table exists before the start of each test
	BeforeEach(func() {
		if err := EnsureTableExists(context.Background(), cfg, &testTable); err != nil {
			panic(err)
		}
	})

	// Ensure that the table is empty at the end of each test (not strictly necessary if test data is isolated)
	AfterEach(func() {
		if err := EmptyTable(context.Background(), cfg, &testTable); err != nil {
			panic(err)
		}
	})

	It("PutItem - Fails - Error", func() {

	})
})
