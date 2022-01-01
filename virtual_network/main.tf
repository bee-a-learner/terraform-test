
resource "azurerm_virtual_network" "this" {
  name                = var.resource_name
  location            = var.location
  resource_group_name = var.resource_group_name
  address_space       = var.address_space
  tags = var.tags

}


output "resource_name" {
  value = azurerm_virtual_network.this.name
}

output "resource_id" {
  value = azurerm_virtual_network.this.id
}

output "vnet_address_space" {
  value = azurerm_virtual_network.this.address_space
}

output "resource_group_name" {
  value = azurerm_virtual_network.this.resource_group_name
}