package freyja

import (
	"bytes"
	"github.com/valyala/fasthttp"
)


type ParsedURI struct {
	bucket []byte
	bucketPath []byte
	opts Opts
}


var validSuffixs = [][]byte{
	[]byte(".jpeg"),
	[]byte(".jpg"),
}

var slashBytes = []byte("/")

func parseURI (uri []byte, validBuckets [][]byte, defaultQuality uint32) (parsed ParsedURI, respStatus int) {

	split := bytes.SplitN(uri, slashBytes, 3)

	if len(split) != 3 {
		// Split should be exactly 3 always if its valid

		return ParsedURI{}, fasthttp.StatusBadRequest
	}

	opts := split[0]
	bucket := split[1]
	bucketPath := split[2]

	if len(validBuckets) == 0 {
		return ParsedURI{}, fasthttp.StatusBadRequest
	}

	if len(validSuffixs) == 0 {
		return ParsedURI{}, fasthttp.StatusBadRequest
	}

	validSuffix := false

	for _, vs := range validSuffixs {
		if bytes.HasSuffix(bucketPath, vs) {
			validSuffix = true
			break
		}
	}

	if !validSuffix {
		return ParsedURI{}, fasthttp.StatusBadRequest
	}

	validBucket := false
	for _, vb := range validBuckets {

		if bytes.Compare(bucket, vb) == 0 {

			validBucket = true
			break
		}
	}

	if !validBucket {
		return ParsedURI{}, fasthttp.StatusBadRequest
	}


	parsedOpts, resp := parseOpts(opts, defaultQuality)
	if resp != fasthttp.StatusOK {
		return ParsedURI{}, resp
	}


	return ParsedURI{
		bucket: bucket,
		bucketPath: bucketPath,
		opts: parsedOpts,
	}, fasthttp.StatusOK

}