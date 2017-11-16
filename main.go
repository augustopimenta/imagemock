package main

import (
	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"strconv"
	"fmt"
	"regexp"
	"os"
	"net/http"
	"github.com/lucasb-eyer/go-colorful"
	"math"
	"image"
	"bytes"
	"image/png"
)

const (
	maxImageSize = 5000
	minImageSize = 50
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		html(c, http.StatusOK, renderHelp())
	})

	r.GET("/:size", func(c *gin.Context) {
		r := regexp.MustCompile(`^(\d+)([xX](\d+))?$`)
		matches := r.FindStringSubmatch(c.Param("size"))

		if len(matches) == 0 {
			html(c, http.StatusInternalServerError, renderError("Parâmetros inválidos"))
			return
		}

		width, _ := strconv.Atoi(matches[1])

		height, err := strconv.Atoi(matches[3])

		if err != nil {
			height = width
		}

		if width > maxImageSize || height > maxImageSize {
			html(c, http.StatusInternalServerError, renderError("A imagem deve ter no máximo 5000 x 5000"))
			return
		}

		if width < minImageSize || height < minImageSize {
			html(c, http.StatusInternalServerError, renderError("A imagem deve ter no mínimo 50 x 50"))
			return
		}

		dc := gg.NewContext(width, height)
		dc.DrawRectangle(0,0, float64(width), float64(height))

		background := c.DefaultQuery("bg", "666666")
		bgColor, err := colorful.Hex("#" + background)

		if err != nil {
			html(c, http.StatusInternalServerError, renderError("Cor de fundo incorreta"))
			return
		}

		dc.SetColor(bgColor)
		dc.Fill()

		fColor := c.DefaultQuery("c", "FFFFFF")
		fontColor, err := colorful.Hex("#" + fColor)

		if err != nil {
			html(c, http.StatusInternalServerError, renderError("Cor do texto incorreta"))
			return
		}

		points := math.Min(float64(width), float64(height)) / 6

		dc.SetColor(fontColor);
		if err := dc.LoadFontFace("Roboto-Medium.ttf", points); err != nil {
			html(c, http.StatusInternalServerError, renderError("Ocorreu um erro para encontrar a fonte"))
			return
		}

		text := c.DefaultQuery("t", fmt.Sprintf("%d x %d", width, height))

		dc.DrawStringWrapped(text, float64(width/2), float64(height/2), 0.5, 0.5, float64(width - 10), 1.3, gg.AlignCenter)
		data, err := imageToBytes(dc.Image())

		if err != nil {
			html(c, http.StatusInternalServerError, renderError("Ocorreu um erro para processar a imagem"))
			return
		}

		c.Data(http.StatusOK, "image/png", data.Bytes())
	})

	r.Run(getPort())
}

func imageToBytes(image image.Image) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, image)

	return buf, err;
}

func getPort() (string) {
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	return fmt.Sprintf(":%s", port)
}

func html(c *gin.Context, code int, text []byte) {
	c.Data(code, "text/html; charset=utf-8", text)
}
