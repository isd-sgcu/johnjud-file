package image

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/isd-sgcu/johnjud-file/constant"
	"github.com/isd-sgcu/johnjud-file/internal/model"
	"github.com/isd-sgcu/johnjud-file/internal/utils"
	"github.com/isd-sgcu/johnjud-file/pkg/client/bucket"
	"github.com/isd-sgcu/johnjud-file/pkg/repository/image"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type serviceImpl struct {
	proto.UnimplementedImageServiceServer
	client     bucket.Client
	repository image.Repository
	random     utils.RandomUtil
}

func NewService(client bucket.Client, repository image.Repository, random utils.RandomUtil) proto.ImageServiceServer {
	return &serviceImpl{
		client:     client,
		repository: repository,
		random:     random,
	}
}

func (s *serviceImpl) FindByPetId(_ context.Context, req *proto.FindImageByPetIdRequest) (res *proto.FindImageByPetIdResponse, err error) {
	var images []*model.Image

	err = s.repository.FindByPetId(req.PetId, &images)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Error().Err(err).
				Str("service", "image").
				Str("module", "find by petId").
				Str("petId", req.PetId).
				Msg(constant.ImageNotFoundErrorMessage)

			return nil, status.Error(codes.NotFound, constant.ImageNotFoundErrorMessage)
		}

		log.Error().Err(err).
			Str("service", "image").
			Str("module", "find by petId").
			Str("petId", req.PetId).
			Msg(constant.InternalServerErrorMessage)

		return nil, status.Error(codes.Internal, constant.InternalServerErrorMessage)
	}

	return &proto.FindImageByPetIdResponse{Images: RawToDtoList(&images)}, nil
}

func (s *serviceImpl) Upload(_ context.Context, req *proto.UploadImageRequest) (res *proto.UploadImageResponse, err error) {
	randomString, err := s.random.GenerateRandomString(10)
	if err != nil {
		log.Error().Err(err).
			Str("service", "image").
			Str("module", "upload").
			Str("petId", req.PetId).
			Msg("Error while generating random string")
		return nil, status.Error(codes.Internal, "Error while generating random string")
	}
	imageUrl, objectKey, err := s.client.Upload(req.Data, req.Filename+"_"+randomString)
	if err != nil {
		log.Error().Err(err).
			Str("service", "image").
			Str("module", "upload").
			Str("petId", req.PetId).
			Msg(constant.UploadToBucketErrorMessage)

		return nil, status.Error(codes.Internal, constant.UploadToBucketErrorMessage)
	}

	raw, _ := DtoToRaw(&proto.Image{
		PetId:     req.PetId,
		ImageUrl:  imageUrl,
		ObjectKey: objectKey,
	})

	err = s.repository.Create(raw)
	if err != nil {
		log.Error().Err(err).
			Str("service", "image").
			Str("module", "upload").
			Str("petId", req.PetId).
			Msg(constant.CreateImageErrorMessage)

		return nil, status.Error(codes.Internal, constant.CreateImageErrorMessage)
	}

	return &proto.UploadImageResponse{Image: RawToDto(raw)}, nil
}

func (s *serviceImpl) AssignPet(_ context.Context, req *proto.AssignPetRequest) (res *proto.AssignPetResponse, err error) {
	petId, err := uuid.Parse(req.PetId)
	if err != nil {
		log.Error().Err(err).
			Str("service", "image").
			Str("module", "assign pet").
			Str("petId", req.PetId).
			Msg(constant.PrimaryKeyRequiredErrorMessage)

		return nil, status.Error(codes.InvalidArgument, constant.PrimaryKeyRequiredErrorMessage)
	}

	for _, id := range req.Ids {
		err = s.repository.Update(id, &model.Image{
			PetID: &petId,
		})
		switch err {
		case nil:
			continue
		case gorm.ErrRecordNotFound:
			log.Error().Err(err).
				Str("service", "image").
				Str("module", "assign pet").
				Str("petId", req.PetId).
				Msg(constant.ImageNotFoundErrorMessage)

			return nil, status.Error(codes.NotFound, constant.ImageNotFoundErrorMessage)
		default:
			log.Error().Err(err).
				Str("service", "image").
				Str("module", "assign pet").
				Str("petId", req.PetId).
				Msg(constant.InternalServerErrorMessage)

			return nil, status.Error(codes.Internal, constant.InternalServerErrorMessage)
		}
	}

	return &proto.AssignPetResponse{Success: true}, nil
}

func (s *serviceImpl) Delete(_ context.Context, req *proto.DeleteImageRequest) (res *proto.DeleteImageResponse, err error) {
	err = s.client.Delete(req.ObjectKey)
	if err != nil {
		log.Error().Err(err).
			Str("service", "image").
			Str("module", "delete").
			Str("id", req.Id).
			Msg(constant.DeleteFromBucketErrorMessage)

		return nil, status.Error(codes.Internal, constant.DeleteFromBucketErrorMessage)
	}

	err = s.repository.Delete(req.Id)
	if err != nil {
		log.Error().Err(err).
			Str("service", "image").
			Str("module", "delete").
			Str("id", req.Id).
			Msg(constant.DeleteImageErrorMessage)

		return nil, status.Error(codes.Internal, constant.DeleteImageErrorMessage)
	}

	return &proto.DeleteImageResponse{Success: true}, nil
}

func DtoToRaw(in *proto.Image) (result *model.Image, err error) {
	var id uuid.UUID
	if in.Id != "" {
		id, err = uuid.Parse(in.Id)
		if err != nil {
			return nil, err
		}
	}

	petId, err := uuid.Parse(in.PetId)
	if err != nil {
		return &model.Image{
			Base: model.Base{
				ID:        id,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				DeletedAt: gorm.DeletedAt{},
			},
			PetID:     nil,
			ImageUrl:  in.ImageUrl,
			ObjectKey: in.ObjectKey,
		}, nil
	}

	return &model.Image{
		Base: model.Base{
			ID:        id,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		PetID:     &petId,
		ImageUrl:  in.ImageUrl,
		ObjectKey: in.ObjectKey,
	}, nil
}

func RawToDtoList(in *[]*model.Image) []*proto.Image {
	var result []*proto.Image
	for _, b := range *in {
		result = append(result, RawToDto(b))
	}

	return result
}

func RawToDto(in *model.Image) *proto.Image {
	var id string
	var petId string
	if in.ID != uuid.Nil {
		id = in.ID.String()
	}
	if in.PetID != nil {
		petId = in.PetID.String()
	}

	return &proto.Image{
		Id:        id,
		PetId:     petId,
		ImageUrl:  in.ImageUrl,
		ObjectKey: in.ObjectKey,
	}
}
