package lib

import (
	"log"
	"os"
	"strconv"
)

// GetEnvStr returns string for environment string with default
func GetEnvStr(name, def string) (res string) {
	res = def

	s := os.Getenv(name)
	if s == "" {
		log.Printf("%s not present in environment", name)
	} else {
		res = s
	}

	log.Printf("%s set to %s", name, res)
	return
}

// GetEnvInt returns string for environment numeric with default
func GetEnvInt(name string, def int) (res int) {
	res = def

	s := os.Getenv(name)
	if s == "" {
		log.Printf("%s not present in environment", name)
	} else {
		i, err := strconv.Atoi(s)
		if err == nil {
			res = i
		} else {
			log.Printf("%s not numeric %s", name, s)
		}
	}

	log.Printf("%s set to %d", name, res)
	return
}

// Trunc ...
func Trunc(in string) (res string) {
	res = in
	if len(res) > 80 {
		res = res[:80] + "..."
	}

	return
}
