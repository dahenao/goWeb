package main

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProductByID(t *testing.T) {

	pExpected := Product{
		Id:           9,
		Name:         "Apple - Delicious, Golden",
		Quantity:     225,
		Code_value:   "S73046D",
		Is_published: true,
		Expiration:   "02/04/2021",
		Price:        976.27}

	createDatabase("products.json")
	pResult, err := getProductByID(9)
	assert.Equal(t, nil, err, "el error debe ser nil")
	assert.Equal(t, pExpected, pResult, "los productos no coinciden")

}

func TestPingEndpoint(t *testing.T) {
	url := "http://localhost:8080/ping"

	// Hacer la petición GET
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Leer la respuesta JSON

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Utiliza Testify para realizar afirmaciones sobre la respuesta HTTP
	assert.Equal(t, http.StatusOK, resp.StatusCode, "código de estado incorrecto") //validamos codigo de respuesta
	assert.Contains(t, string(body), `pong`, "respuesta incorrecta")               //validamos respuesta
}
