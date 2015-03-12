// Copyright 2015 Jake Dahn
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	git "gopkg.in/libgit2/git2go.v22"
	yaml "gopkg.in/yaml.v2"
)

type Repo struct {
	Name         string
	GitUrl       string
	GitRef       string
	GitRepo      *git.Repository
	CheckoutPath string
}

type DispatchFile struct {
	BuildSteps []string               `yaml:"build"`
	Arguments  []DispatchfileArgument `yaml:"arguments"`
}

type DispatchfileArgument struct {
	Key      string
	Presence string
}

type DispatchRequest struct {
	GitUrl    string              `json:"repo_url"`
	GitRef    string              `json:"ref"`
	Arguments []map[string]string `json:"arguments"`
}

func NewRepo(gitUrl, gitRef string) *Repo {
	repoUrlParts := strings.Split(gitUrl, "/")
	repoUrlParts = strings.Split(repoUrlParts[len(repoUrlParts)-1], ".")
	name := repoUrlParts[len(repoUrlParts)-2]

	repo := &Repo{
		GitUrl: gitUrl,
		GitRef: gitRef,
		Name:   name,
	}
	return repo
}

func (r *Repo) Init() error {
	if r.GitUrl == "" || r.GitRef == "" {
		return errors.New("r.GitUrl or r.Gitref are undefined")
	}

	tmpDir := "/tmp/dispatch-tmp"
	checkoutPath, err := ioutil.TempDir(tmpDir, r.Name)
	if err != nil {
		return fmt.Errorf("Unable to create directory: %s", err)
	}

	r.CheckoutPath = checkoutPath

	cloneOpts := &git.CloneOptions{}
	gitRepo, err := git.Clone(r.GitUrl, checkoutPath, cloneOpts)
	if err != nil {
		return fmt.Errorf("Unable to clone the repo: %s", err)
	}
	r.GitRepo = gitRepo

	return nil
}

func (r *Repo) Checkout() error {
	if r.GitUrl == "" || r.GitRef == "" || r.CheckoutPath == "" {
		return errors.New("r.GitUrl or r.Gitref are undefined or r.CheckoutPath")
	}

	sha, err := git.NewOid(r.GitRef)
	if err != nil {
		return err
	}

	commit, err := r.GitRepo.LookupCommit(sha)
	if err != nil {
		return err
	}

	err = r.GitRepo.SetHeadDetached(commit.Id(), nil, "")
	if err != nil {
		return err
	}

	return nil
}

func NewDispatchRequest(requestJson []byte) (*DispatchRequest, error) {
	dispatchRequest := &DispatchRequest{}
	err := json.Unmarshal(requestJson, dispatchRequest)
	if err != nil {
		return nil, err
	}
	return dispatchRequest, nil
}

func ParseDispatchFile(dispatchFileContent []byte) (*DispatchFile, error) {
	df := &DispatchFile{}
	err := yaml.Unmarshal(dispatchFileContent, df)
	if err != nil {
		return nil, err
	}
	return df, nil
}

func (d *DispatchRequest) IsValid(dispatchFile *DispatchFile) bool {
	dfArgs := dispatchFile.Arguments
	reqArgs := d.Arguments[0]

	for _, dfarg := range dfArgs {
		if dfarg.Presence == "required" {
			if reqArgs[dfarg.Key] == "" {
				return false
			}
		}
	}
	return true
}
