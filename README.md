# zctl
ZincObserve CLI tool for easy setup and installation of ZincObserve


Create IAM role, policy and s3 bucket on Amazon EKS and install
> zctl --name=zo1 --k8s=eks install
> zctl --name=zo1 --namespace=zns1 --k8s=eks install
> zctl --name=zo1 --namespace=zns1 --ingress-class=nginx --host=https://myurl.com --k8s=eks install

Will install minio and use it for object storage on any k8s and install
> zctl --name=zo1 --k8s=plain --storage=minio install

> zctl --name=zo1 --k8s=eks --bucket=bucket1 install

> zctl --name=zo1 --k8s=eks --bucket=bucket1 --iam-role=rolearn install

> zctl --name=zo1 delete

> zctl --name=zo1 --image=tag update

