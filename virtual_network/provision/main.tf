terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
    }
  }
}
provider "azurerm" {
  features {}
}


variable "location" {
  type    = string
  default = "uksouth"
}

variable "postfix" {
  type    = string
  default = "kjxhiasuymce"
}

variable "prefix" {
  type    = string
  default = "vnet"
}

variable "resource_group_name" {
  type    = string
  default = "example-resources"
}

variable "address_space" {
  default = "10.0.0.0/16"
}

variable "tags" {
  default = {
      env = "dev"
    }
}

module "core_vnet" {
    source = "../"
    resource_name       = "${var.prefix}-${var.postfix}"
    location            = var.location
    resource_group_name = var.resource_group_name
    address_space       = [var.address_space]
    tags = var.tags
}


output "resource_name" {
  value = module.core_vnet.resource_name
}

output "resource_id" {
  value = module.core_vnet.resource_id
}

output "vnet_address_space" {
  value = module.core_vnet.vnet_address_space
}

output "resource_group_name" {
  value = module.core_vnet.resource_group_name
}
