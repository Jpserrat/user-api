package response

type ApiResponse struct {
	Body   interface{} `json:"body"`
	Status int         `json:"status"`
}
