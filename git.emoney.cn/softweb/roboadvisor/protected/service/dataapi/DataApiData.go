package dataapi

type DataApiData struct {
	Success bool     `json:"success"`
	ErrMsg  string     `json:"errmsg"`
	Index   int        `json:"index"`
	Result  [][]string `json:"result"`
}
