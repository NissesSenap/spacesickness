package s3service

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func CreateSession(awsAccess, awsSecret, awsEndpoint, awsRegion string) *S3 {
// Creating custom client that can ignore TLS
	// For some reason it isn't built in to the tool...
	// https://github.com/aws/aws-sdk-go/issues/2404
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	s := session.New(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(awsAccess, awsSecret, ""),
		Endpoint:         aws.String(awsEndpoint),
		Region:           aws.String(awsRegion),
		DisableSSL:       aws.Bool(false),
		S3ForcePathStyle: aws.Bool(true),
		//LogLevel:         aws.LogLevel(aws.LogDebug | aws.LogDebugWithHTTPBody | aws.LogDebugWithRequestRetries | aws.LogDebugWithRequestErrors | aws.LogDebugWithSigning),
		//Logger:           aws.NewDefaultLogger(),
		HTTPClient: client,
	})
	svc := s3.New(s)
	reteurn svc
}
