package info

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"io"
	"net/http"
	"officeinfo/base"
	"strings"

	"golang.org/x/image/webp"
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
	contentType, err := i.fileContentType(bytes.NewReader(buffer.Bytes()))
	if err != nil {
		return info, err
	}
	var img image.Image
	if contentType == "image/webp" {
		img, err = webp.Decode(buffer)
	} else {
		img, _, err = image.Decode(buffer)
	}
	if err != nil {
		return info, err
	}
	b := img.Bounds()
	var imgInfo = &ImgInfo{}
	imgInfo.Width = b.Max.X
	imgInfo.Height = b.Max.Y
	return json.Marshal(imgInfo)
}

func (i *Info) fileContentType(r io.ReadSeeker) (string, error) {
	// 读取前 512 个字节
	buf := make([]byte, 512)
	_, err := r.Read(buf)
	_, _ = r.Seek(0, 0)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buf)
	return contentType, nil
}
