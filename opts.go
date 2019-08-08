package freyja

import (
	"bytes"
	"github.com/valyala/fasthttp"
	"strconv"
)

type Opts struct {
	width uint32
	height uint32
	quality uint32
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

		quality = uint32(quality)

	} else {
		quality = defaultQuality
	}


	return Opts{
		width: uint32(width),
		height: uint32(height),
		quality: quality,
	}, fasthttp.StatusOK

}