package main

// Import key modules.

type AddressSpace struct {
	AddressPrefixes []string `json:"addressPrefixes"`
}
type Properties struct {
	AddressSpace      AddressSpace `json:"addressSpace"`
	ProvisioningState string       `json:"provisioningState"` //  "Succeeded",
	// "enableDdosProtection": false,
	//"subnets": [],
	//"virtualNetworkPeerings": []
}

type ReponseBase struct {
	ResourceId string     `json:"id"`
	Location   string     `json:"location"`
	Name       string     `json:"name"`
	Properties Properties `json:"properties"`
}
