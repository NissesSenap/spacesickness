package storage

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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
	Svc         *s3.S3
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

// GetObjects grabs all the items in a bucket.
func (ose ObjectStorage) GetObjects() {
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String("something-ec909d91-5794-4acd-ba49-53ec2e2c1f56"),
		MaxKeys: aws.Int64(2),
	}

	result, err := ose.Svc.ListObjectsV2(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				fmt.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	fmt.Println(result)
}
