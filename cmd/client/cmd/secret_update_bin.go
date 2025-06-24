package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/Painkiller675/gophkeeper/internal/client/models"
	pb "github.com/Painkiller675/gophkeeper/internal/proto"
)

var updateBinSecretCmd = &cobra.Command{
	Use:   "bin",
	Short: "Update bin secret",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal().Msgf("Error reading secret name: %v", err)
			return
		}

		file, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatal().Msgf("Error reading file name: %v", err)
			return
		}

		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatal().Msgf("Error reading binary file: %v", err)
			return
		}

		bin := models.Bin{
			Data: data,
		}

		encoded, err := models.EncodeSecret(bin) // use Marshal here
		if err != nil {
			log.Fatal().Msgf("Failed to serialize secret: %v", err)
			return
		}

		resp, err := secretClient.UpdateSecret(context.Background(), &pb.UpdateSecretRequest{
			Name:    name,
			Content: encoded,
		})
		if err != nil {
			log.Fatal().Msgf("Failed to update secret: %v", err)
			return
		}

		fmt.Printf("Secret %s version %v update successfully\n", resp.GetName(), resp.GetVersion())
	},
}

func init() {
	updateSecretCmd.AddCommand(updateBinSecretCmd)

	updateBinSecretCmd.Flags().String("name", "", "Secret name")
	if err := updateBinSecretCmd.MarkFlagRequired("name"); err != nil {
		log.Error().Err(err)
	}
	updateBinSecretCmd.Flags().StringP("file", "f", "", "Binary file")
	if err := updateBinSecretCmd.MarkFlagRequired("file"); err != nil {
		log.Error().Err(err)
	}
}
