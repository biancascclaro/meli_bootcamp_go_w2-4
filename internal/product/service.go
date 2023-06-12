package product

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

type Service interface {
	Create(c context.Context, desc string, expR, freezeR int, height, length, netW float32, code string, freezeTemp, width float32, typeID, sellerID int) (domain.Product, error)
	GetAll(c context.Context) ([]domain.Product, error)
	Get(c context.Context, id int) (domain.Product, error)
	Delete(c context.Context, id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(c context.Context, desc string, expR, freezeR int, height, length, netW float32, code string, freezeTemp, width float32, typeID, sellerID int) (domain.Product, error) {
	ps, err := s.repo.GetAll(c)
	if err != nil {
		return domain.Product{}, NewErrGeneric("error fetching products")
	}

	if !isUniqueProductCode(code, ps) {
		return domain.Product{}, NewErrInvalidProductCode(code)
	}

	p := domain.Product{
		Description:    desc,
		ExpirationRate: expR,
		FreezingRate:   freezeR,
		Height:         height,
		Length:         length,
		Netweight:      netW,
		ProductCode:    code,
		RecomFreezTemp: freezeTemp,
		Width:          width,
		ProductTypeID:  typeID,
		SellerID:       sellerID,
	}
	id, err := s.repo.Save(c, p)
	if err != nil {
		return domain.Product{}, NewErrGeneric("error saving product")
	}

	p.ID = id
	return p, nil
}

func (s *service) GetAll(c context.Context) ([]domain.Product, error) {
	ps, err := s.repo.GetAll(c)
	if err != nil {
		return nil, NewErrGeneric("could not fetch products")
	}
	return ps, nil
}

func (s *service) Get(c context.Context, id int) (domain.Product, error) {
	p, err := s.repo.Get(c, id)
	if err != nil {
		// TODO: Properly handle DB communication error differently
		return domain.Product{}, NewErrNotFound(id)
	}
	return p, nil
}

func (s *service) Delete(c context.Context, id int) error {
	err := s.repo.Delete(c, id)
	if err != nil {
		if errors.Is(err, ErrNotFound{}) {
			return NewErrNotFound(id)
		} else {
			return NewErrGeneric("could not delete product")
		}
	}
	return nil
}

func isUniqueProductCode(code string, ps []domain.Product) bool {
	for _, p := range ps {
		if p.ProductCode == code {
			return false
		}
	}
	return true
}
