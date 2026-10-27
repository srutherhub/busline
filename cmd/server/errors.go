package server

type ErrorCode string

const (
	IllegalArgument     string = "ILLEGAL_ARGUMENT"
	ResourceNotFound    string = "RESOURCE_NOT_FOUND"
	InternalServerError string = "INTERNAL_SERVER_ERROR"
)
