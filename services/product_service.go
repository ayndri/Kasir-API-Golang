package services

import (
	"errors"
	"kasir-api/models"
	"kasir-api/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
	categoryRepo *repositories.CategoryRepository
}

func NewProductService(repo *repositories.ProductRepository, categoryRepo *repositories.CategoryRepository) *ProductService {
	return &ProductService{repo: repo, categoryRepo: categoryRepo}
}

func (s *ProductService) GetAll(name string) ([]models.Product, error) {
	return s.repo.GetAll(name)
}

func (s *ProductService) Create(data *models.Product) (*models.Product, error) {
    exists, err := s.categoryRepo.Exists(data.CategoryID)
    if err != nil {
        return nil, err
    }
    if !exists {
        return nil, errors.New("category not found")
    }
    
    err = s.repo.Create(data)
    if err != nil {
        return nil, err
    }

    return s.repo.GetByID(data.ID)
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Update(product *models.Product) (*models.Product, error) {
    exists, err := s.categoryRepo.Exists(product.CategoryID)
    if err != nil {
        return nil, err
    }
    if !exists {
        return nil, errors.New("category not found")
    }

    err = s.repo.Update(product)
    if err != nil {
        return nil, err
    }

    return s.repo.GetByID(product.ID)
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}