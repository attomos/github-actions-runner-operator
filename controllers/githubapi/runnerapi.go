package githubapi

import (
	"context"
	"github.com/google/go-github/v32/github"
	"github.com/palantir/go-githubapp/githubapp"
	"time"
)

//IRunnerAPI is a service towards GitHubs runners
type IRunnerAPI interface {
	GetRunners(organization string, repository string, token string) ([]*github.Runner, error)
}

type runnerAPI struct {
}

//NewRunnerAPI gets a new instance of the API.
func NewRunnerAPI() runnerAPI {
	return runnerAPI{}
}

func getClient(organization string, token string) (*github.Client, error) {
	config := githubapp.Config{}
	config.SetValuesFromEnv("")

	clientCreator, err := githubapp.NewDefaultCachingClientCreator(config,
		githubapp.WithClientUserAgent("GithubActionsRunnerOperator"),
		githubapp.WithClientTimeout(time.Second*4),
	)
	if err != nil {
		return nil, err
	}

	if config.App.PrivateKey != "" {
		appClient, err := clientCreator.NewAppClient()
		if err != nil {
			return nil, err
		}

		installService := githubapp.NewInstallationsService(appClient)
		installation, err := installService.GetByOwner(context.TODO(), organization)
		if err == nil {
			return clientCreator.NewInstallationClient(installation.ID)
		}
	} else {
		return clientCreator.NewTokenClient(token)
	}

	return nil, err
}

// Return all runners for the org
func (r runnerAPI) GetRunners(organization string, repository string, token string) ([]*github.Runner, error) {
	client, err := getClient(organization, token)
	if err != nil {
		return nil, err
	}

	var allRunners []*github.Runner
	opts := &github.ListOptions{PerPage: 30}

	for {
		var runners *github.Runners
		var response *github.Response
		var err error

		if repository != "" {
			runners, response, err = client.Actions.ListRunners(context.TODO(), organization, repository, opts)
		} else {
			runners, response, err = client.Actions.ListOrganizationRunners(context.TODO(), organization, opts)
		}
		if err != nil {
			return allRunners, err
		}
		allRunners = append(allRunners, runners.Runners...)
		if response.NextPage == 0 {
			break
		}
		opts.Page = response.NextPage
	}

	return allRunners, nil
}
