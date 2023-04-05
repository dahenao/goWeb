package products

import (
	"errors"

	"github.com/dahenao/goWeb/Package_Oriented_Design/internal/domain"
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
	data   []domain.Product
	lastId int
}

func (db *Local_slice_DB) GetAll() (result []domain.Product, err error) { //al estar declarado con nombre en variable de retorno toma valor por defecto

	if len(db.data) > 0 {
		result = db.data
	} else {
		err = ErrProductNotFound
	}
	return
}

func (db *Local_slice_DB) Create(pr *domain.Product) (err error) {

	pr.Id = db.lastId + 1
	db.data = append(db.data, *pr)

	db.lastId++
	return
}

func (db *Local_slice_DB) Update(index int, pr *domain.Product) (err error) {
	pr.Id = index //asignamos el id al producto completo a actualizar

	for i, prod := range db.data {
		if pr.Id == prod.Id { //si encuentra el producto con el id
			db.data[i] = *pr //actualiza en la posicion con el puntero de producto que viene desde el handler
			return           //sale del ciclo
		}
	}

	return ErrProductNotFound // si no encuentra producto retorna error
}

func (db *Local_slice_DB) GetProductByID(id int) (prod domain.Product, err error) {

	for _, p := range db.data {
		if id == p.Id {
			prod = p
			return
		}
	}
	err = ErrProductNotFound
	return
}

func (db *Local_slice_DB) Delete(id int) (err error) {

	for i, p := range db.data {
		if id == p.Id {
			db.data = append(db.data[:i], db.data[i+1:]...)
			return
		}
	}
	err = ErrProductNotFound
	return

}
