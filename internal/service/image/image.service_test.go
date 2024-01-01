package image

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/isd-sgcu/johnjud-file/constant"
	"github.com/isd-sgcu/johnjud-file/internal/model"
	mock_bucket "github.com/isd-sgcu/johnjud-file/mocks/client/bucket"
	mock_image "github.com/isd-sgcu/johnjud-file/mocks/repository/image"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ImageServiceTest struct {
	suite.Suite
	petId      uuid.UUID
	findReq    *proto.FindImageByPetIdRequest
	uploadReq  *proto.UploadImageRequest
	assignReq  *proto.AssignPetRequest
	deleteReq  *proto.DeleteImageRequest
	imageProto *proto.Image
	image      *model.Image
	images     []*model.Image
}

func TestImageService(t *testing.T) {
	suite.Run(t, new(ImageServiceTest))
}

func (t *ImageServiceTest) SetupTest() {
	file := []byte("test")
	id := uuid.New()
	t.petId = uuid.New()
	objectKey := faker.Name()
	imageUrl := faker.URL()

	t.findReq = &proto.FindImageByPetIdRequest{
		PetId: t.petId.String(),
	}
	t.uploadReq = &proto.UploadImageRequest{
		Filename: objectKey,
		Data:     file,
		PetId:    t.petId.String(),
	}
	t.assignReq = &proto.AssignPetRequest{
		Ids:   []string{uuid.New().String(), uuid.New().String()},
		PetId: t.petId.String(),
	}
	t.deleteReq = &proto.DeleteImageRequest{
		Id:        id.String(),
		ObjectKey: objectKey,
	}
	t.imageProto = &proto.Image{
		Id:        id.String(),
		PetId:     t.petId.String(),
		ImageUrl:  imageUrl,
		ObjectKey: objectKey,
	}
	t.image = &model.Image{
		Base: model.Base{
			ID:        id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		PetID:     &t.petId,
		ImageUrl:  imageUrl,
		ObjectKey: objectKey,
	}
	t.images = []*model.Image{
		{
			Base: model.Base{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			PetID:     &t.petId,
			ImageUrl:  faker.URL(),
			ObjectKey: faker.Name(),
		},
		{
			Base: model.Base{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			PetID:     &t.petId,
			ImageUrl:  faker.URL(),
			ObjectKey: faker.Name(),
		},
	}
}

func (t *ImageServiceTest) TestFindByPetIdSuccess() {
	expected := &proto.FindImageByPetIdResponse{
		Images: []*proto.Image{
			{
				Id:        t.images[0].ID.String(),
				PetId:     t.images[0].PetID.String(),
				ImageUrl:  t.images[0].ImageUrl,
				ObjectKey: t.images[0].ObjectKey,
			},
			{
				Id:        t.images[1].ID.String(),
				PetId:     t.images[1].PetID.String(),
				ImageUrl:  t.images[1].ImageUrl,
				ObjectKey: t.images[1].ObjectKey,
			},
		},
	}
	var images []*model.Image

	controller := gomock.NewController(t.T())

	imageRepo := &mock_image.ImageRepositoryMock{}
	bucketClient := mock_bucket.NewMockClient(controller)
	imageRepo.On("FindByPetId", t.petId.String(), &images).Return(&t.images, nil)

	imageService := NewService(bucketClient, imageRepo)
	actual, err := imageService.FindByPetId(context.Background(), t.findReq)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *ImageServiceTest) TestFindByPetIdNotFound() {
	expected := status.Error(codes.NotFound, constant.ImageNotFoundErrorMessage)
	var images []*model.Image

	controller := gomock.NewController(t.T())

	imageRepo := &mock_image.ImageRepositoryMock{}
	bucketClient := mock_bucket.NewMockClient(controller)
	imageRepo.On("FindByPetId", t.petId.String(), &images).Return(nil, gorm.ErrRecordNotFound)

	imageService := NewService(bucketClient, imageRepo)
	actual, err := imageService.FindByPetId(context.Background(), t.findReq)

	status, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, status.Code())
	assert.Equal(t.T(), expected.Error(), err.Error())
}

func (t *ImageServiceTest) TestFindByPetIdInternalErr() {
	expected := status.Error(codes.Internal, constant.InternalServerErrorMessage)
	var images []*model.Image

	controller := gomock.NewController(t.T())

	imageRepo := &mock_image.ImageRepositoryMock{}
	bucketClient := mock_bucket.NewMockClient(controller)
	imageRepo.On("FindByPetId", t.petId.String(), &images).Return(nil, errors.New("Error finding image in db"))

	imageService := NewService(bucketClient, imageRepo)
	actual, err := imageService.FindByPetId(context.Background(), t.findReq)

	status, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Internal, status.Code())
	assert.Equal(t.T(), expected.Error(), err.Error())
}

func (t *ImageServiceTest) TestUploadSuccess() {
	expected := &proto.UploadImageResponse{
		Image: &proto.Image{
			Id:        t.imageProto.Id,
			PetId:     t.imageProto.PetId,
			ImageUrl:  t.imageProto.ImageUrl,
			ObjectKey: t.imageProto.ObjectKey,
		},
	}
	createImage := &model.Image{
		PetID:     t.image.PetID,
		ImageUrl:  t.image.ImageUrl,
		ObjectKey: t.image.ObjectKey,
	}

	controller := gomock.NewController(t.T())

	imageRepo := &mock_image.ImageRepositoryMock{}
	bucketClient := mock_bucket.NewMockClient(controller)
	imageRepo.On("Create", createImage).Return(t.image, nil)
	bucketClient.EXPECT().Upload(t.uploadReq.Data, t.uploadReq.Filename).Return(t.imageProto.ImageUrl, t.imageProto.ObjectKey, nil)

	imageService := NewService(bucketClient, imageRepo)
	actual, err := imageService.Upload(context.Background(), t.uploadReq)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *ImageServiceTest) TestUploadBucketFailed() {
	expected := status.Error(codes.Internal, constant.UploadToBucketErrorMessage)

	controller := gomock.NewController(t.T())

	imageRepo := &mock_image.ImageRepositoryMock{}
	bucketClient := mock_bucket.NewMockClient(controller)
	bucketClient.EXPECT().Upload(t.uploadReq.Data, t.uploadReq.Filename).Return("", "", errors.New("Error uploading to bucket client"))

	imageService := NewService(bucketClient, imageRepo)
	actual, err := imageService.Upload(context.Background(), t.uploadReq)

	status, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Internal, status.Code())
	assert.Equal(t.T(), expected.Error(), err.Error())
}

func (t *ImageServiceTest) TestUploadRepoFailed() {
	expected := status.Error(codes.Internal, constant.CreateImageErrorMessage)
	createImage := &model.Image{
		PetID:     t.image.PetID,
		ImageUrl:  t.image.ImageUrl,
		ObjectKey: t.image.ObjectKey,
	}

	controller := gomock.NewController(t.T())

	imageRepo := &mock_image.ImageRepositoryMock{}
	bucketClient := mock_bucket.NewMockClient(controller)
	imageRepo.On("Create", createImage).Return(nil, errors.New("Error creating image in db"))
	bucketClient.EXPECT().Upload(t.uploadReq.Data, t.uploadReq.Filename).Return(t.imageProto.ImageUrl, t.imageProto.ObjectKey, nil)

	imageService := NewService(bucketClient, imageRepo)
	actual, err := imageService.Upload(context.Background(), t.uploadReq)

	status, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Internal, status.Code())
	assert.Equal(t.T(), expected.Error(), err.Error())
}

func (t *ImageServiceTest) TestAssignPetSuccess() {
	expected := &proto.AssignPetResponse{
		Success: true,
	}
	id1, _ := uuid.Parse(t.assignReq.Ids[0])
	id2, _ := uuid.Parse(t.assignReq.Ids[1])
	petId, _ := uuid.Parse(t.assignReq.PetId)

	updateImages := []*model.Image{
		{
			PetID: &petId,
		},
		{
			PetID: &petId,
		},
	}

	image1 := model.Image{
		Base: model.Base{
			ID:        id1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		PetID:     &id1,
		ImageUrl:  faker.URL(),
		ObjectKey: faker.Name(),
	}
	image2 := model.Image{
		Base: model.Base{
			ID:        id2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		PetID:     &id1,
		ImageUrl:  faker.URL(),
		ObjectKey: faker.Name(),
	}

	controller := gomock.NewController(t.T())

	imageRepo := &mock_image.ImageRepositoryMock{}
	bucketClient := mock_bucket.NewMockClient(controller)
	imageRepo.On("Update", id1.String(), updateImages[0]).Return(&image1, nil)
	imageRepo.On("Update", id2.String(), updateImages[1]).Return(&image2, nil)

	imageService := NewService(bucketClient, imageRepo)
	actual, err := imageService.AssignPet(context.Background(), t.assignReq)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *ImageServiceTest) TestAssignPetNotFound() {
	expected := status.Error(codes.NotFound, constant.ImageNotFoundErrorMessage)

	id1, _ := uuid.Parse(t.assignReq.Ids[0])
	id2, _ := uuid.Parse(t.assignReq.Ids[1])
	petId, _ := uuid.Parse(t.assignReq.PetId)

	updateImages := []*model.Image{
		{
			PetID: &petId,
		},
		{
			PetID: &petId,
		},
	}

	image2 := model.Image{
		Base: model.Base{
			ID:        id2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		PetID:     &id1,
		ImageUrl:  faker.URL(),
		ObjectKey: faker.Name(),
	}

	controller := gomock.NewController(t.T())

	imageRepo := &mock_image.ImageRepositoryMock{}
	bucketClient := mock_bucket.NewMockClient(controller)
	imageRepo.On("Update", id1.String(), updateImages[0]).Return(nil, gorm.ErrRecordNotFound)
	imageRepo.On("Update", id2.String(), updateImages[1]).Return(&image2, nil)

	imageService := NewService(bucketClient, imageRepo)
	actual, err := imageService.AssignPet(context.Background(), t.assignReq)

	status, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, status.Code())
	assert.Equal(t.T(), expected.Error(), err.Error())
}

func (t *ImageServiceTest) TestAssignPetInternalErr() {
	expected := status.Error(codes.Internal, constant.InternalServerErrorMessage)

	id1, _ := uuid.Parse(t.assignReq.Ids[0])
	id2, _ := uuid.Parse(t.assignReq.Ids[1])
	petId, _ := uuid.Parse(t.assignReq.PetId)

	updateImages := []*model.Image{
		{
			PetID: &petId,
		},
		{
			PetID: &petId,
		},
	}

	image2 := model.Image{
		Base: model.Base{
			ID:        id2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		PetID:     &id1,
		ImageUrl:  faker.URL(),
		ObjectKey: faker.Name(),
	}

	controller := gomock.NewController(t.T())

	imageRepo := &mock_image.ImageRepositoryMock{}
	bucketClient := mock_bucket.NewMockClient(controller)
	imageRepo.On("Update", id1.String(), updateImages[0]).Return(nil, errors.New("Error updating image in db"))
	imageRepo.On("Update", id2.String(), updateImages[1]).Return(&image2, nil)

	imageService := NewService(bucketClient, imageRepo)
	actual, err := imageService.AssignPet(context.Background(), t.assignReq)

	status, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Internal, status.Code())
	assert.Equal(t.T(), expected.Error(), err.Error())
}
