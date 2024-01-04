package image

import (
	"github.com/isd-sgcu/johnjud-file/internal/model"
	"github.com/stretchr/testify/mock"
)

type ImageRepositoryMock struct {
	mock.Mock
}

func (m *ImageRepositoryMock) FindOne(id string, image *model.Image) error {
	args := m.Called(id, image)
	if args.Get(0) != nil {
		*image = *args.Get(0).(*model.Image)
		return nil
	}

	return args.Error(1)
}

func (m *ImageRepositoryMock) FindByPetId(id string, image *[]*model.Image) error {
	args := m.Called(id, image)
	if args.Get(0) != nil {
		*image = *args.Get(0).(*[]*model.Image)
		return nil
	}

	return args.Error(1)
}

func (m *ImageRepositoryMock) Create(image *model.Image) error {
	args := m.Called(image)
	if args.Get(0) != nil {
		*image = *args.Get(0).(*model.Image)
		return nil
	}

	return args.Error(1)
}

func (m *ImageRepositoryMock) Update(id string, image *model.Image) error {
	args := m.Called(id, image)
	if args.Get(0) != nil {
		*image = *args.Get(0).(*model.Image)
		return nil
	}

	return args.Error(1)
}

func (m *ImageRepositoryMock) Delete(id string) error {
	args := m.Called(id)

	return args.Error(0)
}
