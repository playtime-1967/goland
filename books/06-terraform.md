Terraform
build, change, and version infrastructure safely and efficiently
manage your infrastructure in a safe, consistent, and repeatable way
manages your infrastructure's lifecycle

Terraform's configuration language is declarative, meaning that it describes the desired end-state for your infrastructure, in contrast to procedural programming languages that require step-by-step instructions to perform tasks. Terraform providers automatically calculate dependencies between resources to create or destroy them in the correct order.

Install:
HashiCorp distributes Terraform
---

terraform.tf 
The terraform {} block configures Terraform itself, including which providers to install, and which version of Terraform to use to provision your infrastructure.
including:
required_providers
required_version

main.tf
When you write a new Terraform configuration, we recommend defining your provider blocks and other primary infrastructure in main.tf as a best practice.
including:
provider "aws" 
data "aws_ami" "ubuntu"
resource "aws_instance" "app_server"

Format configuration:
terraform fmt

Initialize your workspace:
terraform init

Validate configuration:
terraform validate

Create infrastructure:
Terraform makes changes to your infrastructure in two steps.
It creates an execution plan for the changes it will make. Review this plan to ensure that Terraform will make the changes you expect.
Once you approve the execution plan, Terraform applies those changes using your workspace's providers.
Detect and resolve any unexpected problems with your configuration before Terraform makes changes to your infrastructure.
cmd:
terraform apply

Inspect state:
Terraform wrote data about your infrastructure into a file called terraform.tfstate
terraform state list

workspace's entire state(infra detials):
terraform show

Your state file can include sensitive information about your infrastructure, such as passwords or security keys, so you must store your state file securely and restrict access to only those who need to manage your infrastructure with Terraform.

Storing your state remotely using HCP Terraform 
By default, Terraform creates your state file locally.

Variables and outputs
Input variables let you parametrize the behavior of your Terraform configuration.
Variables and outputs also allow you to integrate your Terraform workspaces with other automation tools
variables.tf:
variable "instance_type" {
  description = "The EC2 instance's type."
  type        = string
  default     = "t2.micro"
}
Update the instance configuration in main.tf to refer to these variables instead of hard-coding the argument values:
instance_type = var.instance_type

Run a Terraform plan without applying it to see what would happen if you changed your EC2 instance type from t2.micro to t2.large:
terraform plan -var instance_type=t2.large
result: ~ instance_type                        = "t2.micro" -> "t2.large"

Output values
You can define output values to expose data about the resources you create and consume their values with other automation tools or workflows.
outputs.tf:
output "instance_hostname" {
  description = "Private DNS name of the EC2 instance."
  value       = aws_instance.app_server.private_dns
}

terraform apply
terraform output

Modules
Reusable sets of configuration.
module "vpc" { }
Since Terraform automatically resolves dependencies within your configuration, you can organize your configuration blocks in any order you like
Whenever you add a new module to your configuration, you will need to install it by re-initializing the workspace:
terraform init



terraform destroy
