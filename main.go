package main

import (
	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"github.com/lucasb-eyer/go-colorful"
	"image/png"
	"net/http"
	"strconv"
	"regexp"
	"errors"
	"fmt"
	"math"
	"image"
	"bytes"
	"time"
	"flag"
)

const (
	maxImageSize = 5000
	minImageSize = 50
)

func main() {
	go clearOldCache()

	port := flag.Int("p", 8080, "Porta utilizada pelo servidor web")

	flag.Parse()

	fmt.Printf("Aplicação iniciada na porta %d\n", *port)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	mainEngine(r).Run(fmt.Sprintf(":%d", *port))
}

func mainEngine(r *gin.Engine) *gin.Engine {
	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.GET("/:size", func(c *gin.Context) {
		w, h, r, err := extractSize(c)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": err.Error()})
			return
		}

		bg, fg, err := extractColors(c)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": err.Error()})
			return
		}

		text := c.DefaultQuery("t", fmt.Sprintf("%d x %d", w, h))

		keyMap := makeKey(w, h, bg, fg, r, text)

		if img, ok := getCache(keyMap); ok {
			sendImage(c, img)
			return
		}

		img, err := generateImage(w, h, text, bg, fg, r)

		if err != nil {
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": err.Error()})
			return
		}

		sendImage(c, img)
	})

	return r
}

func sendImage(c *gin.Context, image *bytes.Buffer) {
	c.Header("Pragma", "public")
	c.Header("Cache-Control", "max-age=86400")
	c.Header("Expires", time.Now().AddDate(60, 0, 0).Format(http.TimeFormat))

	c.Data(http.StatusOK, "image/png", image.Bytes())
}

func extractSize(c *gin.Context) (int, int, float64, error) {
	r := regexp.MustCompile(`^(\d+)([xX](\d+))?$`)
	matches := r.FindStringSubmatch(c.Param("size"))

	if len(matches) == 0 {
		return 0, 0, 0, errors.New("Parâmetros inválidos")
	}

	width, _ := strconv.Atoi(matches[1])
	height, err := strconv.Atoi(matches[3])

	if err != nil {
		height = width
	}

	if width > maxImageSize || height > maxImageSize {
		return width, height, 0, errors.New("A imagem deve ter no máximo 5000 x 5000")
	}

	if width < minImageSize || height < minImageSize {
		return width, height, 0, errors.New("A imagem deve ter no mínimo 50 x 50")
	}

	round, err := strconv.ParseFloat(c.DefaultQuery("r", "0"), 64)

	if err != nil || round < 0 {
		return width, height, round, errors.New("O arredondamento deve ser maior ou igual a 0")
	}

	return width, height, round, nil;
}

func extractColors(c *gin.Context) (colorful.Color, colorful.Color, error) {
	bg, err := colorful.Hex("#" + c.DefaultQuery("bg", "666666"))
	if err != nil {
		return colorful.Color{}, colorful.Color{}, errors.New("Cor de fundo incorreta")
	}

	fg, err := colorful.Hex("#" + c.DefaultQuery("fg", "FFFFFF"))
	if err != nil {
		return colorful.Color{}, colorful.Color{}, errors.New("Cor do texto incorreta")
	}

	return bg, fg, nil
}

func generateImage(width, height int, text string, bg, fg colorful.Color, r float64) (*bytes.Buffer, error) {
	dc := gg.NewContext(width, height)
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), r)

	dc.SetColor(bg)
	dc.Fill()

	points := math.Min(float64(width), float64(height)) / 6

	dc.SetColor(fg)
	if err := dc.LoadFontFace("Roboto-Medium.ttf", points); err != nil {
		return nil, errors.New("Ocorreu um erro para encontrar a fonte")
	}

	dc.DrawStringWrapped(text, float64(width/2), float64(height/2), 0.5, 0.5, float64(width-10), 1.3, gg.AlignCenter)
	data, err := imageToBytes(dc.Image())

	if err != nil {
		return nil, errors.New("Ocorreu um erro para processar a imagem")
	}

	putCache(makeKey(width, height, bg, fg, r, text), data)

	return data, nil
}

func imageToBytes(image image.Image) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, image)

	return buf, err;
}

func makeKey(w, h int, bg, fg colorful.Color, r float64, text string) string {
	return fmt.Sprintf("%d;%d;%s;%s;%f;%s", w, h, bg.Hex(), fg.Hex(), r, text)
}
