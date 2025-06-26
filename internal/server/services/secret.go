package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	//sCm "github.com/Painkiller675/gophkeeper/cmd/server/cmd"
	clMod "github.com/Painkiller675/gophkeeper/internal/client/models"
	pb "github.com/Painkiller675/gophkeeper/internal/proto"
	"github.com/Painkiller675/gophkeeper/internal/server/config"
	"github.com/Painkiller675/gophkeeper/internal/server/interceptors"
	"github.com/Painkiller675/gophkeeper/internal/server/models"
	"github.com/Painkiller675/gophkeeper/internal/server/storage"
	"github.com/Painkiller675/gophkeeper/internal/server/storage/pg"
)

//var (
//	BlockCipher cipher.BlockCipher
//)

/*
	func init() {
		// let's create the cipher block
		cipher, err := gcm.New(viper.GetString("encryption.key"))
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create cipher")
		}
		BlockCipher = cipher
	}
*/

// SecretService is an implementation of proto.SecretServiceServer
type SecretService struct {
	SecretStorage storage.SecretStorage
	pb.UnimplementedSecretServiceServer
}

var _ pb.SecretServiceServer = (*SecretService)(nil)
var _ Service = (*SecretService)(nil)

// NewSecretService creates a new service SecretService
func NewSecretService(cfg config.Config) *SecretService {
	secretStorage, err := pg.NewSecretStorage(cfg.DB.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create storage")
	}

	return &SecretService{SecretStorage: secretStorage}
}

// auxiliary functions for encryptions / decryptions

func encryptSecret(s clMod.Secret) ([]byte, error) {
	encoded, err := clMod.EncodeSecret(s) // use Marshal here
	if err != nil {
		return nil, err
	}
	return config.BlockC.Encrypt(encoded)
}

// RegisterService register service SecretService on a gRPC server
func (srv *SecretService) RegisterService(s grpc.ServiceRegistrar) {
	pb.RegisterSecretServiceServer(s, srv)
}

// GetSecret return user's private data by the name in request
func (srv *SecretService) GetSecret(
	ctx context.Context,
	request *pb.GetSecretRequest,
) (*pb.GetSecretResponse, error) {
	if request.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "secret name is empty")
	}

	userID, ok := ctx.Value(interceptors.ContextKeyUserID).(int)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "empty user id")
	}

	secret, err := srv.SecretStorage.GetSecret(ctx, request.GetName(), userID)
	if err != nil {
		if errors.Is(err, storage.ErrSecretNotFound) {
			return nil, status.Error(codes.NotFound, "secret not found")
		}
		return nil, status.Error(codes.Internal, "failed to get secret")
	}
	// ******************
	// decrypt secret
	decr_sec, err := config.BlockC.Decrypt(secret.Content)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to decrypt secret")
	}

	return &pb.GetSecretResponse{
		Name:    secret.Name,
		Content: decr_sec,
		Version: secret.Version.String(),
	}, nil
}

// CreateSecret saves new user's private data on the server
func (srv *SecretService) CreateSecret(
	ctx context.Context,
	request *pb.CreateSecretRequest,
) (*pb.CreateSecretResponse, error) {
	if request.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty secret name")
	}
	if len(request.GetContent()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty secret content")
	}

	userID, ok := ctx.Value(interceptors.ContextKeyUserID).(int)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "empty user id")
	}

	// initial decryption (deserialization)
	//init_dec, err := clMod.DecodeSecret(request.GetContent())
	//if err != nil {
	//	log.Fatal().Msgf("Failed to decrypt secret: %v", err)
	//	return nil, status.Error(codes.Internal, "failed to decode secret")
	//}
	// let's encrypt it to store
	content, err := config.BlockC.Encrypt(request.GetContent()) // TODO: replace
	if err != nil {
		log.Fatal().Msgf("Failed to encrypt secret: %v", err)
		return nil, status.Error(codes.Internal, "failed to create secret (encrypt)")
	}
	secret, err := srv.SecretStorage.CreateSecret(ctx, &models.Secret{
		Name:    request.GetName(),
		Content: content,
		Version: uuid.UUID{},
		OwnerID: userID,
	})
	if err != nil {
		if errors.Is(err, storage.ErrSecretConflict) {
			return nil, status.Error(codes.AlreadyExists, "secret already exists")
		}
		return nil, status.Error(codes.Internal, "failed to create secret")
	}
	return &pb.CreateSecretResponse{
		Name:    request.GetName(),
		Version: secret.Version.String(),
	}, nil
}

// UpdateSecret updates user's private data (secret)
func (srv *SecretService) UpdateSecret(
	ctx context.Context,
	request *pb.UpdateSecretRequest,
) (*pb.UpdateSecretResponse, error) {
	if request.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty secret name")
	}
	if len(request.GetContent()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty secret content")
	}

	userID, ok := ctx.Value(interceptors.ContextKeyUserID).(int)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "empty user id")
	}

	// let's encrypt the updating info
	content, err := config.BlockC.Encrypt(request.GetContent()) // TODO: replace
	if err != nil {
		log.Fatal().Msgf("Failed to encrypt secret: %v", err)
		return nil, status.Error(codes.Internal, "failed to create secret (encrypt)")
	}

	secret, err := srv.SecretStorage.UpdateSecret(ctx, &models.Secret{
		Name:    request.GetName(),
		Content: content,
		OwnerID: userID,
	})
	if err != nil {
		if errors.Is(err, storage.ErrSecretNotFound) {
			return nil, status.Error(codes.NotFound, "secret not found")
		}
		return nil, status.Error(codes.Internal, "failed to create secret")
	}
	return &pb.UpdateSecretResponse{
		Name:    request.GetName(),
		Version: secret.Version.String(),
	}, nil
}

// DeleteSecret deletes user's private data
func (srv *SecretService) DeleteSecret(
	ctx context.Context,
	request *pb.DeleteSecretRequest,
) (*pb.DeleteSecretResponse, error) {
	if request.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty secret name")
	}

	userID, ok := ctx.Value(interceptors.ContextKeyUserID).(int)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "empty user id")
	}

	secret := &models.Secret{
		Name:    request.GetName(),
		OwnerID: userID,
	}

	if err := srv.SecretStorage.DeleteSecret(ctx, secret); err != nil {
		if errors.Is(err, storage.ErrSecretNotFound) {
			return nil, status.Error(codes.NotFound, "secret not found")
		}
		return nil, status.Error(codes.Internal, "failed to create secret")
	}
	return &pb.DeleteSecretResponse{
		Name: request.GetName(),
	}, nil
}

// ListSecrets returns user's secret list
func (srv *SecretService) ListSecrets(
	ctx context.Context,
	_ *pb.ListSecretsRequest,
) (*pb.ListSecretsResponse, error) {
	userID, ok := ctx.Value(interceptors.ContextKeyUserID).(int)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "empty user id")
	}

	secrets, err := srv.SecretStorage.ListSecrets(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list secrets")
	}

	pbSecrets := make([]*pb.SecretInfo, 0, len(secrets))
	for _, secret := range secrets {
		// decrypt each secret
		decr_sec, err := config.BlockC.Decrypt(secret.Content)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to decrypt secret")
		}
		pbSecrets = append(pbSecrets, &pb.SecretInfo{
			Name:    secret.Name,
			Content: decr_sec,
			Version: secret.Version.String(),
		})
	}
	return &pb.ListSecretsResponse{
		Secrets: pbSecrets,
	}, nil
}
