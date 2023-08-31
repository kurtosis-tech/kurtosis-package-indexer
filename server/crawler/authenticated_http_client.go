package crawler

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kurtosis-tech/stacktrace"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	githubUserTokenEnvVarName = "GITHUB_USER_TOKEN"

	awsAccessKeyIdEnvVarName     = "AWS_ACCESS_KEY_ID"
	awsSecretAccessKeyEnvVarName = "AWS_SECRET_ACCESS_KEY"
	awsBucketRegionEnvVarName    = "AWS_BUCKET_REGION"
	awsBucketNameEnvVarName      = "AWS_BUCKET_NAME"
	awsBucketFolderEnvVarName    = "AWS_BUCKET_FOLDER"

	awsS3UserTokenFileName = "github-user-token.txt"
)

type AuthenticatedHttpClient struct {
	*http.Client
}

func AuthenticatedHttpClientFromEnvVar(ctx context.Context) (*AuthenticatedHttpClient, error) {
	githubUserToken := os.Getenv(githubUserTokenEnvVarName)
	if err := mustNotBeEmpty(githubUserToken); err != nil {
		return nil, stacktrace.Propagate(err, "Environment variable '%s' was empty", githubUserTokenEnvVarName)
	}
	return createAuthenticatedHttpClient(ctx, githubUserToken), nil
}

func AuthenticatedHttpClientFromS3BucketContent(ctx context.Context) (*AuthenticatedHttpClient, error) {
	awsAccessKeyId := os.Getenv(awsAccessKeyIdEnvVarName)
	awsSecretAccessKey := os.Getenv(awsSecretAccessKeyEnvVarName)
	awsBucketRegion := os.Getenv(awsBucketRegionEnvVarName)
	awsBucketName := os.Getenv(awsBucketNameEnvVarName)
	awsBucketFolder := os.Getenv(awsBucketFolderEnvVarName)
	if err := mustNotBeEmpty(awsAccessKeyId, awsSecretAccessKey, awsBucketRegion, awsBucketName); err != nil {
		return nil, stacktrace.Propagate(err, "A required environment variable was empty. All of '%s', '%s', '%s', "+
			"'%s' must be set to a non empty value to be able to connect to AWS and retrive the Github user token "+
			"from S3", awsAccessKeyIdEnvVarName, awsSecretAccessKeyEnvVarName, awsBucketRegionEnvVarName,
			awsBucketNameEnvVarName)
	}

	// nolint: exhaustruct
	awsSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsBucketRegion),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occured building the AWS client")
	}

	svc := s3.New(awsSession)

	var awsBucketFileKey string
	if awsBucketFolder == "" {
		awsBucketFileKey = fmt.Sprintf("%s/%s", awsBucketFolder, awsS3UserTokenFileName)
	} else {
		trimmedFolderName := strings.Trim(awsBucketFolder, "/")
		awsBucketFileKey = fmt.Sprintf("%s/%s", trimmedFolderName, awsS3UserTokenFileName)
	}

	// nolint: exhaustruct
	awsGetObjectResult, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(awsBucketName),
		Key:    aws.String(awsBucketFileKey),
	})
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred retrieving the file '%s' in bucket '%s'",
			awsBucketFileKey, awsBucketName)
	}

	rawGithubUserTokenFileContent, err := io.ReadAll(awsGetObjectResult.Body)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred reading the content of the S3 bucket object")
	}

	githubUserToken := strings.TrimSpace(string(rawGithubUserTokenFileContent))
	return createAuthenticatedHttpClient(ctx, githubUserToken), nil
}

func createAuthenticatedHttpClient(ctx context.Context, githubUserToken string) *AuthenticatedHttpClient {
	tokenSource := oauth2.StaticTokenSource(
		// nolint: exhaustruct
		&oauth2.Token{
			AccessToken: githubUserToken,
		},
	)
	return &AuthenticatedHttpClient{
		oauth2.NewClient(ctx, tokenSource),
	}
}

func mustNotBeEmpty(args ...string) error {
	for idx, arg := range args {
		if arg == "" {
			return stacktrace.NewError("Argument number %d was empty", idx)
		}
	}
	return nil
}
