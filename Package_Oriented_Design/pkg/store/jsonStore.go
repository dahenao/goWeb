package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/dahenao/goWeb/Package_Oriented_Design/internal/domain"
)

var (
	ErrReadJson        = errors.New("Error: json file not found o invalid")
	ErrProductNotFound = errors.New("Error: Product not found")
)

type Store interface {
	//Create(product *domain.Product) error
	GetAll() ([]domain.Product, error)
	InicializeDatabase() ([]domain.Product, error)
	Update(index int, product domain.Product) error
	GetProductByID(id int) (domain.Product, error)
	saveProducts(products []domain.Product) error
	Delete(id int) (err error)
	Create(product *domain.Product) (err error)
}

type jsonStore struct {
	Path string
}

// NewJsonStore crea un nuevo store de products
func NewStore(path string) Store {
	return &jsonStore{
		Path: path,
	}
}

func (store *jsonStore) InicializeDatabase() (products []domain.Product, err error) {

	fileJSON, err := os.Open(store.Path)
	if err != nil {
		return
	}
	defer fileJSON.Close()

	err = json.NewDecoder(fileJSON).Decode(&products) //lee el archivo json y lo decodifica en el slice products
	if err != nil {
		return
	}

	return
}

func (store *jsonStore) GetAll() (products []domain.Product, err error) {
	products, err = store.InicializeDatabase()
	if err != nil {
		return nil, err
	}

	if len(products) > 0 {
		return
	} else {
		return nil, errors.New("no found data")
	}

}

func (store *jsonStore) GetProductByID(id int) (product domain.Product, err error) {
	var allProducts []domain.Product
	allProducts, err = store.InicializeDatabase()
	if err != nil {
		err = fmt.Errorf("Error: %w , %w", ErrReadJson, err)
		return
	}

	for _, prod := range allProducts {
		if prod.Id == id {
			product = prod
			return
		}
	}
	err = ErrProductNotFound
	return
}

// saveProducts guarda los productos en un archivo json
func (s *jsonStore) saveProducts(products []domain.Product) error {
	bytes, err := json.Marshal(products)
	if err != nil {
		return err
	}
	return os.WriteFile(s.Path, bytes, 0644)
}

func (store *jsonStore) Update(index int, product domain.Product) (err error) {
	var allProducts []domain.Product
	allProducts, err = store.InicializeDatabase()
	if err != nil {
		err = fmt.Errorf("Error: %w , %w", ErrReadJson, err)
		return
	}

	for i := range allProducts {
		if allProducts[i].Id == index {
			allProducts[i] = product
			store.saveProducts(allProducts)
			return
		}
	}
	err = ErrProductNotFound
	return
}

func (store *jsonStore) Delete(index int) (err error) {
	var allProducts []domain.Product
	allProducts, err = store.InicializeDatabase()
	if err != nil {
		err = fmt.Errorf("Error: %w , %w", ErrReadJson, err)
		return
	}

	for i := range allProducts {
		if allProducts[i].Id == index {
			allProducts = append(allProducts[:i], allProducts[i+1:]...)
			store.saveProducts(allProducts)
			return
		}
	}
	err = ErrProductNotFound
	return
}

func (store *jsonStore) Create(product *domain.Product) (err error) {
	var allProducts []domain.Product
	allProducts, err = store.InicializeDatabase()
	if err != nil {
		err = fmt.Errorf("Error: %w , %w", ErrReadJson, err)
		return
	}
	lastId := allProducts[len(allProducts)-1].Id
	product.Id = lastId + 1
	allProducts = append(allProducts, *product)

	return store.saveProducts(allProducts)
}
