package environments

import (
	"strings"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/jx-promote/pkg/apis/promote/v1alpha1"
	"github.com/jenkins-x/jx-promote/pkg/envctx"
	"github.com/jenkins-x/jx-promote/pkg/helmer"
	"github.com/jenkins-x/jx-apps/pkg/jxapps"
	"github.com/jenkins-x/jx/v2/pkg/auth"
	"github.com/jenkins-x/jx/v2/pkg/gits"
	"github.com/jenkins-x/jx/v2/pkg/util"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

//ValuesFiles is a wrapper for a slice of values files to allow them to be passed around as a pointer
type ValuesFiles struct {
	Items []string
}

// ModifyChartFn callback for modifying a chart, requirements, the chart metadata,
// the values.yaml and all files in templates are unmarshaled, and the root dir for the chart is passed
type ModifyChartFn func(requirements *helmer.Requirements, metadata *chart.Metadata, existingValues map[string]interface{},
	templates map[string]string, dir string, pullRequestDetails *gits.PullRequestDetails) error

// ModifyAppsFn callback for modifying the `jx-apps.yml` in an environment git repository which is using helmfile and helm 3
type ModifyAppsFn func(appsConfig *jxapps.AppConfig, dir string, pullRequestDetails *gits.PullRequestDetails) error

// ModifyKptFn callback for modifying the kpt based installations of resources
type ModifyKptFn func(dir string, promoteConfig *v1alpha1.Promote, pullRequestDetails *gits.PullRequestDetails) error

// EnvironmentPullRequestOptions are options for creating a pull request against an environment.
// The provide a Gitter client for performing git operations, a GitProvider client for talking to the git provider,
// a callback ModifyChartFn which is where the changes you want to make are defined,
type EnvironmentPullRequestOptions struct {
	DevEnvContext     envctx.EnvironmentContext
	Gitter            gits.Gitter
	GitKind           string
	OutDir            string
	Function          func() error
	ModifyChartFn     ModifyChartFn
	ModifyAppsFn      ModifyAppsFn
	ModifyKptFn       ModifyKptFn
	Labels            []string
	BranchName        string
	CommitTitle       string
	CommitMessage     string
	IOFileHandles     *util.IOFileHandles
	ScmClient         *scm.Client
	AuthConfigService auth.ConfigService
	BatchMode         bool
	UseGitHubOAuth    bool
	Fork              bool
	commitBody        strings.Builder
}
