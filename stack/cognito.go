package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/jsii-runtime-go"
)

func setupUserPool(stack awscdk.Stack) {
	pool := awscognito.NewUserPool(stack, jsii.String("HBTFantasyXCPool"), &awscognito.UserPoolProps{
		AccountRecovery: awscognito.AccountRecovery_EMAIL_ONLY,
		AutoVerify: &awscognito.AutoVerifiedAttrs{
			Email: jsii.Bool(true),
		},
		Email: awscognito.UserPoolEmail_WithSES(&awscognito.UserPoolSESOptions{
			FromEmail: jsii.String("techgeek@huntersbogtrotters.com"),
			FromName:  jsii.String("HBT Tech Geek"),
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
			EmailBody:    jsii.String("OFY - Your verification code is {####}"),
			EmailStyle:   awscognito.VerificationEmailStyle_CODE,
		},
	})

	awscognito.NewUserPoolDomain(stack, jsii.String("PoolDomain"), &awscognito.UserPoolDomainProps{
		CognitoDomain: &awscognito.CognitoDomainOptions{
			DomainPrefix: jsii.String("hbt-fantasy-xc"),
		},
		UserPool: pool,
	})

	client := awscognito.NewUserPoolClient(stack, jsii.String("PoolClient"), &awscognito.UserPoolClientProps{
		AccessTokenValidity: awscdk.Duration_Days(jsii.Number(1)),
		GenerateSecret:      jsii.Bool(false), //Not for a public client
		UserPoolClientName:  jsii.String("PoolPublicClient"),
		UserPool:            pool,
		OAuth: &awscognito.OAuthSettings{
			CallbackUrls: jsii.Strings("https://fxc.hbt.ockenden.io/signin"),
			LogoutUrls:   jsii.Strings("https://fxc.hbt.ockenden.io/signout"),
			Flows: &awscognito.OAuthFlows{
				AuthorizationCodeGrant: jsii.Bool(false),
				ImplicitCodeGrant:      jsii.Bool(true),
			},
		},
	})

	awscdk.NewCfnOutput(stack, jsii.String("PoolClientId"), &awscdk.CfnOutputProps{
		Description: jsii.String("Cognito user pool client ID"),
		Value:       client.UserPoolClientId(),
	})
}
