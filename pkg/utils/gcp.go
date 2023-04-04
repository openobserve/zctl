package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	admin "cloud.google.com/go/iam/admin/apiv1"
	"cloud.google.com/go/iam/admin/apiv1/adminpb"
	"cloud.google.com/go/iam/apiv1/iampb"
	"cloud.google.com/go/storage"
)

func SetupGCP(setupData SetupData) (SetupData, error) {
	// 1. Create bucket
	bucketName := "zinc-observe-" + setupData.Identifier + "-" + setupData.ReleaseName
	setupData.BucketName = bucketName

	err := CreateBucket(setupData.GCPProjectId, setupData.BucketName)
	if err != nil {
		fmt.Println(err)
	}

	// 2. Create service account
	serviceAccount, err := CreateGCPServiceAccount(setupData.GCPProjectId, setupData.Identifier)
	if err != nil {
		fmt.Println(err)
	}

	setupData.ServiceAccount = serviceAccount.Email

	// 3. Grant access to service account to the bucket
	err = GrantAllAccessToBucket(setupData.GCPProjectId, setupData.BucketName, serviceAccount.Email)
	if err != nil {
		fmt.Println(err)
	}

	// 4. Create HMAC key
	key, err := CreateHMACKey(setupData.GCPProjectId, serviceAccount.Email)
	if err != nil {
		fmt.Println(err)
	}

	setupData.S3AccessKey = key.AccessID
	setupData.S3SecretKey = key.Secret

	return setupData, nil
}

// DeleteGCPServiceAccount deletes a service account
func DeleteGCPServiceAccount(setupData SetupData) error {
	ctx := context.Background()
	client, err := admin.NewIamClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Delete the service account.
	if err := client.DeleteServiceAccount(ctx, &adminpb.DeleteServiceAccountRequest{
		Name: fmt.Sprintf("projects/%s/serviceAccounts/%s", setupData.GCPProjectId, setupData.ServiceAccount),
	}); err != nil {
		log.Fatalf("Failed to delete service account: %v", err)
	}

	return nil
}

func DeleteGCSBucket(setupData SetupData) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Bucket(setupData.BucketName).Delete(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func GrantAllAccessToBucket(projectID, bucketName, serviceAccountEmail string) error {
	{
		ctx := context.Background()

		role := "roles/storage.objectAdmin" // or any other role you'd like to grant

		client, err := storage.NewClient(ctx)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		defer client.Close()

		policy, err := client.Bucket(bucketName).IAM().V3().Policy(ctx)
		if err != nil {
			log.Fatalf("Failed to get bucket policy: %v", err)
		}

		newBinding := iampb.Binding{
			Role:    role,
			Members: []string{fmt.Sprintf("serviceAccount:%s", serviceAccountEmail)},
		}

		policy.Bindings = append(policy.Bindings, &newBinding)

		if err := client.Bucket(bucketName).IAM().V3().SetPolicy(ctx, policy); err != nil {
			log.Fatalf("Failed to update bucket policy: %v", err)
		}

		fmt.Printf("Successfully granted %s access to %s for service account %s\n", role, bucketName, serviceAccountEmail)
	}

	return nil
}

// CreateHMACKey creates a new HMAC key using the given project and service account.
func CreateHMACKey(projectID string, serviceAccountEmail string) (*storage.HMACKey, error) {
	ctx := context.Background()

	// Initialize client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	key, err := client.CreateHMACKey(ctx, projectID, serviceAccountEmail)
	if err != nil {
		return nil, fmt.Errorf("CreateHMACKey: %v", err)
	}

	fmt.Println("Created HMAC key: ", key)

	return key, nil
}

func CreateBucket(projectID, bucketName string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	bucketAttrs := &storage.BucketAttrs{
		Name:     bucketName,
		Location: "US", // Replace with your desired location
	}

	if err := bucket.Create(ctx, projectID, bucketAttrs); err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
		return fmt.Errorf("Failed to create bucket: %v", err)
	}

	fmt.Printf("Bucket %s created successfully\n", bucketName)

	return nil
}

func CreateGCPServiceAccount(projectID, serviceAccountName string) (*adminpb.ServiceAccount, error) {

	fmt.Println("Creating service account: ", serviceAccountName)
	ctx := context.Background()
	client, err := admin.NewIamClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create IAM client: %v", err)
	}
	defer client.Close()

	fmt.Println("Created IAM client: ", client)

	name := fmt.Sprintf("projects/%s", projectID)

	fmt.Println("Service account name: ", name)

	req := &adminpb.CreateServiceAccountRequest{
		// Parent:    fmt.Sprintf("projects/%s", projectID),
		Name:      name,
		AccountId: "zinc-observe-" + serviceAccountName,
		ServiceAccount: &adminpb.ServiceAccount{
			DisplayName: "zinc-observe-" + serviceAccountName,
			Description: "Allows creating HMAC keys for Zinc Observe for GCS bucket",
		},
	}

	fmt.Println("Created service account request: ", req)

	serviceAccount, err := client.CreateServiceAccount(ctx, req)
	if err != nil {
		fmt.Println("Failed to create service account: ", err)
		return nil, fmt.Errorf("failed to create service account: %v", err)
	}

	fmt.Println("Created service account: ", serviceAccount)

	return serviceAccount, nil
}
