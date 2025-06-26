package cmd

import (
	pb "github.com/Painkiller675/gophkeeper/internal/proto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	authClient pb.AuthServiceClient
)

// authCmd the subcommand
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage user registration, authentication and authorization",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// get the credentials object
		tlsCreds, err := generateTLSCreds()
		if err != nil {
			panic(err)
		}
		connection, err := grpc.NewClient(
			viper.GetString("grpc.address"),
			grpc.WithTransportCredentials(tlsCreds))
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create client connection")
		}

		authClient = pb.NewAuthServiceClient(connection)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
