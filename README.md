# zctl

ZincObserve CLI tool for easy setup and installation of ZincObserve

Create IAM role, policy and s3 bucket on Amazon EKS and install

> zctl --name=zo1 --k8s=eks install
> zctl --name=zo1 --namespace=zns1 --k8s=eks install
> zctl --name=zo1 --namespace=zns1 --ingress-class=nginx --host=https://myurl.com --k8s=eks install

> zctl --name=zo1 --k8s=eks --ingester=2 --router=2 --querier=2 install

Will install minio and use it for object storage on any k8s and install

> zctl --name=zo1 --k8s=plain --storage=minio install

> zctl --name=zo1 --k8s=eks --bucket=bucket1 install

> zctl --name=zo1 --k8s=eks --bucket=bucket1 --iam-role=rolearn install

> zctl --name=zo1 delete

> zctl --name=zo1 --image=tag update

# Steps

1. Check if OIDC provider exists
1. Create a bucket
1. Create IAM policy
1. Create IAM role
1. Install helm chart

Mimnimum Items to specify

1. image.repository
1. image.tag
1. ServiceAccount.Annotations["eks.amazonaws.com/role-arn"]
1. config.ZOS3BUCKETNAME
1.

# Install

1. Check if a configmap exists with the name zincobserve-setup. If the configmap exists then a setup has already been done.
1. If configmap does not exist then proceed
1. get namespace and releaseName
1. Generate a random install identifier.
1. Create a configmap with name zincobserve-setup with releasename and install identifier
1. bucketname should be

configmap should have following values . e.g

setup_data: {
"identifier": "15096452",
"release_name": "zo1"
"bucket_name": "zinc-observe-15096452-dev2-zo1",
"iam_role: "zinc-observe-15096452-dev2-zo1"
}

# What is working today

## Install on EKS

> zctl install --k8s=eks --name=zo1

## Uninstall a release in EKS

> zctl uninstall --k8s=eks --name=zo1


# AWS

## Install
> zctl install --k8s=eks --name=zo1  

## Uninstall

> zctl uninstall --k8s=eks --name=zo1

# GCP

## Install

1. Get project ID

> gcloud auth application-default login 
> gcloud config get-value project

> zctl install --k8s=gke --name=zo1 --namespace=zo1 --gcp_project_id=zinc1-342016

This will create:

1. A GCS bucket
1. An IAM service account
1. Grant IAM service account permissions to the GCS bucket
1. Create HMAC keys (S3 access key and secret) fir the service account




