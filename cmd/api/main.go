package main

import (
	"fmt"
	"os"

	"github.com/NissesSenap/spacesickness/s3service"
	"github.com/aws/aws-sdk-go/aws"
)

func main() {

	// TODO add something that checks if you have a config file the default ~/.aws/credentials
	// Env always wins
	awsAccess := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	awsRegion := os.Getenv("AWS_REGION")

	svc := s3service.CreateSession(awsAccess, awsSecret, awsEndpoint, awsRegion)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
