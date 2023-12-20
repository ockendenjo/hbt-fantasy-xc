package main

import (
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsroute53targets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
)

func setupCloudFront(stack awscdk.Stack) {

	bucket := awss3.NewBucket(stack, jsii.String("Bucket"), &awss3.BucketProps{
		BucketName: jsii.String("fxc.hbt.ockenden.io"),
	})

	hostedZone := awsroute53.HostedZone_FromLookup(stack, jsii.String("HostedZone"), &awsroute53.HostedZoneProviderProps{DomainName: jsii.String("ockenden.io")})

	originAccessIdentity := awscloudfront.NewOriginAccessIdentity(stack, jsii.String("MyOriginAccessIdentity"), nil)
	s3Origin := awscloudfrontorigins.NewS3Origin(bucket, &awscloudfrontorigins.S3OriginProps{OriginAccessIdentity: originAccessIdentity})
	noCacheBehaviour := awscloudfront.BehaviorOptions{
		Origin:               s3Origin,
		Compress:             jsii.Bool(true),
		ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
		CachedMethods:        awscloudfront.CachedMethods_CACHE_GET_HEAD(),
		CachePolicy:          awscloudfront.CachePolicy_CACHING_DISABLED(),
	}

	oidcProviderArn := "arn:aws:iam::574363388371:oidc-provider/token.actions.githubusercontent.com"

	cicdRole := awsiam.NewRole(stack, jsii.String("CICDRole"), &awsiam.RoleProps{
		RoleName: jsii.String("HBTFantasyXCCICDRole"),
		AssumedBy: awsiam.NewWebIdentityPrincipal(&oidcProviderArn, &map[string]any{
			"StringEquals": map[string]string{
				"token.actions.githubusercontent.com:sub": "repo:ockendenjo/hbt-fantasy-xc:ref:refs/heads/main",
				"token.actions.githubusercontent.com:aud": "sts.amazonaws.com",
			},
		}),
		Description:        jsii.String("Used by GitHub actions to upload files to S3 bucket"),
		MaxSessionDuration: awscdk.Duration_Hours(jsii.Number(1)),
	})
	bucket.GrantReadWrite(cicdRole, jsii.String("*"))

	domainNames := []string{
		"fxc.hbt.ockenden.io",
	}

	hbtAndSubCert := "arn:aws:acm:us-east-1:574363388371:certificate/cb2155dc-e84e-4fb2-8d18-5030060c98e3"

	cfDist := awscloudfront.NewDistribution(stack, jsii.String("CFDistribution"), &awscloudfront.DistributionProps{
		DefaultBehavior: &awscloudfront.BehaviorOptions{
			Origin:               s3Origin,
			Compress:             jsii.Bool(true),
			ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
			AllowedMethods:       awscloudfront.AllowedMethods_ALLOW_GET_HEAD(),
			CachedMethods:        awscloudfront.CachedMethods_CACHE_GET_HEAD(),
			CachePolicy:          awscloudfront.CachePolicy_CACHING_OPTIMIZED(),
		},
		AdditionalBehaviors: &map[string]*awscloudfront.BehaviorOptions{
			"live.json":  &noCacheBehaviour,
			"index.html": &noCacheBehaviour,
		},
		PriceClass:  awscloudfront.PriceClass_PRICE_CLASS_100,
		HttpVersion: awscloudfront.HttpVersion_HTTP2_AND_3,
		DomainNames: func() *[]*string {
			a := []*string{}
			for _, name := range domainNames {
				a = append(a, &name)
			}
			return &a
		}(),
		Certificate:       awscertificatemanager.Certificate_FromCertificateArn(stack, jsii.String("Certificate"), &hbtAndSubCert),
		DefaultRootObject: jsii.String("index.html"),
	})
	cfDist.GrantCreateInvalidation(cicdRole)

	for _, name := range domainNames {
		recordId := "ARecord-" + strings.Split(name, ".")[0]
		awsroute53.NewARecord(stack, &recordId, &awsroute53.ARecordProps{
			Target:     awsroute53.RecordTarget_FromAlias(awsroute53targets.NewCloudFrontTarget(cfDist)),
			Zone:       hostedZone,
			RecordName: jsii.String(name + "."),
		})
	}
	cfDist.DomainName()

}
