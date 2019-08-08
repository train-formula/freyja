package main

import (
	"bytes"
	"github.com/valyala/fasthttp"
	"testing"
)

var testValidBuckets = [][]byte{
	[]byte("bucket1"),
	[]byte("bucket2"),
}

var defaultQuality uint32 = 33

func Test_parseURI_noQualityOpts_parsesOK(t *testing.T) {


	okBasicOptsWithSlash := []byte("/192x168/bucket1/Bucket/my/path.bmp.jpeg")

	parsed, resp := parseURI(okBasicOptsWithSlash, testValidBuckets, defaultQuality)
	if resp != fasthttp.StatusOK {
		t.Errorf("Response code should be StatusOK")
	}

	if bytes.Compare(parsed.Bucket, []byte("bucket1")) != 0 {
		t.Errorf("Invalid bucket")
	}

	if bytes.Compare(parsed.BucketPath, []byte("Bucket/my/path.bmp.jpeg")) != 0 {
		t.Errorf("Invalid bucket")
	}

	if parsed.Opts.Width != 192 {
		t.Errorf("Invalid width opt %d", parsed.Opts.Width)
	}

	if parsed.Opts.Height != 168 {
		t.Errorf("Invalid height opt %d", parsed.Opts.Height)
	}

	if parsed.Opts.Quality != defaultQuality {
		t.Errorf("Invalid quality opt %d", parsed.Opts.Quality)
	}
}


func Test_parseURI_qualityOpts_parsesOK(t *testing.T) {


	okBasicOptsWithSlash := []byte("/192x168x55/bucket1/Bucket/my/path.bmp.jpeg")

	parsed, resp := parseURI(okBasicOptsWithSlash, testValidBuckets, defaultQuality)
	if resp != fasthttp.StatusOK {
		t.Errorf("Response code should be StatusOK")
	}

	if bytes.Compare(parsed.Bucket, []byte("bucket1")) != 0 {
		t.Errorf("Invalid bucket")
	}

	if bytes.Compare(parsed.BucketPath, []byte("Bucket/my/path.bmp.jpeg")) != 0 {
		t.Errorf("Invalid bucket")
	}

	if parsed.Opts.Width != 192 {
		t.Errorf("Invalid width opt %d", parsed.Opts.Width)
	}

	if parsed.Opts.Height != 168 {
		t.Errorf("Invalid height opt %d", parsed.Opts.Height)
	}

	if parsed.Opts.Quality != 55 {
		t.Errorf("Invalid quality opt %d", parsed.Opts.Quality)
	}
}
