package responses

type CommonResponse struct {
	StatusCode int32
	StatusMsg  string `json:"status_msg,omitempty"`
}
