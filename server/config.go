package main

import "strings"

type Config struct {
	Host             string `october:"HOST"`
	Port             int    `october:"PORT"`
	ValidBuckets     string `october:"VALID_BUCKETS"`
	MaxContentLength int64  `october:"MAX_CONTENT_LENGTH"`
}

func (c *Config) MustExtractValidBuckets() [][]byte {

	split := strings.Split(c.ValidBuckets, ",")

	var result [][]byte

	for _, s := range split {
		result = append(result, []byte(s))
	}

	return result
}
