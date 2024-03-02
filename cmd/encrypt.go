/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/tnaucoin/gocrypt/internal/gokms"
	"log"

	"github.com/spf13/cobra"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// load args
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		role, _ := cmd.Flags().GetString("role")
		path, _ := cmd.Flags().GetString("path")
		output, _ := cmd.Flags().GetString("output")
		key, _ := cmd.Flags().GetString("key")
		ext, _ := cmd.Flags().GetString("ext")
		// call decrypt
		kms := gokms.New(context.Background(), profile, region, role)
		err := kms.Encrypt(path, output, key, ext)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Encrypted file %s to %s\n", path, output)
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().StringP("profile", "p", "default", "AWS profile")
	decryptCmd.Flags().StringP("region", "r", "us-east-1", "AWS region")
	decryptCmd.Flags().StringP("role", "R", "", "AWS role")
	decryptCmd.Flags().StringP("path", "P", "", "Path to file to encrypt")
	decryptCmd.Flags().StringP("output", "o", "", "Path to file to encrypt")
	decryptCmd.Flags().StringP("key", "k", "", "KMS key alias")
	decryptCmd.Flags().StringP("ext", "e", "encrypted", "Extension of encrypted file")
	decryptCmd.MarkFlagRequired("path")
	decryptCmd.MarkFlagRequired("key")
}
