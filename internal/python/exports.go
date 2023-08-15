package python

import (
	"context"
	"path"

	"dep-tree/internal/language"
	"dep-tree/internal/python/python_grammar"
)

//nolint:gocyclo
func (l *Language) ParseExports(ctx context.Context, file *python_grammar.File) (context.Context, *language.ExportsResult, error) {
	var exports []language.ExportEntry
	var errors []error
	for _, stmt := range file.Statements {
		switch {
		case stmt == nil:
			continue
			// NOTE: it is very typical to do something like:
			// try:
			//   import foo
			// except:
			//   import foo.compat as foo
		case stmt.Import != nil: // && !stmt.Import.Indented:.
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{
					{
						Original: stmt.Import.Path[0],
						Alias:    stmt.Import.Alias,
					},
				},
				Id: file.Path,
			})
			// NOTE: it is very typical to do something like:
			// try:
			//   from foo import bar
			// except:
			//   from foo.compat import bar
		case stmt.FromImport != nil: // && !stmt.FromImport.Indented:.
			entry := language.ExportEntry{
				Names: make([]language.ExportName, len(stmt.FromImport.Names)),
				All:   stmt.FromImport.All,
				Id:    file.Path,
			}
			for i, name := range stmt.FromImport.Names {
				entry.Names[i] = language.ExportName{
					Original: name.Name,
					Alias:    name.Alias,
				}
			}
			var resolved *ResolveResult
			if len(stmt.FromImport.Relative) > 0 {
				var err error
				resolved, err = ResolveRelative(stmt.FromImport.Path, path.Dir(file.Path), len(stmt.FromImport.Relative)-1)
				if err != nil {
					errors = append(errors, err)
					continue
				}
			} else {
				resolved = l.ResolveAbsolute(stmt.FromImport.Path)
			}
			switch {
			case resolved == nil:
			case resolved.InitModule != nil:
			case resolved.Directory != nil:
				// nothing.
			case resolved.File != nil:
				entry.Id = resolved.File.Path
			}
			exports = append(exports, entry)

		case stmt.VariableUnpack != nil && !stmt.VariableUnpack.Indented:
			entry := language.ExportEntry{
				Names: make([]language.ExportName, len(stmt.VariableUnpack.Names)),
				Id:    file.Path,
			}
			for i, name := range stmt.VariableUnpack.Names {
				entry.Names[i] = language.ExportName{Original: name}
			}
			exports = append(exports, entry)
		case stmt.VariableAssign != nil && !stmt.VariableAssign.Indented:
			entry := language.ExportEntry{
				Names: make([]language.ExportName, len(stmt.VariableAssign.Names)),
				Id:    file.Path,
			}
			for i, name := range stmt.VariableAssign.Names {
				entry.Names[i] = language.ExportName{Original: name}
			}
			exports = append(exports, entry)
		case stmt.VariableTyping != nil && !stmt.VariableTyping.Indented:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{{Original: stmt.VariableTyping.Name}},
				Id:    file.Path,
			})
		case stmt.Function != nil && !stmt.Function.Indented:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{
					{
						Original: stmt.Function.Name,
					},
				},
				Id: file.Path,
			})
		case stmt.Class != nil && !stmt.Class.Indented:
			exports = append(exports, language.ExportEntry{
				Names: []language.ExportName{
					{
						Original: stmt.Class.Name,
					},
				},
				Id: file.Path,
			})
		}
	}
	return ctx, &language.ExportsResult{Exports: exports, Errors: errors}, nil
}
