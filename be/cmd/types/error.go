package types

type HttpErrorResponse struct {
	Error string `json:"error" example:"Bad Request"`
}

type HttpNotFound struct {
	Error string `json:"error" example:"Not Found"`
}
