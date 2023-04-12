package responses

type BasicResponse struct {
	Ok        bool   `json:"ok"`
	Error     string `json:"error,omitempty"`
	ErrorCode int    `json:"error_code,omitempty"`
}

func Ok() BasicResponse {
	return BasicResponse{
		Ok: true,
	}
}

func Error(err string, code int) BasicResponse {
	return BasicResponse{
		Ok:        false,
		Error:     err,
		ErrorCode: code,
	}
}
