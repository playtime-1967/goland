# VPC
With Amazon Virtual Private Cloud (Amazon VPC), you can launch AWS resources in a logically isolated virtual network that you've defined.
The VPC has one subnet in each of the Availability Zones in the Region, EC2 instances in each subnet, and an internet gateway to allow communication between the resources in your VPC and the internet.

Your AWS account includes a default VPC in each AWS Region. Your default VPCs are configured such that you can immediately start launching and connecting to EC2 instances.

A subnet is a range of IP addresses in your VPC. You launch AWS resources, such as Amazon EC2 instances, into your subnets. You can connect a subnet to the internet, other VPCs, and your own data centers, and route traffic to and from your subnets using route tables.


network address translation (NAT) 
maps multiple private IPv4 addresses to a single public IPv4 address. 


# EC2
An EC2 instance is a virtual server in the AWS Cloud
provides on-demand, scalable computing capacity in AWS


# Security groups
A virtual firewall that allows you to specify the protocols, ports, and source IP ranges that can reach your instances, and the destination IP ranges to which your instances can connect.
---------

# DB
list of sql dbs:
aws rds describe-db-instances --region us-east-1 --profile glob-admin

SOME INFO:
"DBInstanceIdentifier": "rds-dev-tmm-mssql",
"Address": "rds-dev-tmm-mssql.czk2mki6q3b0.us-east-1.rds.amazonaws.com",
"Port": 1433,
VpcSecurityGroupId
"sg-0aa8f65135bf02346",
"sg-0e5e1e54692125364"

Subnets
subnet-07209f8fcb09a45ff
subnet-02b9a9ed7957443fc

find out db username & password:
1-List of secrets and  extract The Sql instance ARN from the result:
aws secretsmanager list-secrets --profile glob-admin  --region us-east-1

2- retrive the SecretString uincluding  username & password
aws secretsmanager get-secret-value --secret-id arn:aws:secretsmanager:us-east-1:747444338602:secret:sm-dev-db-password-mssql-p7o84nsj-shXAoy --profile glob-admin  --region us-east-1

The Result:
{
    "ARN": "arn:aws:secretsmanager:us-east-1:747444338602:secret:sm-dev-db-password-mssql-p7o84nsj-shXAoy",
    "Name": "sm-dev-db-password-mssql-p7o84nsj",
    "VersionId": "terraform-20250219094008636100000003",
    "SecretString": "{\n  \"username\": \"themodernmilkman\",\n  \"password\": \"kznv9zqYOp4kjWVJ\",\n  \"endpoint\": \"rds-dev-tmm-mssql.czk2mki6q3b0.us-east-1.rds.amazonaws.com\"\n}\n",
    "VersionStages": [
        "AWSCURRENT"
    ],
    "CreatedDate": "2025-02-19T09:40:10.187000+00:00"
}


# EC2 instance
Set up an EC2 instance in the same VPC (Virtual Private Cloud) and subnet group (or at least with routing access) as your SQL Server RDS.
Then you:
Connect to the EC2 from your local machine (via SSH for Linux, RDP for Windows).
From the EC2, connect to the SQL Server (because EC2 is inside the same network as RDS).

#  Launch an EC2 instance
Go to AWS EC2 Console → "Launch Instance".
Choose an OS:
Windows if you want to use SQL Server Management Studio (SSMS) inside the EC2.
Linux if you want to connect via command-line tools.
Select an instance type (t2.micro for testing).
Networking:
Place it in the same VPC and preferably the same subnet as your SQL Server RDS.
Assign a security group that allows inbound RDP (Windows) or SSH (Linux) from your IP.
Allow outbound traffic to the RDS instance’s port (1433 for SQL Server).
Create and download the key pair (for SSH or RDP login).
Launch the instance.

Configure Instance Details
Network: Select the VPC where your DB and EKS are.
Subnet: Select a public subnet — usually, these subnets have a route to the Internet Gateway.
Auto-assign Public IP: Enable (this is crucial to connect from your laptop).

Add Inbound rules:
Type: SSH
Protocol: TCP
Port: 22
/*--

Open an SSH client.
Locate your private key file. The key used to launch this instance is us-milkman.pem
Run this command, if necessary, to ensure your key is not publicly viewable.
chmod 400 "us-milkman.pem"
Connect to your instance using its Public DNS:
ec2-13-222-167-37.compute-1.amazonaws.com
Example:
ssh -i "us-milkman.pem" ubuntu@ec2-13-222-167-37.compute-1.amazonaws.com


after connceting: prepare the env:
Install AWS CLI v2 on Ubuntu (WSL)
sudo apt update && sudo apt upgrade -y
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
sudo apt install unzip -y
unzip awscliv2.zip

(instead of aws configure sso) aws configure  
AWS Access Key ID [None]: ASIA24BZHP6VIMRO5YD3
AWS Secret Access Key [None]: hINSfoPZAU1Zz9z1k82ZvfFAdNCQ963OstzTxBL2
Default region name [None]: eu-west-2
Default output format [None]: json

** You also need to put AWS_SESSION_TOKEN manualy into  ~/.aws/credentials (picture 5)
copy this value from aws console.

Verify the setup:
aws eks list-clusters --region us-east-1
or
aws eks describe-cluster --name eks-dev --region us-east-1


- Install Kubectl

Update kubeconfig (so kubectl can access the cluster):
aws eks update-kubeconfig --name eks-dev --region us-east-1 --profile glob-admin

kubectl get pods -n dev
E0818 15:28:25.606352   15446 memcache.go:265] "Unhandled Error" err="couldn't get current server API group list: Get \"https://C1E5A885E6D62B59641121F1D46C1FFF.gr7.us-east-1.eks.amazonaws.com/api?timeout=32s\": dial tcp 172.30.9.221:443: i/o timeout"

Get error becauase need to define acces of ec2 to eks:
- create a role - picture 6
----------
Update RDS Security Group
Go to your RDS instance → Connectivity & Security.
In the VPC Security Groups, edit inbound rules:
Add Custom TCP Rule for port 1433 (SQL Server default).
Source: your EC2’s security group ID (not an IP address — this is more secure).
