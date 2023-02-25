package js

import (
	"context"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"dep-tree/internal/language"
)

const importsTestFolder = ".imports_test"

func TestParser_parseImports_IsCached(t *testing.T) {
	a := require.New(t)
	ctx := context.Background()
	file := path.Join(importsTestFolder, "index.ts")
	lang, err := MakeJsLanguage(file)
	a.NoError(err)

	start := time.Now()
	ctx, _, err = lang.ParseImports(ctx, file)
	a.NoError(err)
	nonCached := time.Since(start)

	start = time.Now()
	_, _, err = lang.ParseImports(ctx, file)
	a.NoError(err)
	cached := time.Since(start)

	ratio := nonCached.Nanoseconds() / cached.Nanoseconds()
	a.Greater(ratio, int64(10))
}

func TestParser_parseImports(t *testing.T) {
	wd, _ := os.Getwd()

	tests := []struct {
		Name           string
		File           string
		Expected       map[string]language.ImportEntry
		ExpectedErrors []string
	}{
		{
			Name: "test 1",
			File: path.Join(importsTestFolder, "index.ts"),
			Expected: map[string]language.ImportEntry{
				path.Join(wd, importsTestFolder, "2", "2.ts"):      {Names: []string{"a", "b"}},
				path.Join(wd, importsTestFolder, "2", "index.ts"):  {All: true},
				path.Join(wd, importsTestFolder, "1", "a", "a.ts"): {All: true},
			},
			ExpectedErrors: []string{
				"could not perform relative import for './unexisting'",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			a := require.New(t)
			lang, err := MakeJsLanguage(tt.File)
			a.NoError(err)

			_, results, err := lang.ParseImports(context.Background(), tt.File)
			a.NoError(err)
			for expectedPath, expectedNames := range tt.Expected {
				resultNames, ok := results.Imports.Get(expectedPath)
				a.Equal(true, ok)
				a.Equal(expectedNames, resultNames)
			}

			a.Equal(len(tt.ExpectedErrors), len(results.Errors))
			if results.Errors != nil {
				for i, err := range results.Errors {
					a.ErrorContains(err, tt.ExpectedErrors[i])
				}
			}
		})
	}
}
