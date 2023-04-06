package products

import (
	"errors"
	"fmt"

	"github.com/dahenao/goWeb/Package_Oriented_Design/internal/domain"
	"github.com/dahenao/goWeb/Package_Oriented_Design/pkg/store"
)

var (
	ErrProductAlreadyExists = errors.New("product already exist")
	ErrProductNotFound      = errors.New("Product not found")
	ErrInternalServer       = errors.New("internal error")
)

type Repository interface {
	Create(product *domain.Product) error
	GetAll() ([]domain.Product, error)
	Update(index int, product *domain.Product) error
	GetProductByID(id int) (domain.Product, error)
	Delete(id int) (err error)
}

type Local_slice_DB struct {
	Storage store.Store
}

func (db *Local_slice_DB) GetAll() (result []domain.Product, err error) { //al estar declarado con nombre en variable de retorno toma valor por defecto
	result, err = db.Storage.GetAll()
	fmt.Println(len(result))
	if len(result) > 0 {
		return
	} else {
		err = ErrProductNotFound
	}
	return
}

func (db *Local_slice_DB) Create(pr *domain.Product) (err error) {

	err = db.Storage.Create(pr) //como en el storage tambien recibe el puntero no se necesita desreferenciar como en el caso del uodate

	return
}

func (db *Local_slice_DB) Update(index int, pr *domain.Product) (err error) {
	pr.Id = index //asignamos el id al producto completo a actualizar
	return db.Storage.Update(index, *pr)
}
func (db *Local_slice_DB) GetProductByID(id int) (prod domain.Product, err error) {

	return db.Storage.GetProductByID(id)
}

func (db *Local_slice_DB) Delete(id int) (err error) {

	return db.Storage.Delete(id)

}
