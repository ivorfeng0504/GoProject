package dataapi

type DataApiResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token"`
	Data []DataApiData `json:"data"`
}
