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
  profile = "terraform"
}

resource "aws_vpc" "vpc" {
  cidr_block = "110.98.192.0/18"
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.vpc.id
}

data "aws_availability_zones" "available" {}

data "aws_availability_zone" "available" {
  for_each = toset(data.aws_availability_zones.available.names)
  name     = each.value
}

locals {
  az_map = {
    for zone in data.aws_availability_zone.available :
    zone.name => zone.zone_id
  }
}

resource "aws_subnet" "private_subnet_az_1e" {
  vpc_id     = aws_vpc.vpc.id
  cidr_block = "110.98.192.0/26"
  availability_zone_id = local.az_map["us-east-1e"]
}

resource "aws_subnet" "private_subnet_az_1b" {
  vpc_id     = aws_vpc.vpc.id
  cidr_block = "110.98.192.64/26"
  availability_zone_id = local.az_map["us-east-1b"]
}

resource "aws_subnet" "private_subnet_az_1c" {
  vpc_id     = aws_vpc.vpc.id
  cidr_block = "110.98.192.128/26"
  availability_zone_id = local.az_map["us-east-1c"]
}

resource "aws_subnet" "public_subnet_az_1e" {
  vpc_id     = aws_vpc.vpc.id
  cidr_block = "110.98.193.64/26"
  availability_zone_id = local.az_map["us-east-1e"]
}

resource "aws_subnet" "public_subnet_az_1b" {
  vpc_id     = aws_vpc.vpc.id
  cidr_block = "110.98.193.128/26"
  availability_zone_id = local.az_map["us-east-1b"]
}

resource "aws_subnet" "public_subnet_az_1c" {
  vpc_id     = aws_vpc.vpc.id
  cidr_block = "110.98.193.192/26"
  availability_zone_id = local.az_map["us-east-1c"]
}
