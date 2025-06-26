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

var listSecretCmd = &cobra.Command{
	Use:   "list",
	Short: "List secrets",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := secretClient.ListSecrets(context.Background(), &pb.ListSecretsRequest{})
		if err != nil {
			//log.Fatal().Err(err).Msg("Failed to list secret")

			// try to get the list from a local database
			// create local database instance
			localDB, err := sqlite.NewLocalStorage()
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to create local storage")
			}
			localSec, err := localDB.GetLocalList(context.Background())
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get secret from the local database")
			}
			// decode content to show
			for _, sec := range localSec {
				decCon, err := models.DecodeSecret(sec.Content)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to decode secret from the local database")
				}
				// show each secret
				fmt.Printf("%s ---  %s --- %s\n", sec.Name, sec.Version, decCon)
			}
			return

		}

		for _, info := range resp.GetSecrets() {
			secret, err := models.DecodeSecret(info.GetContent())
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to decode(des) secret")
			}

			fmt.Printf("%s\n", secret)
		}
	},
}

func init() {
	secretCmd.AddCommand(listSecretCmd)
}
