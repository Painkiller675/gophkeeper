package cmd

import (
	"context"
	"fmt"
	"github.com/Painkiller675/gophkeeper/internal/client/models"
	"github.com/Painkiller675/gophkeeper/internal/client/storage/sqlite"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	pb "github.com/Painkiller675/gophkeeper/internal/proto"
)

var getSecretCmd = &cobra.Command{
	Use:   "get",
	Short: "Get secret",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal().Msgf("Error reading secret name: %v", err)
		}

		resp, err := secretClient.GetSecret(context.Background(), &pb.GetSecretRequest{
			Name: name,
		})
		if err != nil {
			//log.Fatal().Err(err).Msg("Failed to get secret")
			log.Info().Msgf("Can't get secret from remote database: %v", err)
			// try to get smth from a local one
			// create local database instance
			localDB, err := sqlite.NewLocalStorage()
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to create local storage")
			}
			localSec, err := localDB.GetLocalSecret(context.Background(), name)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get secret from the local database")
			}
			decCon, err := models.DecodeSecret(localSec.Content)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to decode secret from the local database")
			}
			fmt.Printf("%s ---  %s\n", decCon, localSec.Version)
			return
		} //
		// I've changed decrypt to deserialization here
		//secret, err := decryptSecret(resp.GetContent()) // TODO: REMOVE
		//if err != nil {
		//	log.Fatal().Err(err).Msg("Failed to decrypt secret")
		//}
		secret, err := models.DecodeSecret(resp.GetContent())
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to decrypt secret (deserialization)")
		}

		fmt.Printf("%s\n", secret)
	},
}

func init() {
	secretCmd.AddCommand(getSecretCmd)

	getSecretCmd.Flags().String("name", "", "Secret name")
	if err := getSecretCmd.MarkFlagRequired("name"); err != nil {
		log.Error().Err(err)
	}
}
