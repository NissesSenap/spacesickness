package main

import (
	"fmt"
	"os"

	"github.com/NissesSenap/spacesickness/storage"
	"github.com/aws/aws-sdk-go/aws"
)

func main() {

	// TODO add something that checks if you have a config file the default ~/.aws/credentials
	// Env always wins
	awsAccessEnv := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretEnv := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsEndpointEnv := os.Getenv("AWS_ENDPOINT")
	awsRegionEnv := os.Getenv("AWS_REGION")

	u := storage.ObjectStorage{
		AwsAccess:   awsAccessEnv,
		AwsSecret:   awsSecretEnv,
		AwsEndpoint: awsEndpointEnv,
		AwsRegion:   awsRegionEnv,
	}

	storage.NewS3(&u)

	svc := u.CreateSession()

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
