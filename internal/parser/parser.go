package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/isacikgoz/mattermost-suite-utilities/internal/model"
)

const typeName = "MyStruct"

// ParseDirectory reads a directory and generates representation of client struct to be generated.
func ParseDirectory(path string) (*model.Struct, error) {
	packages, err := parser.ParseDir(token.NewFileSet(), filepath.Join(path, "app"), nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("could not parse packages %q: %w", path, err)
	}

	// TODO: make name and abbrv constant
	str := &model.Struct{
		Name:  "Client",
		Abbrv: "c",
	}

	for name, pkg := range packages {
		if name != "app" {
			continue
		}

		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				fn, ok := n.(*ast.FuncDecl)
				if !ok {
					return true
				}
				rc := fn.Recv
				if rc == nil || len(rc.List) == 0 {
					return true
				}
				switch rci := rc.List[0].Type.(type) {
				case *ast.Ident:
					if rci.Name != typeName {
						return false
					}
				case *ast.StarExpr:
					rcn, ok := rci.X.(*ast.Ident)
					if !ok {
						return false
					}
					if rcn.Name != typeName {
						return false
					}
				}

				args := make([]*model.Argument, 0)
				for _, field := range fn.Type.Params.List {
					for _, fieldName := range field.Names {
						argType, isPointer := exprName(field.Type)
						args = append(args, &model.Argument{
							Name:    fieldName.Name,
							Type:    argType,
							Pointer: isPointer,
						})
					}
				}

				// skip the args with 0 param (this shouldn't be a case btw.)
				if len(args) == 0 {
					return false
				}

				// let n be the number of args, n-1 is the actual return type and args[:n-1] will be the actual args
				// we need to unwrap the return type. To do so, we need to find the type spec of it
				// unwrap the last param
				out := fn.Type.Params.List[len(fn.Type.Params.List)-1]
				star, ok := out.Type.(*ast.StarExpr)
				if !ok {
					return false
				}
				sel, ok := star.X.(*ast.SelectorExpr)
				if !ok {
					return false
				}

				// TODO: extract this from types.Uses
				pkgPath := filepath.Join(path, "model")
				returnStr, err := parseModelStruct(pkgPath, "github.com/mattermost/mattermost-server/v5/model", sel.Sel.Name)
				if err != nil {
					fmt.Fprintf(os.Stderr, "could not parse model struct: %s\n", err)
					return false
				}

				// prepare the methods for code gen
				str.Methods = append(str.Methods, &model.Method{
					Abbrv:            str.Abbrv,
					Struct:           str.Name,
					RemoteStructName: typeName,
					RemoteMethodName: fn.Name.Name,
					Name:             strings.TrimSuffix(sel.Sel.Name, "RPCResponse"),
					Arguments:        args[:len(args)-1],
					OutputStruct:     args[len(args)-1],
					ReturnValues:     returnStr.Fields,
				})

				return false
			})
		}
	}

	return str, nil
}

func exprName(expr ast.Expr) (string, bool) {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name, false
	case *ast.StarExpr:
		return exprName(t.X)
	case *ast.SelectorExpr:
		pkgName, _ := exprName(t.X)
		selName, _ := exprName(t.Sel)
		return strings.Join([]string{pkgName, selName}, "."), true
	default:
		return "", false
	}
}

func parseModelStruct(pkgPath, importSpec, typeName string) (*model.Struct, error) {
	packages, err := parser.ParseDir(token.NewFileSet(), pkgPath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("could not parse packages %q: %w", pkgPath, err)
	}

	str := &model.Struct{
		Name:  typeName,
		Abbrv: "c",
	}

	for name, pkg := range packages {
		if name != "model" {
			continue
		}

		str.Fields = make([]*model.Field, 0)
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				gd, ok := n.(*ast.GenDecl)
				if !ok {
					return true
				}

				// does it have eny spec?
				if len(gd.Specs) == 0 {
					return false
				}

				// if not type spec, move on from this node
				ts, ok := gd.Specs[0].(*ast.TypeSpec)
				if !ok || ts.Name.Name != typeName {
					return false
				}

				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					return false
				}

				for _, field := range st.Fields.List {
					for _, name := range field.Names {
						typeName, isPointer := exprName(field.Type)
						str.Fields = append(str.Fields, &model.Field{
							Import:  importSpec,
							Name:    name.Name,
							Type:    strings.Join([]string{"model", typeName}, "."),
							Pointer: isPointer,
						})
					}

				}

				return false
			})
		}
	}

	return str, nil
}
