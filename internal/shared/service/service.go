// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package service

import (
	"fmt"
	"os"
	"path"
	"runtime/debug"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	service struct {
		Name       string `env:"SERVICE_NAME" env-default:"service"`
		Debug      bool   `env:"DEBUG" env-default:"false"`
		Secret     string `env:"SECRET" env-default:""`
		Version    string `env:"VERSION" env-default:"dev"`
		LogSkipper int    `env:"LOG_SKIPPER" env-default:"1"`
	}

	Service interface {
		GetName() string
		GetVersion() string
		GetSecret() string
		GetDebug() bool
		GetLogSkipper() int
	}

	ServiceOption func(Service)
)

func (s *service) GetName() string {
	return s.Name
}

func (s *service) GetVersion() string {
	return s.Version
}

func (s *service) GetSecret() string {
	return s.Secret
}

func (s *service) GetDebug() bool {
	return s.Debug
}

func (s *service) GetLogSkipper() int {
	return s.LogSkipper
}

// WithName sets the service name.
func WithName(name string) ServiceOption {
	return func(s Service) { s.(*service).Name = name }
}

// WithDebug sets the debug flag.
func WithDebug(debug bool) ServiceOption {
	return func(s Service) { s.(*service).Debug = debug }
}

// WithSecret sets the secret. Secret is used to sign JWT and API keys.
func WithSecret(secret string) ServiceOption {
	return func(s Service) { s.(*service).Secret = secret }
}

// WithVersion sets the version.
func WithVersion(version string) ServiceOption {
	return func(s Service) { s.(*service).Version = version }
}

func WithLogSkipper(skipper int) ServiceOption {
	return func(s Service) { s.(*service).LogSkipper = skipper }
}

// WithVersionFromBuildInfo sets the version from the build info.
func WithVersionFromBuildInfo() ServiceOption {
	return func(s Service) {
		if info, ok := debug.ReadBuildInfo(); ok {
			var (
				revision  string
				modified  string
				timestamp time.Time
				version   string
			)

			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					revision = setting.Value
				}

				if setting.Key == "vcs.modified" {
					modified = setting.Value
				}

				if setting.Key == "vcs.time" {
					timestamp, _ = time.Parse(time.RFC3339, setting.Value)
				}
			}

			if len(revision) > 0 && len(modified) > 0 && timestamp.Unix() > 0 {
				version = timestamp.Format("2006.01.02") + "." + revision[:8]
			} else {
				version = "debug"
			}

			if modified == "true" {
				version += "-dev"
			}

			s.(*service).Version = version
		}
	}
}

// WithConfigFromEnv reads the environment variables and sets the config.
func WithConfigFromEnv() ServiceOption {
	return func(s Service) {
		if err := cleanenv.ReadEnv(s.(*service)); err != nil {
			panic(fmt.Errorf("failed to read environment variables: %w", err))
		}
	}
}

// WithConfig reads the config from the given path.
func WithConfig(path string) ServiceOption {
	return func(s Service) {
		if err := cleanenv.ReadConfig(path, s.(*service)); err != nil {
			panic(fmt.Errorf("failed to read config: %w", err))
		}
	}
}

// WithConfigFromDefault reads the config from the default path.
func DefaultConfigFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("failed to get home dir: %w", err))
	}

	return path.Join(home, ".ctrlplane", "config.json")
}

// NewService creates a new instance of the service.
func NewService(opts ...ServiceOption) Service {
	s := &service{}
	for _, opt := range opts {
		opt(s)
	}

	return s
}