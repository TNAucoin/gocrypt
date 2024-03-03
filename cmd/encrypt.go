/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tnaucoin/gocrypt/internal/gokms"
	"os"
)

// encryptCmd represents the decrypt command
var encryptCmd = &cobra.Command{
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
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Encrypted file %s to %s\n", path, output)
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().StringP("profile", "p", "default", "AWS profile")
	encryptCmd.Flags().StringP("region", "r", "us-east-1", "AWS region")
	encryptCmd.Flags().StringP("role", "R", "", "AWS role")
	encryptCmd.Flags().StringP("path", "P", "", "Path to file to encrypt")
	encryptCmd.Flags().StringP("output", "o", "", "Path to file to encrypt")
	encryptCmd.Flags().StringP("key", "k", "", "KMS key alias")
	encryptCmd.Flags().StringP("ext", "e", "encrypted", "Extension of encrypted file")
	encryptCmd.MarkFlagRequired("path")
	encryptCmd.MarkFlagRequired("key")
}
