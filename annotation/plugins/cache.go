package plugins

import (
	"go/ast"
	"go/token"
)

type CachePlugin struct {
}

func (p *CachePlugin) Enhance(fileSet *token.FileSet, astFile *ast.File, nodes []ast.Node) error {
	return nil
}
