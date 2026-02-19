provider "aws" {
  region = var.region
}

module "network" {
  source = "./network"
}

module "eks" {
  source = "./eks"
  vpc_id = module.network.vpc_id
}

module "datastores" {
  source = "./datastores"
  vpc_id = module.network.vpc_id
}
