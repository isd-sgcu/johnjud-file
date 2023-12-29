package image

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/isd-sgcu/johnjud-file/internal/model"
	"github.com/isd-sgcu/johnjud-file/pkg/client/s3"
	"github.com/isd-sgcu/johnjud-file/pkg/repository/image"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type serviceImpl struct {
	proto.UnimplementedImageServiceServer
	client     s3.Client
	repository image.Repository
}

func NewService(client s3.Client, repository image.Repository) *serviceImpl {
	return &serviceImpl{
		client:     client,
		repository: repository,
	}
}

func (s *serviceImpl) FindByPetId(_ context.Context, req *proto.FindImageByPetIdRequest) (res *proto.FindImageByPetIdResponse, err error) {
	var images []*model.Image

	err = s.repository.FindByPetId(req.PetId, &images)
	if err != nil {
		log.Error().Err(err).
			Str("service", "pet").
			Str("module", "find by petId").
			Str("petId", req.PetId).
			Msg("Not found")

		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &proto.FindImageByPetIdResponse{Images: RawToDtoList(&images)}, nil
}

func (s *serviceImpl) Upload(_ context.Context, req *proto.UploadImageRequest) (res *proto.UploadImageResponse, err error) {
	// raw, _ := DtoToRaw(req)

	// err = s.repository.Create(raw)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, "failed to create like")
	// }

	// return &proto.UploadImageResponse{Image: RawToDto(raw)}, nil
	return nil, nil
}

func (s *serviceImpl) Delete(_ context.Context, req *proto.DeleteImageRequest) (res *proto.DeleteImageResponse, err error) {
	err = s.repository.Delete(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "something wrong when deleting like")
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
		return nil, err
	}

	return &model.Image{
		Base: model.Base{
			ID:        id,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		PetID:    &petId,
		ImageUrl: in.ImageUrl,
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
	return &proto.Image{
		Id:       in.ID.String(),
		PetId:    in.PetID.String(),
		ImageUrl: in.ImageUrl,
	}
}
