package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
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

	awscognito.NewUserPool(stack, jsii.String("HBTFantasyXCPool"), &awscognito.UserPoolProps{
		AccountRecovery: awscognito.AccountRecovery_EMAIL_ONLY,
		AutoVerify: &awscognito.AutoVerifiedAttrs{
			Email: jsii.Bool(true),
		},
		Email: awscognito.UserPoolEmail_WithSES(&awscognito.UserPoolSESOptions{
			FromEmail: jsii.String("noreply@hbt.ockenden.io"),
			FromName:  jsii.String("HBT Fantasy XC"),
			ReplyTo:   jsii.String("techgeek@huntersbogtrotters.com"),
		}),
		Mfa:               awscognito.Mfa_OFF,
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
		SelfSignUpEnabled: jsii.Bool(true),
		SignInAliases: &awscognito.SignInAliases{
			Email:             jsii.Bool(true),
			Phone:             jsii.Bool(false),
			PreferredUsername: jsii.Bool(false),
			Username:          jsii.Bool(false),
		},
		StandardAttributes: nil,
		UserPoolName:       jsii.String("HBTFantasyXCPool"),
		UserVerification: &awscognito.UserVerificationConfig{
			EmailSubject: jsii.String("Registration for HBT Fantasy XC"),
			EmailBody:    jsii.String("Your verification code is {####}"),
			EmailStyle:   awscognito.VerificationEmailStyle_CODE,
		},
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
