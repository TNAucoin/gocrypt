package gokms

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"io"
	"os"
	"strings"
)

// KMS is a wrapper around the AWS KMS client.
type KMS struct {
	ctx    context.Context
	cfg    aws.Config
	client *kms.Client
}

// New creates a new KMS client.
func New(ctx context.Context, profile, region, role string) *KMS {
	if ctx == nil {
		ctx = context.Background()
	}
	cfg, err := loadConfig(region, profile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// If no role is specified, use the default credentials.
	if role == "" {
		client := kms.NewFromConfig(cfg)
		return &KMS{client: client, cfg: cfg, ctx: ctx}
	}

	creds, err := assumeRole(ctx, cfg, role)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	client := kms.New(kms.Options{
		Credentials: credentials.NewStaticCredentialsProvider(*creds.AccessKeyId, *creds.SecretAccessKey, *creds.SessionToken),
		Region:      region,
	})

	return &KMS{client: client, cfg: cfg, ctx: ctx}
}

// Encrypt encrypts the plaintext.
func (k *KMS) Encrypt(path, output, key, ext string) error {
	data, err := readFile(path)
	if err != nil {
		return err
	}

	ciphertext, err := k.client.Encrypt(k.ctx, &kms.EncryptInput{
		KeyId:     aws.String(key),
		Plaintext: data,
	})

	// if output is not specified, write the encrypted data to the same path as the plaintext.
	if output != "" {
		err = writeFile(output, ext, ciphertext.CiphertextBlob, true)
	} else {
		err = writeFile(path, ext, ciphertext.CiphertextBlob, true)
	}

	if err != nil {
		return err
	}

	return nil
}

// Decrypt decrypts the ciphertext.
func (k *KMS) Decrypt(path, output, key, ext string) error {
	data, err := readFile(path)
	if err != nil {
		return err
	}
	ciphertext, err := k.client.Decrypt(k.ctx, &kms.DecryptInput{
		CiphertextBlob: data,
	})
	if err != nil {
		return err
	}
	if output != "" {
		err = writeFile(output, ext, ciphertext.Plaintext, false)
	} else {
		err = writeFile(path, ext, ciphertext.Plaintext, false)
	}
	if err != nil {
		return err
	}
	return nil
}

// readFile reads the specified file. If the file does not exist, it returns an error.
func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// writeFile writes the ciphertext to the specified file. If the file already exists, it is renamed and a new file is created.
// If isEncrypted is true, the file extension is.enc. Otherwise, the file extension is removed.
func writeFile(path, ext string, ciphertext []byte, isEncrypted bool) error {
	var writeFilePath string
	if isEncrypted {
		writeFilePath = path + "." + ext
	} else {
		writeFilePath = strings.TrimSuffix(path, "."+ext)
	}

	if _, err := os.Stat(writeFilePath); !os.IsNotExist(err) {
		// File already exists, so don't overwrite it.
		// Rename the existing file so we can write the new file.
		oldPath := writeFilePath + ".old"
		if err := os.Rename(writeFilePath, oldPath); err != nil {
			return err
		}
	}

	if err := os.WriteFile(writeFilePath, ciphertext, 0644); err != nil {
		return err
	}
	return nil
}

// loadConfig loads the configuration.
func loadConfig(region, profile string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile), config.WithRegion(region))
	if err != nil {
		return aws.Config{}, err
	}
	return cfg, nil
}

// assumeRole assumes the role and returns the credentials.
func assumeRole(ctx context.Context, cfg aws.Config, roleArn string) (*types.Credentials, error) {
	svc := sts.NewFromConfig(cfg)
	input := &sts.AssumeRoleInput{
		RoleArn: aws.String(roleArn), RoleSessionName: aws.String("gocrypt"),
	}
	resp, err := svc.AssumeRole(ctx, input)
	if err != nil {
		return nil, err
	}
	return resp.Credentials, nil
}
