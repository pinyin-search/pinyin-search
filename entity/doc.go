package entity

type Doc struct {
	Id    string      `json:"id"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}
