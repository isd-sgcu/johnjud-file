package image

import (
	"github.com/isd-sgcu/johnjud-file/internal/service/image"
	imageRepo "github.com/isd-sgcu/johnjud-file/pkg/repository/image"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
)

func NewService(client image.Client, repo imageRepo.Repository) proto.ImageServiceServer {
	return image.NewService(client, repo)
}
