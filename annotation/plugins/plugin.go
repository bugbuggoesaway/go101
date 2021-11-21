package plugins

import (
	"fmt"
	"go/ast"
	"go/token"
)

var (
	plugins = make(map[string]Plugin)
)

type Plugin interface {
	Enhance(fileSet *token.FileSet, astFile *ast.File, nodes []ast.Node) error
}

func Register(pluginType string, plugin Plugin) {
	plugins[pluginType] = plugin
}

func Get(pluginType string) (Plugin, bool) {
	plugin, exist := plugins[pluginType]
	return plugin, exist
}

func Imports(astFile *ast.File, packages ...*packageInfo) {
	var importDecl *ast.GenDecl
	for _, decl := range astFile.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			if genDecl.Tok == token.IMPORT {
				importDecl = genDecl
				break
			}
		}
	}
	if importDecl == nil {
		importDecl = &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: nil,
		}
		astFile.Decls = append(astFile.Decls, importDecl)
	}

	for _, pkg := range packages {
		pkgName := fmt.Sprintf("\"%s\"", pkg.PkgName)
		importDecl.Specs = append(importDecl.Specs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: pkgName,
			},
			Name: &ast.Ident{Name: pkg.Alias},
		})
	}
}

func init() {
	Register("Log", &LogPlugin{})
	Register("Cache", &CachePlugin{})
}
