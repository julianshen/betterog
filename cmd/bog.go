package main

import (
	"image/color"
	"log"
	"os"

	betterog "github.com/julianshen/betterog/pkg/server"
	"github.com/julianshen/text2img"
	"gopkg.in/go-playground/colors.v1"
)

func env(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultVal
	}
	return val
}

func main() {
	addr := ":" + env("PORT", "8888")
	fontPath := env("FONTPATH", "fonts/SourceHanSansTC-VF.ttf")
	backgroundImage := env("BACKGROUND_IMAGE_PATH", "")
	backgroundColor := env("BACKGROUND_COLOR", "")
	textColor := env("TEXT_COLOR", "")

	param := text2img.Params{
		FontPath: fontPath,
		Width:    800,
		Height:   400,
	}

	if backgroundImage != "" {
		log.Println("BGImage:" + backgroundImage)
		param.BackgroundImagePath = backgroundImage
	}

	if backgroundColor != "" {
		log.Println("BGCOLOR:" + backgroundColor)
		if c, err := colors.ParseHEX(backgroundColor); err == nil {
			rgba := c.ToRGBA()
			param.BackgroundColor = color.RGBA{
				rgba.R, rgba.G, rgba.B, uint8(rgba.A),
			}
		}
	}

	if textColor != "" {
		log.Println("TEXT_COLOR:" + textColor)
		if c, err := colors.ParseHEX(textColor); err == nil {
			rgba := c.ToRGBA()
			param.TextColor = color.RGBA{
				rgba.R, rgba.G, rgba.B, uint8(rgba.A),
			}
		}
	}

	if server, err := betterog.NewServer(addr, param); err == nil {
		server.Start()
	} else {
		log.Fatal(err)
	}
}
