// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2022, 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package db

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Guilospanck/gocqlxmock"
	"github.com/Guilospanck/igocqlx"
	"github.com/avast/retry-go/v4"
	"github.com/go-playground/validator/v10"
	"github.com/gocql/gocql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cassandra"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/scylladb/go-reflectx"

	"go.breu.io/quantm/internal/shared"

	_ "github.com/golang-migrate/migrate/v4/source/file" // required for file:// migrations
)

const (
	NullUUID     = "00000000-0000-0000-0000-000000000000" // RFC 9562 - NULL UUID
	MaxUUID      = "ffffffff-ffff-ffff-ffff-ffffffffffff" // RFC 9562 - MAX UUID
	NullString   = ""
	TestKeyspace = "ctrlplane_test"
)

var (
	db   *Config
	once sync.Once
)

var (
	//go:embed migrations/*.cql
	src embed.FS
)

type (
	// Config holds the information about the database.
	// Config holds the information about the database.
	Config struct {
		Session            igocqlx.ISessionx // Initialized session.
		Hosts              []string          `env:"CASSANDRA_HOSTS" env-default:"database"`     // List of Cassandra nodes to connect to.
		Port               int               `env:"CASSANDRA_PORT" env-default:"9042"`          // Port of the Cassandra cluster.
		User               string            `env:"CASSANDRA_USER" env-default:""`              // Username to authenticate with Cassandra.
		Password           string            `env:"CASSANDRA_PASS" env-default:""`              // Password to authenticate with Cassandra.
		Keyspace           string            `env:"CASSANDRA_KEYSPACE" env-default:"ctrlplane"` // Keyspace to use.
		MigrationSourceURL string            `env:"CASSANDRA_MIGRATION_SOURCE_URL"`             // URL for migrations.
		Timeout            time.Duration     `env:"CASSANDRA_TIMEOUT" env-default:"1m"`         // Default timeout for database operations.
	}

	// MockConfig represents the mock session.
	MockConfig struct {
		*gocqlxmock.SessionxMock
	}

	ConfigOption func(*Config)

	MigrationLogger struct{}
)

// Printf implements the logger interface.
func (l *MigrationLogger) Printf(format string, v ...any) {
	slog.Info(fmt.Sprintf(format, v...))
}

// Verbose implements the logger interface.
func (l *MigrationLogger) Verbose() bool {
	return false
}

// Session implements the ISessionx interface.
func (mc *MockConfig) Session() *igocqlx.Session {
	return nil
}

func (mc *MockConfig) SetMapper(mapper *reflectx.Mapper) {}

// Shutdown closes the database session.
func (c *Config) Shutdown(ctx context.Context) error {
	c.Session.Session().S.Session.Close()

	return nil
}

// WithHosts sets the hosts.
func WithHosts(hosts []string) ConfigOption {
	return func(c *Config) { c.Hosts = hosts }
}

// WithPort sets the port.
func WithPort(port int) ConfigOption {
	return func(c *Config) { c.Port = port }
}

// WithKeyspace sets the keyspace.
func WithKeyspace(keyspace string) ConfigOption {
	return func(c *Config) { c.Keyspace = keyspace }
}

// WithMigrationSourceURL sets the migration source URL.
func WithMigrationSourceURL(url string) ConfigOption {
	return func(c *Config) { c.MigrationSourceURL = url }
}

// WithTimeout sets the timeout.
func WithTimeout(timeout time.Duration) ConfigOption {
	return func(c *Config) { c.Timeout = timeout }
}

// FromEnvironment reads the configuration from the environment.
func FromEnvironment() ConfigOption {
	return func(c *Config) {
		if err := cleanenv.ReadEnv(c); err != nil {
			panic(fmt.Errorf("db: unable to read environment variables, %v", err))
		}
	}
}

// WithSessionCreation initializes the session.
func WithSessionCreation() ConfigOption {
	return func(c *Config) {
		if c.Hosts == nil || c.Keyspace == "" {
			panic(fmt.Errorf("db: hosts & keyspace not set, please set them before initializing session"))
		}

		cluster := gocql.NewCluster(c.Hosts...)
		cluster.Keyspace = c.Keyspace
		cluster.Timeout = c.Timeout
		cluster.ConnectTimeout = c.Timeout

		if c.User != "" && c.Password != "" {
			cluster.Authenticator = gocql.PasswordAuthenticator{
				Username: c.User,
				Password: c.Password,
			}
		}

		createSession := func() error {
			slog.Info("db: connecting ...", "hosts", c.Hosts, "keyspace", c.Keyspace)

			session, err := igocqlx.WrapSession(cluster.CreateSession())
			if err != nil {
				slog.Error("db: failed to connect", "error", err)
				return err
			}

			session.SetMapper(CQLMapper)

			c.Session = session

			slog.Info("db: connected")

			return nil
		}

		if err := retry.Do(
			createSession,
			retry.Attempts(15),
			retry.Delay(6*time.Second),
		); err != nil {
			slog.Error("db: aborting ....", "error", err)
		}
	}
}

// WithE2ESession initializes the session for end-to-end tests.
//
// NOTE: It might appear that the client throws error as explained at [issue], which will eventially point to [gocql github],
// but IRL, it will work. This is a known issue with gocql and it's not a problem for us.
//
// [issue]: https://app.shortcut.com/ctrlplane/story/2509/migrate-testing-to-use-test-containers-instead-of-mocks#activity-2749
// [gocql github]: https://github.com/gocql/gocql/issues/575
func WithE2ESession() ConfigOption {
	return func(c *Config) {
		cluster := gocql.NewCluster(c.Hosts...)
		cluster.Keyspace = c.Keyspace
		cluster.Timeout = 10 * time.Minute
		cluster.ConnectTimeout = 10 * time.Minute
		cluster.Port = c.Port
		// NOTE: Workaround for https://github.com/gocql/gocql/issues/575#issuecomment-172124342
		cluster.IgnorePeerAddr = true
		cluster.DisableInitialHostLookup = true
		cluster.Events.DisableTopologyEvents = true
		cluster.Events.DisableNodeStatusEvents = true
		cluster.Events.DisableSchemaEvents = true

		session, err := igocqlx.WrapSession(cluster.CreateSession())
		if err != nil {
			panic(fmt.Errorf("db: failed to connect to test database, %v", err))
		}

		session.SetMapper(CQLMapper)
		c.Session = session
	}
}

// WithMockSession sets the mock session.
func WithMockSession(session *gocqlxmock.SessionxMock) ConfigOption {
	return func(c *Config) {
		c.Session = &MockConfig{session}
	}
}

// WithMigrations runs the migrations after the session has been initialized.
func WithMigrations() ConfigOption {
	return func(c *Config) {
		dir, err := iofs.New(src, "migrations")
		if err != nil {
			slog.Error("db: failed to initialize migrations ...", "error", err)
			return
		}

		logger := &MigrationLogger{}
		config := &cassandra.Config{KeyspaceName: c.Keyspace, MultiStatementEnabled: true}

		driver, err := cassandra.WithInstance(c.Session.Session().S.Session, config)
		if err != nil {
			slog.Error("db: failed to initialize driver for migrations ...", "error", err)
		}

		migrations, err := migrate.NewWithInstance("iofs", dir, "cassandra", driver)
		if err != nil {
			slog.Error("db: failed to initialize migrations ...", "error", err)
		}

		migrations.Log = logger

		version, dirty, err := migrations.Version()
		if err == migrate.ErrNilVersion {
			slog.Info("db: no migrations have been run, initializing keyspace ...")
		} else if err != nil {
			slog.Error("db: failed to get migration version ...", "error", err)
			return
		}

		slog.Info("db: migrations version ...", "version", version, "dirty", dirty)

		err = migrations.Up()
		if err != nil && err != migrate.ErrNoChange {
			slog.Warn("db: failed to run migrations ...", "error", err)
		}

		if err == migrate.ErrNoChange {
			slog.Info("db: no new migrations to run")
		}

		slog.Info("db: migrations done")
	}
}

// WithValidator registers a validator function.
func WithValidator(name string, validator validator.Func) ConfigOption {
	slog.Info("db: registering validator", "name", name)

	return func(c *Config) {
		if err := shared.Validator().RegisterValidation(name, validator); err != nil {
			panic(fmt.Errorf("db: failed to register validator %s, %v", name, err))
		}
	}
}

// NewSession creates a new session with the given options.
func NewSession(opts ...ConfigOption) *Config {
	db = &Config{}
	for _, opt := range opts {
		opt(db)
	}

	return db
}

// NewE2ESession creates a new session for end-to-end tests.
func NewE2ESession(port int, migrationPath string) {
	db = NewSession(
		WithHosts([]string{"localhost"}),
		WithKeyspace("ctrlplane_test"),
		WithPort(port),
		WithMigrationSourceURL(migrationPath),
		WithTimeout(10*time.Minute),
		WithE2ESession(),
		WithMigrations(),
		WithValidator("db_unique", UniqueField),
	)
}

// NewMockSession creates a new mock session.
func NewMockSession(session *gocqlxmock.SessionxMock) {
	db = NewSession(WithMockSession(session))
}

// DB returns the singleton database session.
func DB() *Config {
	if db == nil {
		slog.Info("db: creating new session")
		once.Do(func() {
			db = NewSession(
				FromEnvironment(),
				WithSessionCreation(),
				WithValidator("db_unique", UniqueField),
			)
		})
	}

	return db
}
