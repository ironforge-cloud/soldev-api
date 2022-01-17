package main

import (
	"api/internal/utils"
	"fmt"
)

func main() {
	asdf := utils.GetImageIfExist("https://dev.to/remiix/getting-website-meta-tags-with-node-js-1li5")
	fmt.Println(asdf)
}
