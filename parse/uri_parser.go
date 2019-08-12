package parse

import (
	"bytes"

	"github.com/valyala/fasthttp"
)

type ParsedURI struct {
	Bucket     []byte
	BucketPath []byte
	Opts       Opts
}

func (p *ParsedURI) S3Uri() []byte {
	return append(slashBytes, append(p.Bucket, append(slashBytes, p.BucketPath...)...)...)
}

var validSuffixs = [][]byte{
	[]byte(".jpeg"),
	[]byte(".jpg"),
}

var slashBytes = []byte("/")

func ParseURI(uri []byte, validBuckets [][]byte, defaultQuality uint32) (parsed ParsedURI, respStatus int) {

	idxModifier := 0
	if bytes.HasPrefix(uri, slashBytes) {
		idxModifier += 1
	}

	split := bytes.SplitN(uri, slashBytes, 3+idxModifier)

	if len(split) != 3+idxModifier {
		// Split should be exactly 3 always if its valid

		return ParsedURI{}, fasthttp.StatusBadRequest
	}

	opts := split[0+idxModifier]
	bucket := split[1+idxModifier]
	bucketPath := split[2+idxModifier]

	if len(validBuckets) == 0 {
		return ParsedURI{}, fasthttp.StatusBadRequest
	}

	if len(validSuffixs) == 0 {
		return ParsedURI{}, fasthttp.StatusBadRequest
	}

	validSuffix := false

	for _, vs := range validSuffixs {
		if bytes.HasSuffix(bucketPath, vs) {
			bucketPath = bucketPath[:len(bucketPath)-len(vs)]
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
		Bucket:     bucket,
		BucketPath: bucketPath,
		Opts:       parsedOpts,
	}, fasthttp.StatusOK

}
