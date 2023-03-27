package utils

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// CreateS3Bucket creates an S3 bucket with the specified name.
func CreateS3Bucket(bucketName, region string) error {
	if region == "" {
		region = "us-west-2"
	}

	fmt.Println(".Creating S3 Bucket............")

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region), // Specify the AWS region
	})
	if err != nil {
		return err
	}

	// Create a new S3 client
	s3Client := s3.New(sess)

	// Create the S3 bucket
	_, err = s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName), // Specify the bucket name
	})
	if err != nil {
		return err
	}

	// Wait for the bucket to exist
	err = s3Client.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return err
	}

	fmt.Println("Bucket created: ", bucketName)

	return nil
}

// delete s3 bucket
func DeleteS3Bucket(bucketName, region string) error {
	if region == "" {
		region = "us-west-2"
	}
	fmt.Println("DeleteS3Bucket............")
	// Create a new session with the AWS configuration
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		fmt.Println("error occured creating new aws session for deleting s3 bucket: ", err)
		return err
	}

	// Create a new S3 client
	s3Client := s3.New(sess)

	// Delete the S3 bucket
	_, err = s3Client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		fmt.Println("error occured deleting s3 bucket: ", bucketName, err)
		return err
	}

	// Wait until the bucket is deleted
	err = s3Client.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return err
	}

	return nil
}
