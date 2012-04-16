package taobaoimg

import (
	"fmt"
)

func ExampleFetchImages() {
	Verbose = true
	imgs, _ := FetchImagesAndDecodeDimension("14754735064")
	fmt.Printf("%+v\n", len(imgs) > 0)
	//Output: true
}
