// Copyright © 2018 github.com/devopsctl authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal
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
	"github.com/spf13/cobra"
	gitlab "github.com/xanzy/go-gitlab"
)

var newReleaseCmd = &cobra.Command{
	Use:     "release",
	Aliases: []string{"r"},
	Short:   "Create a new release for the specified project's tag",
	Example: `# ensure to create the tag where the release will be created from
gitlabctl new tag v1.0 --ref=master --project=groupx/myapp

# create the release
gitlabctl new release v1.0 --project=groupx/myapp --description="Sample Release Note"`,
	SilenceErrors:     true,
	SilenceUsage:      true,
	DisableAutoGenTag: true,
	Args:              cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runNewRelease(cmd, args[0])
	},
}

func init() {
	newCmd.AddCommand(newReleaseCmd)
	addProjectFlag(newReleaseCmd)
	verifyMarkFlagRequired(newReleaseCmd, "project")
	newReleaseCmd.Flags().StringP("description", "d", "",
		"The release note or description")
	verifyMarkFlagRequired(newReleaseCmd, "description")
}

func runNewRelease(cmd *cobra.Command, tag string) error {
	opts := new(gitlab.CreateReleaseOptions)
	opts.Description = gitlab.String(getFlagString(cmd, "description"))
	createdRelease, err := newRelease(getFlagString(cmd, "project"), tag, opts)
	if err != nil {
		return err
	}
	printReleasesOut(getFlagString(cmd, "out"), createdRelease)
	return nil
}

func newRelease(project string, tag string, opts *gitlab.CreateReleaseOptions) (*gitlab.Release, error) {
	git, err := newGitlabClient()
	if err != nil {
		return nil, err
	}
	opts.TagName = &tag
	release, _, err := git.Releases.CreateRelease(project, opts)
	if err != nil {
		return nil, err
	}
	return release, nil
}
