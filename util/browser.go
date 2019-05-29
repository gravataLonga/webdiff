package util

import (
	"bytes"
	"context"
	"image"
	"log"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// TakePicture it take a picture from URL
func TakePicture(url string) (img image.Image, err error) {
	return takePicture(url)
}

func takePicture(url string) (img image.Image, err error) {
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.WindowSize(1200, 600),
	}
	ectx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	// create context
	ctx, cancel := chromedp.NewContext(ectx)
	defer cancel()

	// run task list
	var buf []byte
	err = chromedp.Run(ctx, screenshot(url, &buf))
	if err != nil {
		return
	}
	img, _, err = image.Decode(bytes.NewReader(buf))
	return
}

func screenshot(url string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Sleep(5 * time.Second),
		chromedp.ActionFunc(func(ctxt context.Context) error {
			_, viewLayout, contentRect, err := page.GetLayoutMetrics().Do(ctxt)
			if err != nil {
				return err
			}

			v := page.Viewport{
				X:      contentRect.X,
				Y:      contentRect.Y,
				Width:  viewLayout.ClientWidth, // or contentRect.Width,
				Height: contentRect.Height,
				Scale:  1,
			}
			log.Printf("Capture %#v", v)
			buf, err := page.CaptureScreenshot().WithClip(&v).Do(ctxt)
			if err != nil {
				return err
			}
			*res = buf
			return err
		}),
	}
}
