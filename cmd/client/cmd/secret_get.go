package cmd

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/Painkiller675/gophkeeper/internal/client/models"
	"github.com/Painkiller675/gophkeeper/internal/client/storage/sqlite"
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
			// decode content to show
			decCon, err := models.DecodeSecret(localSec.Content)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to decode secret from the local database")
			}
			// show the info needed
			fmt.Printf("%s ---  %s --- %s\n", localSec.Name, localSec.Version, decCon)
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
