package parse

import (
	"bytes"
	"github.com/valyala/fasthttp"
	"strconv"
)

type Opts struct {
	Width   uint32
	Height  uint32
	Quality uint32
}

const maxQuality = 100
const maxWidth = 1920*2
const maxHeight = 1080*2

var optsSep = []byte("x")

func parseOpts(opts []byte, defaultQuality uint32) (opt Opts, respStatus int) {

	split := bytes.SplitN(opts, optsSep, 3)


	width, err := strconv.ParseUint(string(split[0]), 10, 32)
	if err != nil {
		return Opts{}, fasthttp.StatusBadRequest
	}

	if width > maxWidth {
		return Opts{}, fasthttp.StatusBadRequest
	}

	height, err := strconv.ParseUint(string(split[1]), 10, 32)
	if err != nil {
		return Opts{}, fasthttp.StatusBadRequest
	}

	if height > maxHeight {
		return Opts{}, fasthttp.StatusBadRequest
	}


	var quality uint32

	if len(split) == 3 {
		qual, err := strconv.ParseUint(string(split[2]), 10, 32)
		if err != nil {
			return Opts{}, fasthttp.StatusBadRequest
		}

		if qual > maxQuality {
			return Opts{}, fasthttp.StatusBadRequest
		}

		quality = uint32(qual)

	} else {
		quality = defaultQuality
	}


	return Opts{
		Width:   uint32(width),
		Height:  uint32(height),
		Quality: quality,
	}, fasthttp.StatusOK

}