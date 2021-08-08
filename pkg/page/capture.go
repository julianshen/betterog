package page

import (
	"context"
	"encoding/base64"
	"log"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func FullScreenshotInViewport(res *[]byte, quality int) chromedp.EmulateAction {
	if res == nil {
		panic("res cannot be nil")
	}
	return chromedp.ActionFunc(func(ctx context.Context) error {
		format := page.CaptureScreenshotFormatJpeg

		var err error
		// capture screenshot
		*res, err = page.CaptureScreenshot().
			WithCaptureBeyondViewport(true).
			WithFormat(format).
			WithQuality(int64(quality)).WithClip(&page.Viewport{
			X:      0,
			Y:      0,
			Width:  1200,
			Height: 630,
			Scale:  1,
		}).Do(ctx)

		if err != nil {
			return err
		}
		return nil
	})
}

func Capture(encodedurl string) ([]byte, error) {
	var err error
	var decoded []byte

	if decoded, err = base64.RawURLEncoding.DecodeString(encodedurl); err == nil {
		url := string(decoded)
		log.Printf("capture URL:%s\n", url)

		ctx, _ := chromedp.NewExecAllocator(context.Background(), chromedp.NoSandbox)
		ctx, cancel := chromedp.NewContext(
			ctx,
			// chromedp.WithDebugf(log.Printf),
		)

		defer cancel()

		var buf []byte

		if err = chromedp.Run(ctx, chromedp.Tasks{
			chromedp.EmulateViewport(1200, 630),
			chromedp.Navigate(url),
			FullScreenshotInViewport(&buf, 90),
		}); err != nil {
			return nil, err
		}

		return buf, nil
	}
	return nil, err
}
