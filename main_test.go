package main

import (
	"net/http/httptest"
	"net/http"
	"testing"
	"io/ioutil"
	"strings"
	"image/png"
	"github.com/gin-gonic/gin"
	"fmt"
)

func getRequest(url string) (*http.Response, error) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	ts := httptest.NewServer(mainEngine(r))
	defer ts.Close()

	c := ts.Client()

	res, err := c.Get(ts.URL + url)

	if err != nil {
		return nil, fmt.Errorf("Não foi possível solicitar a url / (%v)", err)
	}

	return res, nil
}

// TestMainPage - Testa se a tela inicial/ajuda está aparecendo corretamente
func TestMainPage(t *testing.T) {
	res, err := getRequest("/")

	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Fatalf("Não foi possível ler o corpo da solicitação / (%v)", err)
	}

	text := "<h1>Image Mock</h1>"
	if !strings.Contains(string(body), text) {
		t.Errorf("\"%s\" não encontrada no retorno da requisição. (Obtido %s)", text, string(body))
	}
}

// TestCreateBasicImage - Testa se o app cria com suceso uma imagem png
func TestCreateBasicImage(t *testing.T) {
	res, err := getRequest("/300")

	if err != nil {
		t.Fatal(err)
	}

	_, err =  png.Decode(res.Body)
	res.Body.Close()

	if err != nil {
		t.Fatalf("Não foi gerado um png válido. Motivo: %v", err)
	}
}

type TestCase struct {
	url string
	message string
}

// TestParamsConfigValidation - Testa se o app valida entrada de parametros incorretas
func TestColorsConfigValidation(t *testing.T) {
	testCases := []TestCase{
		{url: "/wrong", message: "Parâmetros inválidos"},
		{url: "/300?bg=wrong", message: "Cor de fundo incorreta"},
		{url: "/300?fg=wrong", message: "Cor do texto incorreta"},
		{url: "/300?r=x", message: "O arredondamento deve ser maior ou igual a 0"},
		{url: "/1", message: "A imagem deve ter no mínimo 50 x 50"},
		{url: "/15000", message: "A imagem deve ter no máximo 5000 x 5000"},
	}

	for _, test := range testCases {
		res, err := getRequest(test.url)

		if err != nil {
			t.Fatal(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()

		if !strings.Contains(string(body), test.message) {
			t.Errorf("\"%s\" não encontrada no retorno da requisição. (Obtido %s)", test.message, string(body))
		}
	}
}

// TestParamsConfigValidation - Verifica se a aplicação utiliza o cache de imagens identicas
func TestCacheImageRequest(t *testing.T) {
	_, err := getRequest("/300")

	if err != nil {
		t.Fatal(err)
	}

	_, err = getRequest("/300")

	if err != nil {
		t.Fatal(err)
	}

	if len(cache) == 0 {
		t.Errorf("Não foi atualizado o cache")
	}
}
