package db

import (
	"testing"

	"go.breu.io/ctrlplane/internal/shared"
)

func TestEntityGetTable(expect string, entity Entity) func(*testing.T) {
	return func(t *testing.T) {
		if expect != entity.GetTable().Metadata().M.Name {
			t.Errorf("expected %s, got %s", expect, entity.GetTable().Metadata().M.Name)
		}
	}
}

func TestEntityOps(entity Entity, tests shared.TestFnMap) func(*testing.T) {
	return func(t *testing.T) {
		for name, test := range tests {
			t.Run(name, test.Run(test.Args, test.Want))
		}
	}
}
