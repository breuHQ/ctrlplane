package github

import (
	"go.breu.io/quantm/internal/hooks/github/activities"
	"go.breu.io/quantm/internal/hooks/github/config"
	"go.breu.io/quantm/internal/hooks/github/web"
	"go.breu.io/quantm/internal/hooks/github/workflows"
)

type (
	InstallActivity      = activities.Install
	InstallReposActivity = activities.InstallRepos
	PushActivity         = activities.Push

	Config  = config.Config
	Webhook = web.Webhook
)

var (
	Install   = workflows.Install
	Push      = workflows.Push
	SyncRepos = workflows.SyncRepos
	Configure = config.Configure

	WithConfig = config.WithConfig
	Instance   = config.Instance
)