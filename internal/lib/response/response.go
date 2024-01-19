package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const StatusOk = "OK"

func OK() Response {
	return Response{
		Status: StatusOk,
	}
}
