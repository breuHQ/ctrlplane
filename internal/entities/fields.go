// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
//
// This software is made available by Breu, Inc., under the terms of the Breu  
// Community License Agreement, Version 1.0 located at  
// http://www.breu.io/breu-community-license/v1. BY INSTALLING, DOWNLOADING,  
// ACCESSING, USING OR DISTRIBUTING ANY OF THE SOFTWARE, YOU AGREE TO THE TERMS  
// OF SUCH LICENSE AGREEMENT. 

package entities

import (
	"encoding/json"

	"github.com/gocql/gocql"
)

type (
	AppConfig struct{}

	BluePrintRegions struct {
		GCP     []string `json:"gcp"`
		AWS     []string `json:"aws"`
		Azure   []string `json:"azure"`
		Default string   `json:"default"`
	}

	RolloutArtifact struct{}

	RolloutArtifacts map[string]RolloutArtifact
)

func (config AppConfig) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(config)
}

func (config *AppConfig) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, config)
}

func (regions BluePrintRegions) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(regions)
}

func (regions *BluePrintRegions) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, regions)
}

func (artifacts RolloutArtifacts) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(artifacts)
}

func (artifacts *RolloutArtifacts) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return json.Unmarshal(data, artifacts)
}
