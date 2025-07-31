package db

import (
	"context"
	"fmt"
	"os"

	infisical "github.com/infisical/go-sdk"
)

const (
	identityIdSecret  = "INFISICAL_IDENTITY_ID"
	projectSlugSecret = "INFISICAL_PROJECT_SLUG"
)

func initializeClient() (infisical.InfisicalClientInterface, error) {
	client := infisical.NewInfisicalClient(context.Background(), infisical.Config{})

	identityId, err := getEnv(identityIdSecret)
	if err != nil {
		return nil, err
	}

	_, err = client.Auth().AwsIamAuthLogin(identityId)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("environment variable %s not set", key)
	}
	return value, nil
}

func GetSecret(secretName string) (string, error) {
	client, err := initializeClient()
	if err != nil {
		return "", fmt.Errorf("failed to initialize Infisical client: %w", err)
	}

	projectSlug, err := getEnv(projectSlugSecret)
	if err != nil {
		return "", fmt.Errorf("failed to get project slug: %w", err)
	}

	// Fetch the secret using the Infisical client
	secret, err := client.Secrets().Retrieve(infisical.RetrieveSecretOptions{
		SecretKey:   secretName,
		ProjectSlug: projectSlug,
		Environment: "prod",
	})
	if err != nil {
		return "", fmt.Errorf("failed to retrieve secret: %w", err)
	}

	return secret.SecretValue, nil
}
