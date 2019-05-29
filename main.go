package main

import (
	"flag"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/gravatalonga/webdiff/util"
)

var (
	picture, output string
	compare         bool
)

func main() {

	log.SetPrefix("WebDiff - ")

	flag.StringVar(&picture, "picture", "", "Take a picture of url.")
	flag.BoolVar(&compare, "compare", false, "Is to compare two image?")
	flag.StringVar(&output, "output", "output.png", "The name of file to output")
	flag.StringVar(&output, "o", "output.png", "The name of file to output")
	flag.Parse()

	// take a picture
	// go run main.go -picture "https://www.google.pt" -o "google.png"
	// compare picture
	// go run main.go -o "diff.png" -compare "google.png" "google1.png"

	// Compare
	if compare {
		args := flag.Args()
		if len(args) < 2 {
			log.Fatal("In order to compare you must provider which file to compare.")
		}

		img := util.DiffImage(mustLoadImage(args[0]), mustLoadImage(args[1]))
		mustSaveImage(img, output)
		log.Println("Done")
		return
	}

	// Take a picture
	img, err := util.TakePicture(picture)
	if err != nil {
		log.Fatal("Got an error. ", err)
	}
	mustSaveImage(img, output)
}

func mustSaveImage(img image.Image, output string) {
	f, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	png.Encode(f, img)
}

func mustOpen(filename string) *os.File {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	return f
}

func mustLoadImage(filename string) image.Image {
	f := mustOpen(filename)
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	return img
}
