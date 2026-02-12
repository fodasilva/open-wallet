package utils

import (
	"os"
	"strconv"
)

type Env struct {
	GETRpmLimit    int
	POSTRpmLimit   int
	PATCHRpmLimit  int
	DELETERpmLimit int
}

var (
	EnvConfig Env
)

func LoadEnvs() {
	getRpmLimitStr := os.Getenv("GET_RPM_RATE_LIMIT")
	postRpmLimitStr := os.Getenv("POST_RPM_RATE_LIMIT")
	patchRpmLimitStr := os.Getenv("PATCH_RPM_RATE_LIMIT")
	deleteRpmLimitStr := os.Getenv("DELETE_RPM_RATE_LIMIT")

	if getRpmLimitStr == "" {
		panic("GET_RPM_RATE_LIMIT environment variable cannot be empty")
	}

	if postRpmLimitStr == "" {
		panic("POST_RPM_RATE_LIMIT environment variable cannot be empty")
	}

	if patchRpmLimitStr == "" {
		panic("PATCH_RPM_RATE_LIMIT environment variable cannot be empty")
	}

	if deleteRpmLimitStr == "" {
		panic("DELETE_RPM_RATE_LIMIT environment variable cannot be empty")
	}

	getRpmLimit, err := strconv.Atoi(getRpmLimitStr)
	if err != nil {
		panic("GET_RPM_RATE_LIMIT must be a valid integer")
	}

	postRpmLimit, err := strconv.Atoi(postRpmLimitStr)
	if err != nil {
		panic("POST_RPM_RATE_LIMIT must be a valid integer")
	}

	patchRpmLimit, err := strconv.Atoi(patchRpmLimitStr)
	if err != nil {
		panic("PATCH_RPM_RATE_LIMIT must be a valid integer")
	}

	deleteRpmLimit, err := strconv.Atoi(deleteRpmLimitStr)
	if err != nil {
		panic("DELETE_RPM_RATE_LIMIT must be a valid integer")
	}

	EnvConfig = Env{
		GETRpmLimit:    getRpmLimit,
		POSTRpmLimit:   postRpmLimit,
		PATCHRpmLimit:  patchRpmLimit,
		DELETERpmLimit: deleteRpmLimit,
	}
}
