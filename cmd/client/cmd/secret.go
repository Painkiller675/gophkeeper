package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"

	"github.com/Painkiller675/gophkeeper/internal/client/interceptors"
	"github.com/Painkiller675/gophkeeper/internal/client/models"
	pb "github.com/Painkiller675/gophkeeper/internal/proto"
	"github.com/Painkiller675/gophkeeper/pkg/cipher"
	"github.com/Painkiller675/gophkeeper/pkg/cipher/aes/gcm"
)

var (
	secretClient pb.SecretServiceClient
	blockCipher  cipher.BlockCipher
)

func encryptSecret(s models.Secret) ([]byte, error) {
	encoded, err := models.EncodeSecret(s)
	if err != nil {
		return nil, err
	}
	return blockCipher.Encrypt(encoded)
}

func decryptSecret(b []byte) (models.Secret, error) {
	encoded, err := blockCipher.Decrypt(b)
	if err != nil {
		return nil, err
	}
	return models.DecodeSecret(encoded)
}

func LaadTLSCredentials() (credentials.TransportCredentials, error) {
	// load the certificate of the CA who signed server's certificate
	pemServerCA, err := os.ReadFile("../../internal/cert/ca-cert.pem") // viper.GetString("tls.server-ca")
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA to cert pool")
	}

	config := &tls.Config{
		RootCAs: certPool,
	}
	return credentials.NewTLS(config), nil
}

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage user private data",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		accessToken, err := tokenStorage.Load()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to load access token")
		}
		if accessToken == "" {
			log.Fatal().Msg("Empty access token")
		}
		interceptor := interceptors.NewAuthInterceptor(viper.GetString("token"))

		// get the credentials object
		tlsCredentials, err := LaadTLSCredentials()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to load TLS credentials")
		}

		connection, err := grpc.NewClient(
			viper.GetString("grpc.address"),
			grpc.WithTransportCredentials(tlsCredentials),
			grpc.WithUnaryInterceptor(interceptor.Unary()),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create client connection")
		}

		secretClient = pb.NewSecretServiceClient(connection)
		cipher, err := gcm.New(viper.GetString("encryption.key"))
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create cipher")
		}
		blockCipher = cipher
	},
}

func init() {
	rootCmd.AddCommand(secretCmd)
}
