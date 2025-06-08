package config

import (
	"os"
	"strconv"
)

func getString(key, def string) string {
	res := os.Getenv(key)
	if res == "" {
		return def
	}

	return res
}

func getBool(key string, def bool) bool {
	str := getString(key, strconv.FormatBool(def))
	res, err := strconv.ParseBool(str)
	if err != nil {
		return def
	}

	return res
}
