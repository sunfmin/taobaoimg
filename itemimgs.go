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

func DecodeImage(imgu string) (img *Image, err error) {
	r3, err := http.Get(imgu)
	defer close(r3)
	if err != nil {
		log.Printf("taobaoimg: get taobao img url %s error: %s\n", imgu, err)
		return
	}
	config, format, err := image.DecodeConfig(r3.Body)
	if err != nil {
		log.Printf("taobaoimg: decode image %s error: %s\n", imgu, err)
		return
	}
	img = &Image{}
	img.URL = imgu
	img.Format = format
	img.Width = config.Width
	img.Height = config.Height
	return
}

func fetch(num_iid string, dimension bool) (imgs []*Image, err error) {
	s := integrationtest.NewSession()
	// r, err := http.Get(fmt.Sprintf("http://item.taobao.com/item.htm?id=%s", num_iid))
	r1, err := s.Get(fmt.Sprintf("http://item.taobao.com/item.htm?id=%s", num_iid))
	defer close(r1)

	if err != nil {
		log.Printf("taobaoimg: get taobao item page error: %s\n", err)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	matchIndex := descUrlRegexp.FindReaderIndex(bufio.NewReader(io.TeeReader(r1.Body, buf)))
	if len(matchIndex) == 0 {
		return
	}
	descURL := buf.String()[matchIndex[0]:matchIndex[1]]
	r2, err := http.Get(descURL)
	if err != nil {
		log.Printf("taobaoimg: get taobao desc url %s error: %s\n", descURL, err)
		return
	}
	b, _ := ioutil.ReadAll(r2.Body)
	defer close(r2)
	matches := imgTagRegexp.FindAllStringSubmatch(string(b), -1)
	for _, match := range matches {
		var img *Image
		if dimension {
			img, err = DecodeImage(match[1])
		} else {
			img = &Image{
				URL: match[1],
			}
		}
		if img != nil {
			imgs = append(imgs, img)
		}
	}

	return
}

func close(r *http.Response) {
	if r != nil && r.Body != nil {
		r.Body.Close()
	}
}
