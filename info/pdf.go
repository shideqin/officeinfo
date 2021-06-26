package info

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"officeinfo/base"
	"strconv"
	"strings"

	pdf "github.com/unidoc/unipdf/v3/model"
)

type PdfInfo struct {
	Page   int `json:"page"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (i *Info) GetPdfInfo(args []string) ([]byte, error) {
	var info = make([]byte, 0)
	if len(args) < 2 {
		return info, errors.New("miss parameters")
	}
	var input = args[1]
	var pageNum = 1
	if len(args) > 2 {
		pageNum, _ = strconv.Atoi(args[2])
	}
	var exitChan = make(chan bool)
	defer close(exitChan)
	var buffer = &bytes.Buffer{}
	var _, err = base.CURL2Reader(input, "GET", map[string]string{}, strings.NewReader(""), buffer, exitChan)
	if err != nil {
		return info, err
	}
	return i.pdfPageProperties(bytes.NewReader(buffer.Bytes()), pageNum)
}

func (i *Info) pdfPageProperties(rs io.ReadSeeker, pageNum int) ([]byte, error) {
	var info = make([]byte, 0)
	pdfReader, err := pdf.NewPdfReader(rs)
	if err != nil {
		return info, err
	}
	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return info, err
	}
	// Try decrypting with an empty one.
	if isEncrypted {
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return info, err
		}
		if !auth {
			return info, errors.New("encrypted unable to access")
		}
	}
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return info, err
	}
	// If invalid pageNum.
	if pageNum <= 0 || pageNum > numPages {
		pageNum = 1
	}
	page, err := pdfReader.GetPage(pageNum)
	if err != nil {
		return info, err
	}
	mBox, err := page.GetMediaBox()
	if err != nil {
		return info, err
	}
	var pdfInfo = &PdfInfo{}
	pdfInfo.Width = int(mBox.Urx - mBox.Llx)
	pdfInfo.Height = int(mBox.Ury - mBox.Lly)
	pdfInfo.Page = numPages
	return json.Marshal(pdfInfo)
}
