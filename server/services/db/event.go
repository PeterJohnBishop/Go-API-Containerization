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

func CreateEvent(client *dynamodb.Client, tableName string, item map[string]types.AttributeValue) error {
	_, err := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	return err
}

func GetEventById(client *dynamodb.Client, tableName, id string) (map[string]types.AttributeValue, error) {
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

func GetAllEvents(client *dynamodb.Client, tableName string) ([]map[string]types.AttributeValue, error) {
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

func UpdateEvent(client *dynamodb.Client, tableName string, event Event) error {
	updateBuilder := expression.UpdateBuilder{}
	updatedFields := 0 // Track the number of fields updated

	if event.Name != "" {
		updateBuilder = updateBuilder.Set(expression.Name("name"), expression.Value(event.Name))
		updatedFields++
	}
	if event.Description != "" {
		updateBuilder = updateBuilder.Set(expression.Name("description"), expression.Value(event.Description))
		updatedFields++
	}
	if event.AssignedTo != nil {
		updateBuilder = updateBuilder.Set(expression.Name("assigned_to"), expression.Value(event.AssignedTo))
		updatedFields++
	}
	if event.StartDate != 0 {
		updateBuilder = updateBuilder.Set(expression.Name("start_date"), expression.Value(event.StartDate))
		updatedFields++
	}
	if event.EndDate != 0 {
		updateBuilder = updateBuilder.Set(expression.Name("end_date"), expression.Value(event.EndDate))
		updatedFields++
	}
	if event.LocationName != "" {
		updateBuilder = updateBuilder.Set(expression.Name("location_name"), expression.Value(event.LocationName))
		updatedFields++
	}
	if event.LocationAddress != "" {
		updateBuilder = updateBuilder.Set(expression.Name("location_address"), expression.Value(event.LocationAddress))
		updatedFields++
	}
	if event.LocationLong != 0 {
		updateBuilder = updateBuilder.Set(expression.Name("location_long"), expression.Value(event.LocationLong))
		updatedFields++
	}
	if event.LocationLat != 0 {
		updateBuilder = updateBuilder.Set(expression.Name("location_lat"), expression.Value(event.LocationLat))
		updatedFields++
	}
	if event.Notes != "" {
		updateBuilder = updateBuilder.Set(expression.Name("notes"), expression.Value(event.Notes))
		updatedFields++
	}
	if event.FirstNotification != 0 {
		updateBuilder = updateBuilder.Set(expression.Name("first_notification"), expression.Value(event.FirstNotification))
		updatedFields++
	}
	if event.SecondNotification != 0 {
		updateBuilder = updateBuilder.Set(expression.Name("second_notification"), expression.Value(event.SecondNotification))
		updatedFields++
	}
	updateBuilder = updateBuilder.Set(expression.Name("active"), expression.Value(event.Active))
	updatedFields++
	updateBuilder = updateBuilder.Set(expression.Name("updated_at"), expression.Value(time.Now().Unix()))
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
			"id": &types.AttributeValueMemberS{Value: event.ID},
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

func DeleteEvent(client *dynamodb.Client, tableName, id string) error {
	_, err := client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	return err
}
