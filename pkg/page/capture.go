package page

import (
	"context"
	"encoding/base64"
	"log"

	"github.com/chromedp/chromedp"
)

func Capture(encodedurl string) ([]byte, error) {
	var err error
	var decoded []byte

	log.Println(encodedurl)
	if decoded, err = base64.URLEncoding.DecodeString(encodedurl); err == nil {
		url := string(decoded)
		log.Printf("capture URL:%s\n", url)

		ctx, cancel := chromedp.NewContext(
			context.Background(),
			// chromedp.WithDebugf(log.Printf),
		)
		defer cancel()

		var buf []byte

		if err = chromedp.Run(ctx, chromedp.Tasks{
			chromedp.EmulateViewport(1200, 630),
			chromedp.Navigate(url),
			chromedp.FullScreenshot(&buf, 90),
		}); err != nil {
			return nil, err
		}

		return buf, nil
	}
	return nil, err
}
