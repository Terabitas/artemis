package config

import "strings"

// Config type
type Config struct {
	Verbosity int
	IP        string
	Port      string
	Secret    string

	CORSAllowedOrigins     []string
	CORSAllowedMethods     []string
	CORSAllowedHeaders     []string
	CORSExposedHeaders     []string
	CORSAllowCredentials   bool
	CORSMaxAge             int
	CORSOptionsPassThrough bool
	CORSDebug              bool
}

// StringToSlice return slice from "x,y,z"
func StringToSlice(in string) []string {
	rez := []string{}
	for _, val := range strings.Split(in, ",") {
		rez = append(rez, val)
	}

	return rez
}
