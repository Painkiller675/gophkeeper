package cmd

import (
	"context"
	"fmt"
	"github.com/Painkiller675/gophkeeper/internal/client/storage/sqlite"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	pb "github.com/Painkiller675/gophkeeper/internal/proto"
)

// TODO: usage??
var localDB sqlite.LocalStorage

var syncSecretCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync local DB secrets",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := secretClient.ListSecrets(context.Background(), &pb.ListSecretsRequest{})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to sync local DB secrets\n")
		}

		for _, info := range resp.GetSecrets() {
			//secret, err := models.DecodeSecret(info.GetContent())
			//if err != nil {
			//	log.Fatal().Err(err).Msg("Failed to decode(des) secret")
			//}
			//name := info.GetName()
			//version := info.GetVersion()
			// create local database instance
			localDB, err := sqlite.NewLocalStorage()
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to create local storage\n")
			}
			// feed remote data to the local sqlite database
			err = localDB.SyncLocalSecrets(context.Background(), info)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to sync local DB secrets\n")
			}

		}
		fmt.Printf("[info] local database synced with %d secrets\n", len(resp.GetSecrets()))
	},
}

func init() {
	secretCmd.AddCommand(syncSecretCmd)
}
