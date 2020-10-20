package controllers

import (
	"github.com/google/go-github/v32/github"
	"github.com/stretchr/testify/mock"
	"testing"
)

func (r *mockAPI) GetRunners(organization string, repository string, token string) ([]*github.Runner, error) {
	args := r.Called(organization, repository, token)
	return args.Get(0).([]*github.Runner), args.Error(1)
}

type mockAPI struct {
	mock.Mock
}

func TestGithubactionRunnerController(t *testing.T) {

}
