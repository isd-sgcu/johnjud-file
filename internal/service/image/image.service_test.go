package image

import (
	"context"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/isd-sgcu/johnjud-file/internal/model"
	mock_bucket "github.com/isd-sgcu/johnjud-file/mocks/client/bucket"
	mock_image "github.com/isd-sgcu/johnjud-file/mocks/repository/image"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ImageServiceTest struct {
	suite.Suite
	uploadReq  *proto.UploadImageRequest
	deleteReq  *proto.DeleteImageRequest
	imageProto *proto.Image
	image      *model.Image
}

func TestImageService(t *testing.T) {
	suite.Run(t, new(ImageServiceTest))
}

func (t *ImageServiceTest) SetupTest() {
	file := []byte("test")
	id := uuid.New()
	petId := uuid.New()
	objectKey := faker.Name()
	imageUrl := faker.URL()

	t.uploadReq = &proto.UploadImageRequest{
		Filename: objectKey,
		Data:     file,
		PetId:    petId.String(),
	}
	t.deleteReq = &proto.DeleteImageRequest{
		Id:        id.String(),
		ObjectKey: objectKey,
	}
	t.imageProto = &proto.Image{
		Id:        id.String(),
		PetId:     petId.String(),
		ImageUrl:  imageUrl,
		ObjectKey: objectKey,
	}
	t.image = &model.Image{
		Base: model.Base{
			ID:        id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		PetID:     &petId,
		ImageUrl:  imageUrl,
		ObjectKey: objectKey,
	}
}

func (t *ImageServiceTest) TestUploadSuccess() {
	var uuid uuid.UUID
	expected := &proto.UploadImageResponse{
		Image: &proto.Image{
			Id:        uuid.String(),
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

	imageRepo := mock_image.NewMockRepository(controller)
	imageRepo.EXPECT().Create(createImage).Return(nil)
	bucketClient := mock_bucket.NewMockClient(controller)
	bucketClient.EXPECT().Upload(t.uploadReq.Data, t.uploadReq.Filename).Return(t.imageProto.ImageUrl, t.imageProto.ObjectKey, nil)

	imageService := NewService(bucketClient, imageRepo)
	actual, err := imageService.Upload(context.Background(), t.uploadReq)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

// func (t *AuthServiceTest) TestSignupHashPasswordFailed() {
// 	hashPasswordErr := errors.New("Hash password error")

// 	expected := status.Error(codes.Internal, constant.InternalServerErrorMessage)

// 	controller := gomock.NewController(t.T())

// 	authRepo := mock_auth.NewMockRepository(controller)
// 	userRepo := user.UserRepositoryMock{}
// 	tokenService := token.TokenServiceMock{}
// 	bcryptUtil := utils.BcryptUtilMock{}

// 	bcryptUtil.On("GenerateHashedPassword", t.signupRequest.Password).Return("", hashPasswordErr)

// 	authSvc := NewService(authRepo, &userRepo, &tokenService, &bcryptUtil)
// 	actual, err := authSvc.SignUp(t.ctx, t.signupRequest)

// 	status, ok := status.FromError(err)

// 	assert.Nil(t.T(), actual)
// 	assert.Equal(t.T(), codes.Internal, status.Code())
// 	assert.True(t.T(), ok)
// 	assert.Equal(t.T(), expected.Error(), err.Error())
// }

// func (t *AuthServiceTest) TestSignupCreateUserDuplicateConstraint() {
// 	hashedPassword := faker.Password()
// 	newUser := &model.User{
// 		Email:     t.signupRequest.Email,
// 		Password:  hashedPassword,
// 		Firstname: t.signupRequest.FirstName,
// 		Lastname:  t.signupRequest.LastName,
// 		Role:      constant.USER,
// 	}
// 	createUserErr := gorm.ErrDuplicatedKey

// 	expected := status.Error(codes.AlreadyExists, constant.DuplicateEmailErrorMessage)

// 	controller := gomock.NewController(t.T())

// 	authRepo := mock_auth.NewMockRepository(controller)
// 	userRepo := user.UserRepositoryMock{}
// 	tokenService := token.TokenServiceMock{}
// 	bcryptUtil := utils.BcryptUtilMock{}

// 	bcryptUtil.On("GenerateHashedPassword", t.signupRequest.Password).Return(hashedPassword, nil)
// 	userRepo.On("Create", newUser).Return(nil, createUserErr)

// 	authSvc := NewService(authRepo, &userRepo, &tokenService, &bcryptUtil)
// 	actual, err := authSvc.SignUp(t.ctx, t.signupRequest)

// 	status, ok := status.FromError(err)

// 	assert.Nil(t.T(), actual)
// 	assert.Equal(t.T(), codes.AlreadyExists, status.Code())
// 	assert.True(t.T(), ok)
// 	assert.Equal(t.T(), expected.Error(), err.Error())
// }

// func (t *AuthServiceTest) TestSignupCreateUserInternalFailed() {
// 	hashedPassword := faker.Password()
// 	newUser := &model.User{
// 		Email:     t.signupRequest.Email,
// 		Password:  hashedPassword,
// 		Firstname: t.signupRequest.FirstName,
// 		Lastname:  t.signupRequest.LastName,
// 		Role:      constant.USER,
// 	}
// 	createUserErr := errors.New("Internal server error")

// 	expected := status.Error(codes.Internal, constant.InternalServerErrorMessage)

// 	controller := gomock.NewController(t.T())

// 	authRepo := mock_auth.NewMockRepository(controller)
// 	userRepo := user.UserRepositoryMock{}
// 	tokenService := token.TokenServiceMock{}
// 	bcryptUtil := utils.BcryptUtilMock{}

// 	bcryptUtil.On("GenerateHashedPassword", t.signupRequest.Password).Return(hashedPassword, nil)
// 	userRepo.On("Create", newUser).Return(nil, createUserErr)

// 	authSvc := NewService(authRepo, &userRepo, &tokenService, &bcryptUtil)
// 	actual, err := authSvc.SignUp(t.ctx, t.signupRequest)

// 	status, ok := status.FromError(err)

// 	assert.Nil(t.T(), actual)
// 	assert.Equal(t.T(), codes.Internal, status.Code())
// 	assert.True(t.T(), ok)
// 	assert.Equal(t.T(), expected.Error(), err.Error())
// }

// func (t *AuthServiceTest) TestSignInSuccess() {
// 	existUser := &model.User{
// 		Base: model.Base{
// 			ID: uuid.New(),
// 		},
// 		Email:     t.signInRequest.Email,
// 		Password:  faker.Password(),
// 		Firstname: faker.FirstName(),
// 		Lastname:  faker.LastName(),
// 		Role:      constant.USER,
// 	}
// 	newAuthSession := &model.AuthSession{
// 		UserID: existUser.ID,
// 	}
// 	credential := &authProto.Credential{
// 		AccessToken:  faker.Word(),
// 		RefreshToken: faker.Word(),
// 		ExpiresIn:    3600,
// 	}

// 	expected := &authProto.SignInResponse{Credential: credential}

// 	controller := gomock.NewController(t.T())

// 	authRepo := mock_auth.NewMockRepository(controller)
// 	userRepo := user.UserRepositoryMock{}
// 	tokenService := token.TokenServiceMock{}
// 	bcryptUtil := utils.BcryptUtilMock{}

// 	userRepo.On("FindByEmail", t.signInRequest.Email, &model.User{}).Return(existUser, nil)
// 	bcryptUtil.On("CompareHashedPassword", existUser.Password, t.signInRequest.Password).Return(nil)
// 	authRepo.EXPECT().Create(newAuthSession).Return(nil)
// 	tokenService.On("CreateCredential", existUser.ID.String(), existUser.Role, newAuthSession.ID.String()).Return(credential, nil)

// 	authSvc := NewService(authRepo, &userRepo, &tokenService, &bcryptUtil)
// 	actual, err := authSvc.SignIn(t.ctx, t.signInRequest)

// 	assert.Nil(t.T(), err)
// 	assert.Equal(t.T(), expected.Credential.AccessToken, actual.Credential.AccessToken)
// 	assert.Equal(t.T(), expected.Credential.RefreshToken, actual.Credential.RefreshToken)
// }

// func (t *AuthServiceTest) TestSignInUserNotFound() {
// 	findUserErr := gorm.ErrRecordNotFound

// 	expected := status.Error(codes.PermissionDenied, constant.IncorrectEmailPasswordErrorMessage)

// 	controller := gomock.NewController(t.T())

// 	authRepo := mock_auth.NewMockRepository(controller)
// 	userRepo := user.UserRepositoryMock{}
// 	tokenService := token.TokenServiceMock{}
// 	bcryptUtil := utils.BcryptUtilMock{}

// 	userRepo.On("FindByEmail", t.signInRequest.Email, &model.User{}).Return(nil, findUserErr)

// 	authSvc := NewService(authRepo, &userRepo, &tokenService, &bcryptUtil)
// 	actual, err := authSvc.SignIn(t.ctx, t.signInRequest)

// 	status, ok := status.FromError(err)
// 	assert.Nil(t.T(), actual)
// 	assert.Equal(t.T(), codes.PermissionDenied, status.Code())
// 	assert.True(t.T(), ok)
// 	assert.Equal(t.T(), expected.Error(), err.Error())
// }

// func (t *AuthServiceTest) TestSignInUnmatchedPassword() {
// 	existUser := &model.User{
// 		Base: model.Base{
// 			ID: uuid.New(),
// 		},
// 		Email:     t.signInRequest.Email,
// 		Password:  faker.Password(),
// 		Firstname: faker.FirstName(),
// 		Lastname:  faker.LastName(),
// 		Role:      constant.USER,
// 	}
// 	comparePwdErr := errors.New("Unmatched password")

// 	expected := status.Error(codes.PermissionDenied, constant.IncorrectEmailPasswordErrorMessage)

// 	controller := gomock.NewController(t.T())

// 	authRepo := mock_auth.NewMockRepository(controller)
// 	userRepo := user.UserRepositoryMock{}
// 	tokenService := token.TokenServiceMock{}
// 	bcryptUtil := utils.BcryptUtilMock{}

// 	userRepo.On("FindByEmail", t.signInRequest.Email, &model.User{}).Return(existUser, nil)
// 	bcryptUtil.On("CompareHashedPassword", existUser.Password, t.signInRequest.Password).Return(comparePwdErr)

// 	authSvc := NewService(authRepo, &userRepo, &tokenService, &bcryptUtil)
// 	actual, err := authSvc.SignIn(t.ctx, t.signInRequest)

// 	status, ok := status.FromError(err)
// 	assert.Nil(t.T(), actual)
// 	assert.Equal(t.T(), codes.PermissionDenied, status.Code())
// 	assert.True(t.T(), ok)
// 	assert.Equal(t.T(), expected.Error(), err.Error())
// }

// func (t *AuthServiceTest) TestSignInCreateAuthSessionFailed() {
// 	existUser := &model.User{
// 		Base: model.Base{
// 			ID: uuid.New(),
// 		},
// 		Email:     t.signInRequest.Email,
// 		Password:  faker.Password(),
// 		Firstname: faker.FirstName(),
// 		Lastname:  faker.LastName(),
// 		Role:      constant.USER,
// 	}
// 	newAuthSession := &model.AuthSession{
// 		UserID: existUser.ID,
// 	}
// 	createAuthSessionErr := errors.New("Internal server error")

// 	expected := status.Error(codes.Internal, constant.InternalServerErrorMessage)

// 	controller := gomock.NewController(t.T())

// 	authRepo := mock_auth.NewMockRepository(controller)
// 	userRepo := user.UserRepositoryMock{}
// 	tokenService := token.TokenServiceMock{}
// 	bcryptUtil := utils.BcryptUtilMock{}

// 	userRepo.On("FindByEmail", t.signInRequest.Email, &model.User{}).Return(existUser, nil)
// 	bcryptUtil.On("CompareHashedPassword", existUser.Password, t.signInRequest.Password).Return(nil)
// 	authRepo.EXPECT().Create(newAuthSession).Return(createAuthSessionErr)

// 	authSvc := NewService(authRepo, &userRepo, &tokenService, &bcryptUtil)
// 	actual, err := authSvc.SignIn(t.ctx, t.signInRequest)

// 	st, ok := status.FromError(err)
// 	assert.Nil(t.T(), actual)
// 	assert.Equal(t.T(), codes.Internal, st.Code())
// 	assert.True(t.T(), ok)
// 	assert.Equal(t.T(), expected.Error(), err.Error())
// }

// func (t *AuthServiceTest) TestSignInCreateCredentialFailed() {
// 	existUser := &model.User{
// 		Base: model.Base{
// 			ID: uuid.New(),
// 		},
// 		Email:     t.signInRequest.Email,
// 		Password:  faker.Password(),
// 		Firstname: faker.FirstName(),
// 		Lastname:  faker.LastName(),
// 		Role:      constant.USER,
// 	}
// 	newAuthSession := &model.AuthSession{
// 		UserID: existUser.ID,
// 	}
// 	createCredentialErr := errors.New("Failed to create credential")

// 	expected := status.Error(codes.Internal, constant.InternalServerErrorMessage)

// 	controller := gomock.NewController(t.T())

// 	authRepo := mock_auth.NewMockRepository(controller)
// 	userRepo := user.UserRepositoryMock{}
// 	tokenService := token.TokenServiceMock{}
// 	bcryptUtil := utils.BcryptUtilMock{}

// 	userRepo.On("FindByEmail", t.signInRequest.Email, &model.User{}).Return(existUser, nil)
// 	bcryptUtil.On("CompareHashedPassword", existUser.Password, t.signInRequest.Password).Return(nil)
// 	authRepo.EXPECT().Create(newAuthSession).Return(nil)
// 	tokenService.On("CreateCredential", existUser.ID.String(), existUser.Role, newAuthSession.ID.String()).Return(nil, createCredentialErr)

// 	authSvc := NewService(authRepo, &userRepo, &tokenService, &bcryptUtil)
// 	actual, err := authSvc.SignIn(t.ctx, t.signInRequest)

// 	status, ok := status.FromError(err)
// 	assert.Nil(t.T(), actual)
// 	assert.Equal(t.T(), codes.Internal, status.Code())
// 	assert.True(t.T(), ok)
// 	assert.Equal(t.T(), expected.Error(), err.Error())
// }

// func (t *AuthServiceTest) TestValidateSuccess() {}

// func (t *AuthServiceTest) TestValidateFailed() {}

// func (t *AuthServiceTest) TestRefreshTokenSuccess() {}

// func (t *AuthServiceTest) TestRefreshTokenNotFound() {}

// func (t *AuthServiceTest) TestRefreshTokenCreateCredentialFailed() {}

// func (t *AuthServiceTest) TestRefreshTokenUpdateTokenFailed() {}
