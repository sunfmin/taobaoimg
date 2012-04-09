package taobaoimg

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/sunfmin/integrationtest"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Image struct {
	URL    string
	Format string
	Width  int
	Height int
}

var descUrlRegexp = regexp.MustCompile(`http\:\/\/dsc\.taobaocdn\.com([\w\.\%\/]+)`)

var imgTagRegexp = regexp.MustCompile(`\<img[^<>]+(http[^\'\"]+jpg)`)

func FetchImagesAndDecodeDimension(num_iid string) (imgs []*Image, err error) {
	return fetch(num_iid, true)
}

func FetchImages(num_iid string) (imgs []*Image, err error) {
	return fetch(num_iid, false)
}

func fetch(num_iid string, dimension bool) (imgs []*Image, err error) {
	s := integrationtest.NewSession()
	// r, err := http.Get(fmt.Sprintf("http://item.taobao.com/item.htm?id=%s", num_iid))
	r, err := s.Get(fmt.Sprintf("http://item.taobao.com/item.htm?id=%s", num_iid))
	if err != nil {
		log.Printf("taobaoimg: get taobao item page error: %s\n", err)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	matchIndex := descUrlRegexp.FindReaderIndex(bufio.NewReader(io.TeeReader(r.Body, buf)))
	if len(matchIndex) == 0 {
		return
	}
	descURL := buf.String()[matchIndex[0]:matchIndex[1]]
	r, err = http.Get(descURL)
	if err != nil {
		log.Printf("taobaoimg: get taobao desc url %s error: %s\n", descURL, err)
		return
	}
	b, _ := ioutil.ReadAll(r.Body)
	matches := imgTagRegexp.FindAllStringSubmatch(string(b), -1)
	for _, match := range matches {
		img := &Image{
			URL: match[1],
		}

		if dimension {
			r, err = http.Get(match[1])
			if err != nil {
				log.Printf("taobaoimg: get taobao img url %s error: %s\n", match[1], err)
				continue
			}
			config, format, err := image.DecodeConfig(r.Body)
			if err != nil {
				log.Printf("taobaoimg: decode image %s error: %s\n", match[1], err)
				continue
			}
			img.Format = format
			img.Width = config.Width
			img.Height = config.Height
		}

		imgs = append(imgs, img)
	}

	return
}
