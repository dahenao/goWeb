package products

import (
	"errors"

	"github.com/dahenao/goWeb/Package_Oriented_Design/internal/domain"
)

type Service interface {
	Create(product *domain.Product) error
	GetAll() ([]domain.Product, error)
	Update(index int, product *domain.Product) error
	GetProductByID(id int) (domain.Product, error)
	Delete(id int) (err error)
}

type ServiceDefault struct {
	BD Repository //instanciar repositorio llamando la interface????
}

func (s ServiceDefault) Create(product *domain.Product) (err error) {

	err = s.BD.Create(product)
	if err != nil {

	}
	return
}

func (s ServiceDefault) GetAll() (allProducts []domain.Product, err error) {

	allProducts, err = s.BD.GetAll()
	if err != nil {

	}
	return
}

func (s ServiceDefault) GetProductByID(id int) (prod domain.Product, err error) {

	prod, err = s.BD.GetProductByID(id)
	if err != nil {

	}
	return
}

func (s ServiceDefault) Update(index int, product *domain.Product) (err error) {

	err = s.BD.Update(index, product)
	if err != nil {

	}
	return
}

func (s ServiceDefault) Delete(id int) (err error) {

	if err = s.BD.Delete(id); err != nil {
		if errors.Is(err, ErrProductNotFound) {
			return
		} else {
			err = ErrInternalServer
		}
	}
	return
}
