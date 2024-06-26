// Copyright © 2018 github.com/devopsctl authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	gitlab "github.com/xanzy/go-gitlab"
)

var getProjectsCmd = &cobra.Command{
	Use:               "projects",
	Aliases:           []string{"p"},
	SuggestFor:        []string{"project"},
	Short:             "List projects of the authenticated user or of a group",
	SilenceErrors:     true,
	SilenceUsage:      true,
	DisableAutoGenTag: true,
	Example: `# get all projects
gitlabctl get projects

# get all projects from a group
gitlabctl get projects --from-group=Group1`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := validateSortFlagValue(cmd); err != nil {
			return err
		}
		if err := validateProjectOrderByFlagValue(cmd); err != nil {
			return err
		}
		return validateVisibilityFlagValue(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if getFlagString(cmd, "from-group") != "" {
			return runGetProjectsFromGroup(cmd)
		}
		return runGetProjects(cmd)
	},
}

func init() {
	getCmd.AddCommand(getProjectsCmd)
	addGetProjectsFlags(getProjectsCmd)
}

func runGetProjects(cmd *cobra.Command) error {
	opts := assignListProjectOptions(cmd)
	projects, err := getProjects(opts)
	if err != nil {
		return err
	}
	printProjectsOut(getFlagString(cmd, "out"), projects...)
	return nil
}

func getProjects(opts *gitlab.ListProjectsOptions) ([]*gitlab.Project, error) {
	git, err := newGitlabClient()
	if err != nil {
		return nil, err
	}
	projects, _, err := git.Projects.ListProjects(opts)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func runGetProjectsFromGroup(cmd *cobra.Command) error {
	optstr, err := json.Marshal(assignListProjectOptions(cmd))
	if err != nil {
		return err
	}
	opt := &gitlab.ListGroupProjectsOptions{}
	if err = json.Unmarshal(optstr, opt); err != nil {
		return err
	}
	projects, err := getProjectsFromGroup(getFlagString(cmd, "from-group"), opt)
	if err != nil {
		return err
	}
	printProjectsOut(getFlagString(cmd, "out"), projects...)
	return nil
}

func getProjectsFromGroup(group string, opts *gitlab.ListGroupProjectsOptions) ([]*gitlab.Project, error) {
	git, err := newGitlabClient()
	if err != nil {
		return nil, err
	}
	projects, _, err := git.Groups.ListGroupProjects(group, opts)
	if err != nil {
		return nil, err
	}
	return projects, nil
}
