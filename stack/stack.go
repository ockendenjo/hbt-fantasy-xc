package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type StackProps struct {
	awscdk.StackProps
}

func NewStack(scope constructs.Construct, id string, props *StackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	setupCloudFront(stack)

	awsdynamodb.NewTable(stack, jsii.String("users"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("email"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode:         awsdynamodb.BillingMode_PAY_PER_REQUEST,
		TableClass:          awsdynamodb.TableClass_STANDARD,
		TimeToLiveAttribute: jsii.String("TTL"),
		TableName:           jsii.String("HbtFantasyXcUsers"),
	})

	awsdynamodb.NewTable(stack, jsii.String("challenges"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("challenge_id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode:         awsdynamodb.BillingMode_PAY_PER_REQUEST,
		TableClass:          awsdynamodb.TableClass_STANDARD,
		TimeToLiveAttribute: jsii.String("TTL"),
		TableName:           jsii.String("HbtFantasyXcChallenges"),
	})

	awsdynamodb.NewTable(stack, jsii.String("sessions"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("session_id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode: awsdynamodb.BillingMode_PAY_PER_REQUEST,
		TableClass:  awsdynamodb.TableClass_STANDARD,
		TableName:   jsii.String("HbtFantasyXcSessions"),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewStack(app, "HbtFantasyXCStack", &StackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String("574363388371"),
		Region:  jsii.String("eu-west-1"),
	}
}
