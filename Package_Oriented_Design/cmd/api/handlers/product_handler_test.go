package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dahenao/goWeb/Package_Oriented_Design/internal/domain"
	"github.com/dahenao/goWeb/Package_Oriented_Design/internal/products"
	"github.com/dahenao/goWeb/Package_Oriented_Design/pkg/store"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type Resp struct {
	Data domain.Product `json:"data"`
}

type respList struct {
	Data []domain.Product `json:"data"`
}

func createServerForProductFunctionalTest() *gin.Engine { //returnamos un servidor de web gin
	godotenv.Load("../../../cmd/api/.env")
	gin.SetMode(gin.TestMode) //seteamos testmode para evitar logs innecesarios
	//server := gin.New() crea un servidor gin vacio sin middlewarez por defecto o rutas
	store := store.NewStore("../../../../products.json")
	repository := products.NewRepository(store)
	service := products.NewService(repository)
	handler := ProductHandler{
		Service: service,
	}

	router := gin.New() //gin.Default() //crear un router de gin

	routerGroup := router.Group("/products") //grear un agrupador de rutas que comparten middlewarez y cosas en comun
	{
		routerGroup.Use(TokenMiddlewareValidate())
		routerGroup.GET("", handler.GetAll()) //crear metodo get
		routerGroup.GET(":id", handler.getProductByID())
		routerGroup.POST("", handler.Create())
		routerGroup.DELETE(":id", handler.Delete())
		routerGroup.PUT(":id", handler.Update())
		routerGroup.PATCH(":id", handler.UpdatePartial())
	}
	return router //retorna el router ya configurado como un servidor de gin

}

func createRequestResponse(method, url, body string) (*http.Request, *httptest.ResponseRecorder) {

	request := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body))) //
	response := httptest.NewRecorder()
	request.Header.Set("token", "909090")
	if body != "" {
		request.Header.Set("Content-Type", "application/json")

	}

	return request, response
}
func createDatabase(path string) []domain.Product {
	fileJSON, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fileJSON.Close()
	var Products []domain.Product
	err = json.NewDecoder(fileJSON).Decode(&Products) //lee el archivo json y lo decodifica en el slice products
	if err != nil {
		panic(err)
	}
	return Products
}
func TestProductGetByID(t *testing.T) {

	t.Run("test response with one product", func(t *testing.T) {

		//implement AAA test method
		//Arrange
		ExpectedBody := domain.Product{Id: 503,
			Name:         "Product",
			Quantity:     8,
			Code_value:   "000990",
			Is_published: true,
			Expiration:   "23/04/2044",
			Price:        0}
		ExpectedStatusCode := http.StatusOK

		//Act

		//creamos servidor
		server := createServerForProductFunctionalTest()

		req, res := createRequestResponse(http.MethodGet, "/products/503", "")

		server.ServeHTTP(res, req)
		fmt.Println(res.Body)
		bodyResponse := Resp{ //creamos el tipo de dato donde vamos a parsear el json de respuesta
			Data: domain.Product{},
		}

		json.Unmarshal(res.Body.Bytes(), &bodyResponse)

		//Assert
		assert.Equal(t, ExpectedStatusCode, res.Code)
		assert.Equal(t, ExpectedBody, bodyResponse.Data)

	})

}
func TestProductsGetAll(t *testing.T) {

	t.Run("test response with all products", func(t *testing.T) {

		//implement AAA test method
		//Arrange
		ExpectedBody := createDatabase("../../../../products.json")
		ExpectedStatusCode := http.StatusOK
		//ExpectedQuantityProds := 503
		//Act

		//creamos servidor
		server := createServerForProductFunctionalTest()

		req, res := createRequestResponse(http.MethodGet, "/products", "")

		server.ServeHTTP(res, req)
		//fmt.Println(res.Body)
		bodyResponse := respList{ //creamos el tipo de dato donde vamos a parsear el json de respuesta
			Data: []domain.Product{},
		}

		json.Unmarshal(res.Body.Bytes(), &bodyResponse)

		//Assert
		assert.Equal(t, ExpectedStatusCode, res.Code)
		//assert.Equal(t, ExpectedQuantityProds, len(bodyResponse.Data))
		assert.Equal(t, ExpectedBody, bodyResponse.Data)

	})

}

func TestCreateProduct(t *testing.T) {

	t.Run("test response create an product", func(t *testing.T) {

		//implement AAA test method
		//Arrange
		ExpectedBody := domain.Product{
			Id:           504,
			Name:         "ProductTest",
			Quantity:     888,
			Code_value:   "000990",
			Is_published: true,
			Expiration:   "23/04/2044",
			Price:        230.0}
		ExpectedStatusCode := http.StatusCreated
		//Act

		productToCreate := domain.Product{
			Name:         "ProductTest",
			Quantity:     888,
			Code_value:   "000990",
			Is_published: true,
			Expiration:   "23/04/2044",
			Price:        230.0}

		sjson, _ := json.Marshal(productToCreate) //convierte el objeto a bytes

		//creamos servidor
		server := createServerForProductFunctionalTest()

		req, res := createRequestResponse(http.MethodPost, "/products", string(sjson))

		server.ServeHTTP(res, req)
		//fmt.Println(res.Body)
		bodyResponse := Resp{ //creamos el tipo de dato donde vamos a parsear el json de respuesta
			Data: domain.Product{},
		}

		json.Unmarshal(res.Body.Bytes(), &bodyResponse)

		//Assert
		assert.Equal(t, ExpectedStatusCode, res.Code)
		assert.Equal(t, ExpectedBody, bodyResponse.Data)

	})

}

func TestProductDelete(t *testing.T) {

	t.Run("test response delete product", func(t *testing.T) {

		//implement AAA test method
		//Arrange

		ExpectedStatusCode := http.StatusNoContent

		//Act

		//creamos servidor
		server := createServerForProductFunctionalTest()

		req, res := createRequestResponse(http.MethodDelete, "/products/504", "")

		server.ServeHTTP(res, req)
		fmt.Println(res.Body)
		bodyResponse := Resp{ //creamos el tipo de dato donde vamos a parsear el json de respuesta
			Data: domain.Product{},
		}

		json.Unmarshal(res.Body.Bytes(), &bodyResponse)

		//Assert
		assert.Equal(t, ExpectedStatusCode, res.Code)

	})

}

func TestProductBadRequest(t *testing.T) {

	t.Run("test response delete product bad id", func(t *testing.T) {

		//implement AAA test method
		//Arrange
		ExpectedStatusCode := http.StatusBadRequest

		//Act
		server := createServerForProductFunctionalTest()

		req, res := createRequestResponse(http.MethodDelete, "/products/k", "")

		server.ServeHTTP(res, req)

		//Assert
		assert.Equal(t, ExpectedStatusCode, res.Code)

	})

	t.Run("test response update product bad id", func(t *testing.T) {

		//implement AAA test method
		//Arrange
		ExpectedStatusCode := http.StatusBadRequest

		//Act
		server := createServerForProductFunctionalTest()

		req, res := createRequestResponse(http.MethodPut, "/products/k", "")

		server.ServeHTTP(res, req)

		//Assert
		assert.Equal(t, ExpectedStatusCode, res.Code)

	})

	t.Run("test response partial update product bad id", func(t *testing.T) {

		//implement AAA test method
		//Arrange
		ExpectedStatusCode := http.StatusBadRequest

		//Act
		server := createServerForProductFunctionalTest()

		req, res := createRequestResponse(http.MethodPatch, "/products/k", "")

		server.ServeHTTP(res, req)

		//Assert
		assert.Equal(t, ExpectedStatusCode, res.Code)

	})

	t.Run("test response get by id product bad id", func(t *testing.T) {

		//implement AAA test method
		//Arrange
		ExpectedStatusCode := http.StatusBadRequest

		//Act
		server := createServerForProductFunctionalTest()

		req, res := createRequestResponse(http.MethodGet, "/products/k", "")

		server.ServeHTTP(res, req)

		//Assert
		assert.Equal(t, ExpectedStatusCode, res.Code)

	})

}

func TestProductNotFound(t *testing.T) {

	t.Run("test response error getbyid when product no exist", func(t *testing.T) {

		//implement AAA test method
		//Arrange

		ExpectedStatusCode := http.StatusNotFound

		//Act

		//creamos servidor
		server := createServerForProductFunctionalTest()

		req, res := createRequestResponse(http.MethodGet, "/products/888", "")

		server.ServeHTTP(res, req)
		fmt.Println(res.Body)

		//Assert
		assert.Equal(t, ExpectedStatusCode, res.Code)

	})
	t.Run("test response error update partial when product no exist", func(t *testing.T) {

		//implement AAA test method
		//Arrange

		ExpectedStatusCode := http.StatusNotFound

		//Act

		//creamos servidor
		server := createServerForProductFunctionalTest()

		req, res := createRequestResponse(http.MethodPatch, "/products/888", "")

		server.ServeHTTP(res, req)
		fmt.Println(res.Body)

		//Assert
		assert.Equal(t, ExpectedStatusCode, res.Code)

	})

	t.Run("test response error Update when product no exist", func(t *testing.T) {

		//implement AAA test method
		//Arrange

		ExpectedStatusCode := http.StatusNotFound

		//Act
		productToCreate := domain.Product{
			Name:         "ProductTest",
			Quantity:     888,
			Code_value:   "000990",
			Is_published: true,
			Expiration:   "23/04/2044",
			Price:        230.0}

		sjson, _ := json.Marshal(productToCreate) //convierte el objeto a bytes

		//creamos servidor
		server := createServerForProductFunctionalTest()

		req, res := createRequestResponse(http.MethodPut, "/products/888", string(sjson))

		server.ServeHTTP(res, req)
		fmt.Println("update")
		fmt.Println(res.Body)

		//Assert
		assert.Equal(t, ExpectedStatusCode, res.Code)

	})

	t.Run("test response error delete when product no exist", func(t *testing.T) {

		//implement AAA test method
		//Arrange

		ExpectedStatusCode := http.StatusNotFound

		//Act

		//creamos servidor
		server := createServerForProductFunctionalTest()

		req, res := createRequestResponse(http.MethodDelete, "/products/888", "")

		server.ServeHTTP(res, req)
		fmt.Println(res.Body)

		//Assert
		assert.Equal(t, ExpectedStatusCode, res.Code)

	})

}

func TestProductUnAuthorized(t *testing.T) {

	t.Run("test response error getbyid when invalid token", func(t *testing.T) {

		//implement AAA test method
		//Arrange

		ExpectedStatusCode := http.StatusUnauthorized

		//Act

		//creamos servidor
		server := createServerForProductFunctionalTest()

		req, res := createRequestResponse(http.MethodGet, "/products/888", "")
		req.Header.Del("token")
		server.ServeHTTP(res, req)

		//Assert
		assert.Equal(t, ExpectedStatusCode, res.Code)

	})
	t.Run("test response error update partial  when invalid token", func(t *testing.T) {

		//implement AAA test method
		//Arrange

		ExpectedStatusCode := http.StatusUnauthorized

		//Act

		//creamos servidor
		server := createServerForProductFunctionalTest()

		req, res := createRequestResponse(http.MethodPatch, "/products/888", "")
		req.Header.Del("token")
		server.ServeHTTP(res, req)

		//Assert
		assert.Equal(t, ExpectedStatusCode, res.Code)

	})

	t.Run("test response error Update when invalid token", func(t *testing.T) {

		//implement AAA test method
		//Arrange

		ExpectedStatusCode := http.StatusUnauthorized

		//Act
		productToCreate := domain.Product{
			Name:         "ProductTest",
			Quantity:     888,
			Code_value:   "000990",
			Is_published: true,
			Expiration:   "23/04/2044",
			Price:        230.0}

		sjson, _ := json.Marshal(productToCreate) //convierte el objeto a bytes

		//Act

		//creamos servidor
		server := createServerForProductFunctionalTest()

		req, res := createRequestResponse(http.MethodPut, "/products/888", string(sjson))
		req.Header.Del("token")
		server.ServeHTTP(res, req)

		//Assert
		assert.Equal(t, ExpectedStatusCode, res.Code)

	})

	t.Run("test response error delete when invalid tkn", func(t *testing.T) {

		//implement AAA test method
		//Arrange

		ExpectedStatusCode := http.StatusUnauthorized

		//Act

		//creamos servidor
		server := createServerForProductFunctionalTest()

		req, res := createRequestResponse(http.MethodDelete, "/products/888", "")
		req.Header.Del("token")
		server.ServeHTTP(res, req)

		//Assert
		assert.Equal(t, ExpectedStatusCode, res.Code)

	})

}
