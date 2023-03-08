terraform {

    backend "s3" {}

}

provider "aws" {

    profile = "blank-devops"
    region  = "eu-central-1"
    
}

module "ecr" {
    source  = "app.terraform.io/blank/ecr/aws"
    version = "0.0.1"
    name = "go-template"
}