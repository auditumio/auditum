// Copyright 2023 Igor Zibarev
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

package auditum

import (
	"fmt"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

const (
	storeTypeSQLite   = "sqlite"
	storeTypePostgres = "postgres"
)

type StoreConfig struct {
	Type     string         `yaml:"type" json:"type"`
	SQLite   SQLiteConfig   `yaml:"sqlite" json:"sqlite"`
	Postgres PostgresConfig `yaml:"postgres" json:"postgres"`
}

func (c StoreConfig) Validate() error {
	err := validation.ValidateStruct(&c,
		validation.Field(
			&c.Type,
			validation.Required,
			validation.In(storeTypeSQLite, storeTypePostgres),
		),
	)
	if err != nil {
		return err
	}

	switch c.Type {
	case storeTypeSQLite:
		if err := c.SQLite.Validate(); err != nil {
			return fmt.Errorf("invalid 'sqlite': %v", err)
		}
		return nil
	case storeTypePostgres:
		if err := c.Postgres.Validate(); err != nil {
			return fmt.Errorf("invalid 'postgres': %v", err)
		}
		return nil
	default:
		return fmt.Errorf("unknown 'type': %s", c.Type)
	}
}

type SQLiteConfig struct {
	DatabasePath   string `yaml:"databasePath" json:"databasePath"`
	MigrationsPath string `yaml:"migrationsPath" json:"migrationsPath"`
	LogQueries     bool   `yaml:"logQueries" json:"logQueries"`
}

func (c SQLiteConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.DatabasePath, validation.Required),
		validation.Field(&c.MigrationsPath, validation.Required),
	)
}

type PostgresConfig struct {
	Host           string `yaml:"host" json:"host"`
	Port           string `yaml:"port" json:"port"`
	Database       string `yaml:"database" json:"database"`
	Username       string `yaml:"username" json:"username"`
	Password       string `yaml:"password" json:"password"`
	SSLMode        string `yaml:"sslmode" json:"sslmode"`
	MigrationsPath string `yaml:"migrationsPath" json:"migrationsPath"`
	LogQueries     bool   `yaml:"logQueries" json:"logQueries"`
}

func (c PostgresConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Host, validation.Required, is.Host),
		validation.Field(&c.Port, validation.Required, is.Port),
		validation.Field(&c.Database, validation.Required),
		validation.Field(&c.Username, validation.Required),
		validation.Field(&c.Password, validation.Required),
		validation.Field(&c.SSLMode, validation.Required),
		validation.Field(&c.MigrationsPath, validation.Required),
	)
}

var defaultStoreConfig = StoreConfig{
	Type: storeTypeSQLite,
	SQLite: SQLiteConfig{
		DatabasePath:   ":memory:",
		MigrationsPath: "./internal/sql/sqlite/migrations",
		LogQueries:     false,
	},
	Postgres: PostgresConfig{
		Host:           "",
		Port:           "5432",
		Database:       "auditum_db",
		Username:       "",
		Password:       "",
		SSLMode:        "require",
		MigrationsPath: "./internal/sql/postgres/migrations",
		LogQueries:     false,
	},
}
