package github

import (
	"context"
	github_lib "github.com/google/go-github/v54/github"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
)

func CreateGithubClient(ctx context.Context) (*github_lib.Client, error) {
	authenticatedHttpClientObj, err := authenticatedHttpClientFromEnvVar(ctx)
	if err != nil {
		logrus.Warnf("Unable to build authenticated Github client from environment variable. It will now try "+
			"from AWS S3 bucket. Error was:\n%v", err.Error())
		authenticatedHttpClientObj, err = authenticatedHttpClientFromS3BucketContent(ctx)
		if err != nil {
			logrus.Warnf("Unable to build authenticated Github client from S3 bucket. Error was:\n%v", err.Error())
			return nil, stacktrace.NewError("Unable to build authenticated Github client.")
		}
	}
	githubClient := github_lib.NewClient(authenticatedHttpClientObj.Client)
	return githubClient, nil
}
