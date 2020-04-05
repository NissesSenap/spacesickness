package storage

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const bucketname string = "something-ec909d91-5794-4acd-ba49-53ec2e2c1f56"

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

// GetPreSign will generate a presigned url
func (ose ObjectStorage) GetPreSign() {
	req, _ := ose.Svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketname),
		Key:    aws.String("debbuild-19.11.0-1.el8.noarch.rpm"),
	})
	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		log.Println("Failed to sign request", err)
	}

	fmt.Printf("The URL is: \n %v \n", urlStr)
}

// GetBucketPolicyStatus returns a bool if the system is public or not
func (ose ObjectStorage) GetBucketPolicyStatus() (bool, error) {
	status, err := ose.Svc.GetBucketPolicyStatus(&s3.GetBucketPolicyStatusInput{
		Bucket: aws.String(bucketname),
	})
	if err != nil {
		fmt.Println("Can't get the status...")
		return false, err
	}
	fmt.Println(status.String())
	/*
		if status {
			return true, nil
		}
		return false, nil
	*/
	return false, nil
}

// GetBucketPolicy grabs the current bucket policy
func (ose ObjectStorage) GetBucketPolicy() {
	result, err := ose.Svc.GetBucketPolicy(&s3.GetBucketPolicyInput{
		Bucket: aws.String(bucketname),
	})
	if err != nil {
		// Special error handling for the when the bucket doesn't
		// exists so we can give a more direct error message from the CLI.
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				fmt.Printf("Bucket %q does not exist", bucketname)
			case "NoSuchBucketPolicy":
				fmt.Printf("Bucket %q does not have a policy", bucketname)
			}
		}
		fmt.Printf("Unable to get bucket %q policy, %v", bucketname, err)
	}
	out := bytes.Buffer{}
	policyStr := aws.StringValue(result.Policy)
	if err := json.Indent(&out, []byte(policyStr), "", "  "); err != nil {
		fmt.Printf("Failed to pretty print bucket policy, %v", err)
	}

	fmt.Printf("%q's Bucket Policy:\n", bucketname)
	fmt.Println(out.String())

}

// ReadAllPolicyBucket allows anonymous users  to read a specific bucket
func (ose ObjectStorage) ReadAllPolicyBucket() {
	readOnlyAnonUserPolicy := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Sid":       "PublicReadForGetBucketObjects",
				"Effect":    "Allow",
				"Principal": "*",
				"Action": []string{
					"s3:GetObject",
				},
				"Resource": []string{
					fmt.Sprintf("arn:aws:s3:::%s/*", bucketname),
				},
			},
		},
	}
	policy, err := json.Marshal(readOnlyAnonUserPolicy)
	if err != nil {
		exitErrorf("Failed to marshal policy, %v", err)
	}
	fmt.Println(string(policy))

	// Set the bucket policy
	_, err = ose.Svc.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketname),
		Policy: aws.String(string(policy)),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == s3.ErrCodeNoSuchBucket {
			// Special error handling for the when the bucket doesn't
			// exists so we can give a more direct error message from the CLI.
			exitErrorf("Bucket %q does not exist", bucketname)
		}
		exitErrorf("Unable to set bucket %q policy, %v", bucketname, err)
	}

	fmt.Printf("Successfully set bucket %q's policy\n", bucketname)

}

// GetAllObjects grabs all the items in a bucket
func (ose ObjectStorage) GetAllObjects() *s3.ListObjectsV2Output {
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucketname),
		MaxKeys: aws.Int64(1000), // Default value is 1000, need to look in to pageination long-term
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
		return nil
	}
	return result
}

// GetObjectList is returns a list of strings containing all names of the objects
func (ose ObjectStorage) GetObjectList(baseList *s3.ListObjectsV2Output) []string {
	var myList []string
	for _, b := range baseList.Contents {
		myList = append(myList, *b.Key)
	}
	return myList
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
