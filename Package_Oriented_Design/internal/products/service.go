package products

import (
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

func NewService(repo Repository) Service {
	return &ServiceDefault{BD: repo}
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
		switch err {
		case ErrProductNotFound:
			return prod, err
		default:
			return prod, ErrInternalServer
		}
	}
	return
}

func (s ServiceDefault) Update(index int, product *domain.Product) (err error) {

	err = s.BD.Update(index, product)
	if err != nil {
		switch err {
		case ErrProductNotFound:
			return err
		default:
			return ErrInternalServer
		}
	}
	return
}

func (s ServiceDefault) Delete(id int) (err error) {

	if err = s.BD.Delete(id); err != nil {
		switch err {
		case ErrProductNotFound:
			return err
		default:
			return ErrInternalServer
		}
	}

	return
}
