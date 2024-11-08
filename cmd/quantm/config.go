package main

import (
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/durable"
	githubcfg "go.breu.io/quantm/internal/hooks/github/config"
	pkg_slack "go.breu.io/quantm/internal/hooks/slack/config"
	"go.breu.io/quantm/internal/nomad"
)

type (
	Config struct {
		DB      *db.Config        `koanf:"DB"`      // Configuration for the database.
		Durable *durable.Config   `koanf:"DURABLE"` // Configuration for the durable.
		Nomad   *nomad.Config     `koanf:"NOMAD"`   // Configuration for Nomad.
		Github  *githubcfg.Config `koanf:"GITHUB"`  // Configuration for the github.
		Slack   *pkg_slack.Config `koanf:"SLACK"`   // Configuration for the slack.

		Secret  string `koanf:"SECRET"`  // Secret key for JWE.
		Debug   bool   `koanf:"DEBUG"`   // Flag to enable debug mode.
		Migrate bool   `koanf:"MIGRATE"` // Flag to enable database migration.
	}
)

func (c *Config) Load() {
	c.DB = &db.DefaultConfig
	c.Durable = &durable.DefaultConfig
	c.Nomad = &nomad.DefaultConfig
	c.Github = &githubcfg.Config{}
	c.Slack = &pkg_slack.Config{}

	k := koanf.New("__")

	// Load default values from the Config struct.
	if err := k.Load(structs.Provider(c, "__"), nil); err != nil {
		panic(err)
	}

	// Load environment variables with the "__" delimiter.
	if err := k.Load(env.Provider("", "__", nil), nil); err != nil {
		panic(err)
	}

	// Unmarshal configuration from the Koanf instance to the Config struct.
	if err := k.Unmarshal("", c); err != nil {
		panic(err)
	}

	// Add -m or --migrate flag to enable database migration.
	if !c.Migrate {
		flag.BoolVarP(&c.Migrate, "migrate", "m", false, "run database migrations")
		flag.Parse()
	}
}