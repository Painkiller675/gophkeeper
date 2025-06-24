package cmd

import (
	"github.com/Painkiller675/gophkeeper/internal/client/interceptors"
	pb "github.com/Painkiller675/gophkeeper/internal/proto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	secretClient pb.SecretServiceClient
)

func generateTLSCreds() (credentials.TransportCredentials, error) {
	//
	certFile := "../../internal/cert/ca-cert.pem"
	//certFile := "internal/cert/ca-cert.pem"

	return credentials.NewClientTLSFromFile(certFile, "")
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
		tlsCreds, err := generateTLSCreds()
		if err != nil {
			panic(err)
		}

		connection, err := grpc.NewClient(
			viper.GetString("grpc.address"),
			grpc.WithTransportCredentials(tlsCreds),
			grpc.WithUnaryInterceptor(interceptor.Unary()),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create client connection")
		}

		secretClient = pb.NewSecretServiceClient(connection)
	},
}

func init() {
	rootCmd.AddCommand(secretCmd)
}
