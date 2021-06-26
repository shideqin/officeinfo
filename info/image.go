package info

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"officeinfo/base"
	"strings"
)

type ImgInfo struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (i *Info) GetImageInfo(args []string) ([]byte, error) {
	var info = make([]byte, 0)
	if len(args) < 2 {
		return info, errors.New("miss parameters")
	}
	var input = args[1]
	var exitChan = make(chan bool)
	defer close(exitChan)
	var buffer = &bytes.Buffer{}

	var _, err = base.CURL2Reader(input, "GET", map[string]string{}, strings.NewReader(""), buffer, exitChan)
	if err != nil {
		return info, err
	}
	img, _, err := image.Decode(buffer)
	if err != nil {
		return info, err
	}
	b := img.Bounds()
	var imgInfo = &ImgInfo{}
	imgInfo.Width = b.Max.X
	imgInfo.Height = b.Max.Y
	return json.Marshal(imgInfo)
}
