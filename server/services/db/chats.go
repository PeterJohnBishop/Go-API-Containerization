package db

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func CreateChat(client *dynamodb.Client, tableName string, item map[string]types.AttributeValue) error {
	_, err := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	return err
}

func GetChatById(client *dynamodb.Client, tableName, id string) (map[string]types.AttributeValue, error) {
	result, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}
	return result.Item, nil
}

func GetAllChats(client *dynamodb.Client, tableName string) ([]map[string]types.AttributeValue, error) {
	var items []map[string]types.AttributeValue
	var lastEvaluatedKey map[string]types.AttributeValue

	for {
		out, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
			TableName:         aws.String(tableName),
			ExclusiveStartKey: lastEvaluatedKey,
		})
		if err != nil {
			return nil, err
		}

		items = append(items, out.Items...)

		if out.LastEvaluatedKey == nil {
			break
		}
		lastEvaluatedKey = out.LastEvaluatedKey
	}

	return items, nil
}

func UpdateChat(client *dynamodb.Client, tableName string, chat Chat) error {
	updateBuilder := expression.UpdateBuilder{}
	updatedFields := 0 // Track the number of fields updated

	if chat.Users != nil {
		updateBuilder = updateBuilder.Set(expression.Name("users"), expression.Value(chat.Users))
		updatedFields++
	}
	if chat.Messages != nil {
		updateBuilder = updateBuilder.Set(expression.Name("messages"), expression.Value(chat.Messages))
		updatedFields++
	}
	updateBuilder = updateBuilder.Set(expression.Name("active"), expression.Value(time.Now().Unix()))
	updatedFields++

	// Ensure at least one field is being updated
	if updatedFields == 0 {
		return fmt.Errorf("must update at least one field")
	}

	expr, err := expression.NewBuilder().WithUpdate(updateBuilder).Build()
	if err != nil {
		fmt.Println("Error in expression builder:", err)
		return err
	}

	_, err = client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: chat.ID},
		},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              types.ReturnValueUpdatedNew,
	})

	if err != nil {
		fmt.Println("Error in client updater:", err)
	}
	return err
}

func DeleteChat(client *dynamodb.Client, tableName, id string) error {
	_, err := client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	return err
}
