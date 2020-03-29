package storage

import (
	"crypto/tls"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ObjectStorage this is the base struct that is used everywhere.
type ObjectStorage struct {
	AwsAccess   string
	AwsSecret   string
	AwsEndpoint string
	AwsRegion   string
	client      *http.Client
}

//NewS3 creates the method
func NewS3(ose *ObjectStorage) {
	// Creating custom client that can ignore TLS
	// For some reason it isn't built in to the tool...
	// https://github.com/aws/aws-sdk-go/issues/2404

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpclient := &http.Client{Transport: tr}

	ose.client = httpclient
}

// CreateSession returns the *S3 svc
func (ose ObjectStorage) CreateSession() *s3.S3 {
	s := session.New(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(ose.AwsAccess, ose.AwsSecret, ""),
		Endpoint:         aws.String(ose.AwsEndpoint),
		Region:           aws.String(ose.AwsRegion),
		DisableSSL:       aws.Bool(false),
		S3ForcePathStyle: aws.Bool(true),
		//LogLevel:         aws.LogLevel(aws.LogDebug | aws.LogDebugWithHTTPBody | aws.LogDebugWithRequestRetries | aws.LogDebugWithRequestErrors | aws.LogDebugWithSigning),
		//Logger:           aws.NewDefaultLogger(),
		HTTPClient: ose.client,
	})
	svc := s3.New(s)
	return svc
}
