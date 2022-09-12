package dynamodb

import (
	"context"

	"github.com/Woody1193/goutils/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database Connection Tests", Ordered, func() {

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

	// Test that, if the inner PutItem request fails, then calling PutItem will return an error
	It("PutItem - Fails - Error", func() {

		// First, create our test database connection from our test config
		conn := createTestConnection(cfg)

		// Next, create our test object with some fake data
		data := testObject{
			ID:      "test_id",
			SortKey: "test|sort|key",
			Data:    "test",
		}

		// Attempt to marshal the test object into a DynamoDB item structure
		attrs, err := attributevalue.MarshalMapWithOptions(&data,
			func(eo *attributevalue.EncoderOptions) { eo.TagKey = "json" })
		Expect(err).ShouldNot(HaveOccurred())

		// Now, create our put-item input from our attribute data
		input := dynamodb.PutItemInput{
			Item:                        attrs,
			TableName:                   aws.String("FAKE_TABLE"),
			ReturnConsumedCapacity:      types.ReturnConsumedCapacityNone,
			ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsNone,
			ReturnValues:                types.ReturnValueAllNew,
		}

		// Finally, attempt to put the item to the database; this should fail
		output, err := conn.PutItem(context.Background(), &input)

		// Verify the failure
		Expect(output).Should(BeNil())
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(And(
			HavePrefix("operation error DynamoDB: PutItem, https response error StatusCode: 400, RequestID: "),
			HaveSuffix(", ResourceNotFoundException: ")))
	})

	// Test that, if no failure occurs, then calling PutItem will result in the item being written
	// to the associated table in the database
	It("PutItem - No failures - Data exists", func() {

		// First, create our test database connection from our test config
		conn := createTestConnection(cfg)

		// Create our test object with some fake data
		data := testObject{
			ID:      "test_id",
			SortKey: "test|sort|key",
			Data:    "test",
		}

		// Attempt to marshal the test object into a DynamoDB item structure
		attrs, err := attributevalue.MarshalMapWithOptions(&data,
			func(eo *attributevalue.EncoderOptions) { eo.TagKey = "json" })
		Expect(err).ShouldNot(HaveOccurred())

		// Next, create our put-item input from our attribute data
		input := dynamodb.PutItemInput{
			Item:                        attrs,
			TableName:                   aws.String("TEST_TABLE"),
			ReturnConsumedCapacity:      types.ReturnConsumedCapacityNone,
			ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsNone,
			ReturnValues:                types.ReturnValueNone,
		}

		// Now, attempt to put the item to the database; this should not fail
		_, err = conn.PutItem(context.Background(), &input)
		Expect(err).ShouldNot(HaveOccurred())

		// Finally, attempt to retrieve the item as it exists in the database; this should not fail
		gOut, err := conn.db.GetItem(context.Background(), &dynamodb.GetItemInput{
			TableName: aws.String("TEST_TABLE"),
			Key: map[string]types.AttributeValue{
				"id":       &types.AttributeValueMemberS{Value: "test_id"},
				"sort_key": &types.AttributeValueMemberS{Value: "test|sort|key"}}})
		Expect(err).ShouldNot(HaveOccurred())

		// Attempt to extract our test object from the response
		var read *testObject
		err = attributevalue.UnmarshalMapWithOptions(gOut.Item, &read,
			func(eo *attributevalue.DecoderOptions) { eo.TagKey = "json" })
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the data on the test object
		Expect(read.ID).Should(Equal("test_id"))
		Expect(read.SortKey).Should(Equal("test|sort|key"))
		Expect(read.Data).Should(Equal("test"))
	})

	// Test that, if the inner GetItem request fails, then calling GetItem will return an error
	It("GetItem - Fails - Error", func() {

		// First, create our test database connection from our test config
		conn := createTestConnection(cfg)

		// Create our test object with some fake data
		data := testObject{
			ID:      "test_id",
			SortKey: "test|sort|key",
			Data:    "test",
		}

		// Attempt to marshal the test object into a DynamoDB item structure
		attrs, err := attributevalue.MarshalMapWithOptions(&data,
			func(eo *attributevalue.EncoderOptions) { eo.TagKey = "json" })
		Expect(err).ShouldNot(HaveOccurred())

		// Next, attempt to write a test object to the database; this should not fail
		_, err = conn.db.PutItem(context.Background(), &dynamodb.PutItemInput{
			Item:                        attrs,
			TableName:                   aws.String("TEST_TABLE"),
			ReturnConsumedCapacity:      types.ReturnConsumedCapacityNone,
			ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsNone})
		Expect(err).ShouldNot(HaveOccurred())

		// Now, attempt to get an item from the database; this should fail
		output, err := conn.GetItem(context.Background(), &dynamodb.GetItemInput{
			TableName:              aws.String("FAKE_TABLE"),
			ConsistentRead:         aws.Bool(false),
			ReturnConsumedCapacity: types.ReturnConsumedCapacityNone,
			Key: map[string]types.AttributeValue{
				"id":       &types.AttributeValueMemberS{Value: "test_id"},
				"sort_key": &types.AttributeValueMemberS{Value: "test|sort|key"},
			}})

		// Finally, verify the details of the error
		Expect(output).Should(BeNil())
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(And(
			HavePrefix("operation error DynamoDB: GetItem, https response error StatusCode: 400, RequestID: "),
			HaveSuffix(", ResourceNotFoundException: ")))
	})

	// Test that, if no failure occurs, then calling GetItem will result in the item being read
	// from the associated table in the database
	It("GetItem - No failures - Data returned", func() {

		// First, create our test database connection from our test config
		conn := createTestConnection(cfg)

		// Create our test object with some fake data
		data := testObject{
			ID:      "test_id",
			SortKey: "test|sort|key",
			Data:    "test",
		}

		// Attempt to marshal the test object into a DynamoDB item structure
		attrs, err := attributevalue.MarshalMapWithOptions(&data,
			func(eo *attributevalue.EncoderOptions) { eo.TagKey = "json" })
		Expect(err).ShouldNot(HaveOccurred())

		// Next, attempt to write a test object to the database; this should not fail
		_, err = conn.db.PutItem(context.Background(), &dynamodb.PutItemInput{
			Item:                        attrs,
			TableName:                   aws.String("TEST_TABLE"),
			ReturnConsumedCapacity:      types.ReturnConsumedCapacityNone,
			ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsNone})
		Expect(err).ShouldNot(HaveOccurred())

		// Now, attempt to get an item from the database; this should not fail
		output, err := conn.GetItem(context.Background(), &dynamodb.GetItemInput{
			TableName:              aws.String("TEST_TABLE"),
			ConsistentRead:         aws.Bool(false),
			ReturnConsumedCapacity: types.ReturnConsumedCapacityNone,
			Key: map[string]types.AttributeValue{
				"id":       &types.AttributeValueMemberS{Value: "test_id"},
				"sort_key": &types.AttributeValueMemberS{Value: "test|sort|key"}}})
		Expect(err).ShouldNot(HaveOccurred())

		// Finally, unmarshal the output response into a test object
		var written *testObject
		err = attributevalue.UnmarshalMapWithOptions(output.Item, &written,
			func(eo *attributevalue.DecoderOptions) { eo.TagKey = "json" })
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the data on the test object
		Expect(written.ID).Should(Equal("test_id"))
		Expect(written.SortKey).Should(Equal("test|sort|key"))
		Expect(written.Data).Should(Equal("test"))
	})

	// Test that, if the inner UpdateItem request fails, then calling UpdateItem will return an error
	It("UpdateItem - Fails - Error", func() {

		// First, create our test database connection from our test config
		conn := createTestConnection(cfg)

		// Create our test object with some fake data
		data := testObject{
			ID:      "test_id",
			SortKey: "test|sort|key",
			Data:    "test",
		}

		// Attempt to marshal the test object into a DynamoDB item structure
		attrs, err := attributevalue.MarshalMapWithOptions(&data,
			func(eo *attributevalue.EncoderOptions) { eo.TagKey = "json" })
		Expect(err).ShouldNot(HaveOccurred())

		// Next, attempt to write a test object to the database; this should not fail
		_, err = conn.db.PutItem(context.Background(), &dynamodb.PutItemInput{
			Item:                        attrs,
			TableName:                   aws.String("TEST_TABLE"),
			ReturnConsumedCapacity:      types.ReturnConsumedCapacityNone,
			ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsNone})
		Expect(err).ShouldNot(HaveOccurred())

		// Now, attempt to update the item we just put into the table; this should fail
		output, err := conn.UpdateItem(context.Background(), &dynamodb.UpdateItemInput{
			TableName:                   aws.String("FAKE_TABLE"),
			UpdateExpression:            aws.String("SET data = :val"),
			ExpressionAttributeValues:   map[string]types.AttributeValue{":val": &types.AttributeValueMemberS{Value: "test2"}},
			ReturnConsumedCapacity:      types.ReturnConsumedCapacityNone,
			ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsNone,
			ReturnValues:                types.ReturnValueAllNew,
			Key: map[string]types.AttributeValue{
				"id":       &types.AttributeValueMemberS{Value: "test_id"},
				"sort_key": &types.AttributeValueMemberS{Value: "test|sort|key"}}})

		// Finally, verify the details of the error
		Expect(output).Should(BeNil())
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(And(
			HavePrefix("operation error DynamoDB: UpdateItem, https response error StatusCode: 400, RequestID: "),
			HaveSuffix(", ResourceNotFoundException: ")))
	})

	// Test that, if no failure occurs, then calling UpdateItem will result in the item being
	// updated in the associated table in the database
	It("UpdateItem - No failures - Data updated", func() {

		// First, create our test database connection from our test config
		conn := createTestConnection(cfg)

		// Create our test object with some fake data
		data := testObject{
			ID:      "test_id",
			SortKey: "test|sort|key",
			Data:    "test",
		}

		// Attempt to marshal the test object into a DynamoDB item structure
		attrs, err := attributevalue.MarshalMapWithOptions(&data,
			func(eo *attributevalue.EncoderOptions) { eo.TagKey = "json" })
		Expect(err).ShouldNot(HaveOccurred())

		// Next, attempt to write a test object to the database; this should not fail
		_, err = conn.db.PutItem(context.Background(), &dynamodb.PutItemInput{
			Item:                        attrs,
			TableName:                   aws.String("TEST_TABLE"),
			ReturnConsumedCapacity:      types.ReturnConsumedCapacityNone,
			ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsNone})
		Expect(err).ShouldNot(HaveOccurred())

		// Now, attempt to update the item we just put into the table; this should not fail
		output, err := conn.UpdateItem(context.Background(), &dynamodb.UpdateItemInput{
			TableName:                   aws.String("TEST_TABLE"),
			UpdateExpression:            aws.String("SET #d = :val"),
			ExpressionAttributeNames:    map[string]string{"#d": "data"},
			ExpressionAttributeValues:   map[string]types.AttributeValue{":val": &types.AttributeValueMemberS{Value: "test2"}},
			ReturnConsumedCapacity:      types.ReturnConsumedCapacityNone,
			ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsNone,
			ReturnValues:                types.ReturnValueAllNew,
			Key: map[string]types.AttributeValue{
				"id":       &types.AttributeValueMemberS{Value: "test_id"},
				"sort_key": &types.AttributeValueMemberS{Value: "test|sort|key"}}})
		Expect(err).ShouldNot(HaveOccurred())

		// Finally, unmarshal the output response into a test object
		var updated *testObject
		err = attributevalue.UnmarshalMapWithOptions(output.Attributes, &updated,
			func(eo *attributevalue.DecoderOptions) { eo.TagKey = "json" })
		Expect(err).ShouldNot(HaveOccurred())

		// Verify the data on the test object
		Expect(updated.ID).Should(Equal("test_id"))
		Expect(updated.SortKey).Should(Equal("test|sort|key"))
		Expect(updated.Data).Should(Equal("test2"))
	})
})

// Helper function that creates a test connection from an AWS config for
func createTestConnection(cfg aws.Config) *DatabaseConnection {
	logger := utils.NewLogger("testd", "test")
	logger.Discard()
	return NewDatabaseConnection(cfg, logger,
		WithBackoffStart(1), WithBackoffEnd(5), WithBackoffMaxElapsed(10))
}

// Helper type that we'll use to test DynamoDB functionality
type testObject struct {
	ID      string `json:"id"`
	SortKey string `json:"sort_key"`
	Data    string `json:"data"`
}
