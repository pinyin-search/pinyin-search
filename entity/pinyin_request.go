package entity

type PinYinRequest struct {
	Tenant    string `json:"tenant"`
	IndexName string `json:"indexName"`
	DataId    string `json:"dataId"`
	Data      string `json:"data"`
}
