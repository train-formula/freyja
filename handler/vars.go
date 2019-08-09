package handler

import "time"

var healthResp = []byte("OK")
var cacheControlHeader = []byte("Cache-Control")
var noCacheControlValue = []byte("private, no-cache, no-store, must-revalidate, max-age=0")
var pragmaHeader = []byte("Pragma")
var noCachePragmaValue = []byte("no-cache")

var internalServerError = []byte("Internal Server Error")
var s3InternalServerError = []byte("Failed to build s3 request through varnish")

var faviconURI = []byte("/favicon.ico")
var healthCheckURI = []byte("/_health")

var etagEnclosing = "\""
var etagHeaderString = "ETag"
var etagHeader = []byte(etagHeaderString)

var expiresHeader = []byte("Expires")
var expiresDuration = time.Second * 2592000

var okCacheControlValue = []byte("public, must-revalidate, max-age=2592000")
