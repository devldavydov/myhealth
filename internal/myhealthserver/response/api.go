package response

type APIResponse struct {
	Error string `json:"error"`
	Data  any    `json:"data"`
}

func NewErrorAPIResponse(err string) *APIResponse {
	return &APIResponse{Error: err}
}

func NewDataAPIResponse(data any) *APIResponse {
	return &APIResponse{Data: data}
}

func NewOKAPIResponse() *APIResponse {
	return &APIResponse{Data: "ok"}
}
