// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
//
// This software is made available by Breu, Inc., under the terms of the Breu  
// Community License Agreement, Version 1.0 located at  
// http://www.breu.io/breu-community-license/v1. BY INSTALLING, DOWNLOADING,  
// ACCESSING, USING OR DISTRIBUTING ANY OF THE SOFTWARE, YOU AGREE TO THE TERMS  
// OF SUCH LICENSE AGREEMENT. 

package entities

import (
	"time"

	itable "github.com/Guilospanck/igocqlx/table"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/table"
)

var (
	teamUserColumns = []string{
		"id",
		"user_id",
		"team_id",
		"created_at",
		"updated_at",
	}

	teamUserMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "team_users",
			Columns: teamUserColumns,
		},
	}

	teamUserTable = itable.New(*teamUserMeta.M)
)

type (
	TeamUser struct {
		ID        gocql.UUID `json:"id" cql:"id"`
		UserID    gocql.UUID `json:"user_id" cql:"user_id"`
		TeamID    gocql.UUID `json:"team_id" cql:"team_id"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
	}
)

func (tu *TeamUser) GetTable() itable.ITable { return teamUserTable }
func (tu *TeamUser) PreCreate() error        { return nil }
func (tu *TeamUser) PreUpdate() error        { return nil }
