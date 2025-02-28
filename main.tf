terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

resource "aws_vpc" "vpc" {
  cidr_block = "110.98.192.0/18"
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.vpc.id
}