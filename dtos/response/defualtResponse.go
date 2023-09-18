package dtos

type DefaultResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Date    string      `json:"date"`
	Data    interface{} `json:"data"`
}

type StatusResponse string

const (
	Success StatusResponse = "success"
	Failed  StatusResponse = "failed"
)
