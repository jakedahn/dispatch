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

package core

import git "gopkg.in/libgit2/git2go.v22"

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
