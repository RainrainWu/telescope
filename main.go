package main

import (
	"telescope/telescope"
)

func main() {

	atlas := telescope.NewAtlas("poetry.lock")
	atlas.Query()
	atlas.Report()
}
