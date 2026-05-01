package util

var StatusMessages = map[int]string{
	400: "bad request",
	401: "unauthorized",
	403: "forbidden",
	404: "not found",
	429: "too many requests",
	500: "internal server error",
}
