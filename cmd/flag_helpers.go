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
	"fmt"
	"strings"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	gitlab "github.com/xanzy/go-gitlab"
)

// addGetGroupsFlags adds common flags for `get groups` and `get subgroups` commands
// Flags usage reference:
// https://docs.gitlab.com/ce/api/groups.html#list-groups
// https://docs.gitlab.com/ce/api/groups.html#list-a-groups-s-subgroups
func addGetGroupsFlags(cmd *cobra.Command) {
	addAllAvailableFlag(cmd)
	addGroupOrderByFlag(cmd)
	addOwnedFlag(cmd)
	addSortFlag(cmd)
	addStatisticsFlag(cmd)
	addSearchFlag(cmd)
}

// addGetProjectsFlags adds common flags for `get projects` commands
// Flags usage reference:
// https://docs.gitlab.com/ce/api/groups.html#list-a-group-39-s-projects
func addGetProjectsFlags(cmd *cobra.Command) {
	addFromGroupFlag(cmd)
	addProjectOrderByFlag(cmd)
	addSortFlag(cmd)
	addSearchFlag(cmd)
	addStatisticsFlag(cmd)
	addVisibilityFlag(cmd)
	addOwnedFlag(cmd)
	cmd.Flags().Bool("archived", false,
		"Limit by archived status")
	cmd.Flags().Bool("simple", false,
		"Return only the ID, URL, name, and path of each project")
	cmd.Flags().Bool("membership", false,
		"Limit by projects that the current user is a member of")
	cmd.Flags().Bool("starred", false,
		"Limit by projects starred by the current user")
	cmd.Flags().Bool("with-issues-enabled", false,
		"Limit by enabled issues feature")
	cmd.Flags().Bool("with-merge-requests-enabled", false,
		"Limit by enabled merge requests feature")
}

// addNewGroupFlags add the required flags for creating a new group
// Flag usage reference: https://docs.gitlab.com/ce/api/groups.html#new-group
func addNewGroupFlags(cmd *cobra.Command) {
	addNamespaceFlag(cmd)
	addDescriptionFlag(cmd)
	addLFSenabled(cmd)
	addRequestAccessEnabledFlag(cmd)
	addVisibilityFlag(cmd)
}

// addEditGroupFlags add the required flags for updating an existing group
// Flag usage reference: https://docs.gitlab.com/ce/api/groups.html#update-group
func addEditGroupFlags(cmd *cobra.Command) {
	addChangeNameFlag(cmd)
	addChangePathFlag(cmd)
	addDescriptionFlag(cmd)
	addLFSenabled(cmd)
	addRequestAccessEnabledFlag(cmd)
	addVisibilityFlag(cmd)
}

// addEditProjectFlags add the required flags for creating a new project
// Flag usage reference:
// https://docs.gitlab.com/ce/api/projects.html#edit-project
func addEditProjectFlags(cmd *cobra.Command) {
	addNewProjectFlags(cmd)
	addChangeNameFlag(cmd)
	addChangePathFlag(cmd)
	cmd.Flags().String("default-branch", "master", "The default branch")
}

// addNewProjectFlags add the required flags for creating a new project
// Flag usage reference: https://docs.gitlab.com/ce/api/projects.html#create-project
func addNewProjectFlags(cmd *cobra.Command) {
	addDescriptionFlag(cmd)
	addLFSenabled(cmd)
	addRequestAccessEnabledFlag(cmd)
	addVisibilityFlag(cmd)
	// unique flags for projects
	cmd.Flags().Bool("issues-enabled", true, "Enable issues")
	cmd.Flags().Bool("merge-requests-enabled", true, "Enable merge requests")
	cmd.Flags().Bool("jobs-enabled", true, "Enable jobs")
	cmd.Flags().Bool("wiki-enabled", true, "Enable wiki")
	cmd.Flags().Bool("snippets-enabled", true, "Enable snippets")
	cmd.Flags().Bool("resolve-outdated-diff-discussions", false,
		"Automatically resolve merge request diffs discussions on lines "+
			"changed with a push")
	cmd.Flags().Bool("container-registry-enabled", false,
		"Enable container registry for this project")
	cmd.Flags().Bool("shared-runners-enabled", false,
		"Enable shared runners for this project")
	cmd.Flags().Bool("public-jobs", false, "If true, jobs can be viewed "+
		"by non-project-members")
	cmd.Flags().Bool("only-allow-merge-if-pipeline-succeeds", false,
		"Set whether merge requests can only be merged with successful jobs")
	cmd.Flags().Bool("only-allow-merge-if-discussion-are-resolved", false,
		"Set whether merge requests can only be merged "+
			"when all the discussions are resolved")
	cmd.Flags().String("merge-method", "merge",
		"Set the merge method used. (available: 'merge', 'rebase_merge', 'ff')")
	cmd.Flags().StringSlice("tag-list", []string{},
		"The list of tags for a project; put array of tags, "+
			"that should be finally assigned to a project.\n"+
			"Example: --tag-list='tag1,tag2'")
	cmd.Flags().Bool("printing-merge-request-link-enabled", true,
		"Show link to create/view merge request "+
			"when pushing from the command line")
	cmd.Flags().String("ci-config-path", "", "The path to CI config file")
}

func validateMergeMethodValue(cmd *cobra.Command) error {
	return validateFlagStringValue(
		[]string{"merge", "ff", "rebase_merge"},
		cmd, "merge-method")
}

func addDescriptionFlag(cmd *cobra.Command) {
	cmd.Flags().String("desc", "", "The description of the resource")
}

func addChangeNameFlag(cmd *cobra.Command) {
	cmd.Flags().String("change-name", "",
		"Use this flag to change the resource name that is "+
			"displayed in the web user interface")
}

func addChangePathFlag(cmd *cobra.Command) {
	cmd.Flags().String("change-path", "",
		"Use this flag to change the path name that is "+
			"used when accessing the resource via http or ssh url")
}

func addFromGroupFlag(cmd *cobra.Command) {
	cmd.Flags().String("from-group", "",
		"Use a group as the target namespace when performing the command")
}

func addAllAvailableFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("all-available", false,
		"Show all the groups you have access to "+
			"(defaults to false for authenticated users, true for admin)")
}

func addOwnedFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("owned", false,
		"Limit to resources owned by the current user")
}

func addGroupOrderByFlag(cmd *cobra.Command) {
	cmd.Flags().String("order-by", "name",
		"Order groups by name or path. Default is name")
}

func validateGroupOrderByFlagValue(cmd *cobra.Command) error {
	return validateFlagStringValue([]string{"path", "name"},
		cmd, "order-by")
}

func addSearchFlag(cmd *cobra.Command) {
	cmd.Flags().String("search", "",
		"Return the list of resources matching the search criteria")
}

func addStatisticsFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("statistics", false,
		"Include resource statistics (admins only)")
}

func addSortFlag(cmd *cobra.Command) {
	cmd.Flags().String("sort", "asc",
		"Order resources in asc or desc order. Default is asc")
}

func validateSortFlagValue(cmd *cobra.Command) error {
	return validateFlagStringValue([]string{"asc", "desc"},
		cmd, "sort")
}

func addProjectOrderByFlag(cmd *cobra.Command) {
	cmd.Flags().String("order-by", "created_at",
		"Return projects ordered by id, name, path, created_at, updated_at, "+
			"or last_activity_at fields. Default is created_at")
}

func validateProjectOrderByFlagValue(cmd *cobra.Command) error {
	return validateFlagStringValue([]string{"id", "name", "path",
		"created_at", "updated_at", "last_activity_at"},
		cmd, "order-by")
}

//
// NOTE(@bzon): All addFlags* helpers should be added below
//

func addNamespaceFlag(cmd *cobra.Command) {
	cmd.Flags().String("namespace", "",
		"This can be the parent namespace ID, group path, or user path. "+
			"(defaults to current user namespace)")
}

func addVisibilityFlag(cmd *cobra.Command) {
	cmd.Flags().String("visibility", "private", "public, internal or private")
}

func validateVisibilityFlagValue(cmd *cobra.Command) error {
	return validateFlagStringValue([]string{"public", "private", "internal"},
		cmd, "visibility")
}

func addRequestAccessEnabledFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("request-access-enabled", false, "Enable request access")
}

func addLFSenabled(cmd *cobra.Command) {
	cmd.Flags().Bool("lfs-enabled", false, "Enable LFS")
}

// TODO(@bzon): to be deleted soon
// currently used by group-member that will be refactored
func addPathFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("path", "p", "",
		"the group name, id or full the path "+
			"including the parent group (path/to/group)")
	if err := cmd.MarkFlagRequired("path"); err != nil {
		panic(err)
	}
}

func addOutFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("out", "o", "simple",
		"Print the command output to the "+
			"desired format. (json, yaml, simple)")
}

func validateOutFlagValue(cmd *cobra.Command) error {
	return validateFlagStringValue([]string{JSON, YAML, "simple"},
		cmd, "out")
}

func validateFlagStringValue(stringSlice []string,
	cmd *cobra.Command, fName string) error {
	fValue := getFlagString(cmd, fName)
	for _, v := range stringSlice {
		if fValue == v {
			return nil
		}
	}
	return fmt.Errorf("'%s' is not a recognized value of '%s' flag. "+
		"Please choose from: [%s]\n",
		fValue, fName, strings.Join(stringSlice, ", "))
}

//
// NOTE(@bzon): All getFlag* helpers should be added below
//

// getFlagVisibility converts the string flag visiblity to gitlab.VisibilityValue.
func getFlagVisibility(cmd *cobra.Command) *gitlab.VisibilityValue {
	v := getFlagString(cmd, "visibility")
	return gitlab.Visibility(gitlab.VisibilityValue(v))
}

// getFlagMergeMethod converts the string flag merge-method to gitlab.MergeMethod
func getFlagMergeMethod(cmd *cobra.Command) *gitlab.MergeMethodValue {
	v := getFlagString(cmd, "merge-method")
	return gitlab.MergeMethod(gitlab.MergeMethodValue(v))
}

func getFlagStringSlice(cmd *cobra.Command, flag string) []string {
	s, err := cmd.Flags().GetStringSlice(flag)
	if err != nil {
		glog.Fatalf("error accessing flag %s for command %s: %v",
			flag, cmd.Name(), err)
	}
	return s
}

func getFlagString(cmd *cobra.Command, flag string) string {
	s, err := cmd.Flags().GetString(flag)
	if err != nil {
		glog.Fatalf("error accessing flag %s for command %s: %v",
			flag, cmd.Name(), err)
	}
	return s
}

func getFlagBool(cmd *cobra.Command, flag string) bool {
	b, err := cmd.Flags().GetBool(flag)
	if err != nil {
		glog.Fatalf("error accessing flag %s for command %s: %v",
			flag, cmd.Name(), err)
	}
	return b
}
