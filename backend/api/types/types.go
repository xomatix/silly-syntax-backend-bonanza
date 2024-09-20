package types

type ResponseMessage struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    any    `json:"data"`
}
