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

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypts a file",
	Long: ` Decrypts a file using the KMS key associated with the encryption method
specified in the file's metadata. The decrypted file is saved to the output path.

Example:
 gocrypt decrypt -p default -r us-east-1 -k <KEY> -e.enc -o test.txt

`,
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
		err := kms.Decrypt(path, output, key, ext)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Decrypted file %s to %s\n", path, output)
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().StringP("profile", "p", "default", "AWS profile")
	decryptCmd.Flags().StringP("region", "r", "us-east-1", "AWS region")
	decryptCmd.Flags().StringP("role", "R", "", "AWS role")
	decryptCmd.Flags().StringP("path", "P", "", "Path to file to encrypt")
	decryptCmd.Flags().StringP("output", "o", "", "Path to file to encrypt")
	decryptCmd.Flags().StringP("ext", "e", "encrypted", "Extension of encrypted file")
	decryptCmd.MarkFlagRequired("path")
}
