package model

type NFTCollection struct {
	ID              string `json:"id"`
	Blockchain      string `json:"blockchain"`
	Network         string `json:"network"`
	Name            string `json:"name"`
	ContractAddress string `json:"contract_address"`
}

type CollectionRegistrationRequestData struct {
	Blockchain      string `json:"blockchain"`
	Network         string `json:"network"`
	Name            string `json:"name,omitempty"`
	ContractAddress string `json:"contract_address"`
}

type CollectionDeregistrationRequestData struct {
	ID string `json:"id"`
}

type CollectionListResponse struct {
	Items []NFTCollection `json:"items"`
}
