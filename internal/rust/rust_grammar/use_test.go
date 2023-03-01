package rust_grammar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUse(t *testing.T) {
	tests := []struct {
		Name        string
		ExpectedUse []FlattenUse
	}{
		{
			Name: "use something::something;",
			ExpectedUse: []FlattenUse{{
				PathSlices: []string{"something"},
				Name:       Name{Original: "something"},
			}},
		},
		{
			Name: "use something::something as another;",
			ExpectedUse: []FlattenUse{{
				PathSlices: []string{"something"},
				Name:       Name{Original: "something", Alias: "another"},
			}},
		},
		{
			Name: "pub use something::something;",
			ExpectedUse: []FlattenUse{{
				Pub:        true,
				PathSlices: []string{"something"},
				Name:       Name{Original: "something"},
			}},
		},
		{
			Name: "pub (   crate) use something  ::something  ;",
			ExpectedUse: []FlattenUse{{
				Pub:        true,
				PathSlices: []string{"something"},
				Name:       Name{Original: "something"},
			}},
		},
		{
			Name: "use something::{something};",
			ExpectedUse: []FlattenUse{{
				PathSlices: []string{"something"},
				Name:       Name{Original: "something"},
			}},
		},
		{
			Name: "use something::{one, OrAnother};",
			ExpectedUse: []FlattenUse{
				{
					PathSlices: []string{"something"},
					Name:       Name{Original: "one"},
				},
				{
					PathSlices: []string{"something"},
					Name:       Name{Original: "OrAnother"},
				},
			},
		},
		{
			Name: "use something::{one as two, OrAnother};",
			ExpectedUse: []FlattenUse{
				{
					PathSlices: []string{"something"},
					Name:       Name{Original: "one", Alias: "two"},
				},
				{
					PathSlices: []string{"something"},
					Name:       Name{Original: "OrAnother"},
				},
			},
		},
		{
			Name: "use one::very_long::veeery_long::path::something;",
			ExpectedUse: []FlattenUse{{
				PathSlices: []string{"one", "very_long", "veeery_long", "path"},
				Name:       Name{Original: "something"},
			}},
		},
		{
			Name: "use one::very_long::veeery_long::path::*;",
			ExpectedUse: []FlattenUse{{
				PathSlices: []string{"one", "very_long", "veeery_long", "path"},
				All:        true,
			}},
		},
		{
			Name: "pub use crate::ast::{Node, Operator};\n",
			ExpectedUse: []FlattenUse{
				{
					Pub:        true,
					PathSlices: []string{"crate", "ast"},
					Name:       Name{Original: "Node"},
				},
				{
					Pub:        true,
					PathSlices: []string{"crate", "ast"},
					Name:       Name{Original: "Operator"},
				},
			},
		},
		{
			Name: "use crate::{one::One, another::*};",
			ExpectedUse: []FlattenUse{
				{
					PathSlices: []string{"crate", "one"},
					Name:       Name{Original: "One"},
				},
				{
					PathSlices: []string{"crate", "another"},
					All:        true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)

			content := []byte(tt.Name)
			parsed, err := parser.ParseBytes("", content)
			a.NoError(err)

			var uses []FlattenUse
			for _, stmt := range parsed.Statements {
				if stmt.Use != nil {
					uses = append(uses, stmt.Use.Flatten()...)
				}
			}

			a.Equal(tt.ExpectedUse, uses)
		})
	}
}
