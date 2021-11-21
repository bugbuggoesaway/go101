package plugins

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"text/template"
)

type LogPlugin struct {
}

type funcInfo struct {
	FileName string   `json:"file_name"`
	FuncName string   `json:"func_name"`
	Request  string   `json:"request"`
	Response string   `json:"response"`
	Params   []string `json:"params"`
}

type packageInfo struct {
	PkgName string `json:"pkg_name"`
	Alias   string `json:"alias"`
}

func (p *LogPlugin) Enhance(fileSet *token.FileSet, astFile *ast.File, nodes []ast.Node) error {
	Imports(astFile, &packageInfo{
		PkgName: "encoding/json",
		Alias:   "_json",
	})

	for _, node := range nodes {
		if funcDecl, ok := node.(*ast.FuncDecl); ok {
			funcBody, err := p.buildFuncBody(funcDecl)
			if err != nil {
				fmt.Printf("buildFuncBody. err=[%v]\n", err)
				continue
			}
			funcInfo := &funcInfo{
				FileName: astFile.Name.Name,
				FuncName: funcDecl.Name.Name,
			}
			fmt.Printf("funcInfo=[%+v]\n", funcInfo)
			fmt.Printf("func:\n%s\n", funcBody)
			if expr, err := parser.ParseExpr(funcBody); err != nil {
				fmt.Printf("ParseExpr. err=[%v]\n", err)
			} else if callExpr, ok := expr.(*ast.CallExpr); ok {
				funcDecl.Body.List = append([]ast.Stmt{
					&ast.DeferStmt{
						Call: callExpr,
					},
				}, funcDecl.Body.List...)
			} else {
				fmt.Printf("Not CallExpr. expr=[%+v]\n", expr)
			}

		}
	}
	return nil
}

func (p *LogPlugin) buildFuncBody(funcDecl *ast.FuncDecl) (string, error) {
	logTemplate := `func() {
	paramMap := make(map[string]interface{})
	{{range .Params}}paramMap["{{.}}"] = {{.}}
	{{end}}
	params, _ := _json.Marshal(paramMap)
	fmt.Printf("%s.%s params=[%s]\n", "{{.FileName}}", "{{.FuncName}}", params)
}()`

	params := make([]string, len(funcDecl.Type.Params.List))
	for i, field := range funcDecl.Type.Params.List {
		params[i] = field.Names[0].Name
	}

	funcInfo := &funcInfo{
		FuncName: funcDecl.Name.Name,
		Params:   params,
	}
	funcBody := bytes.NewBufferString("")
	if template, err := template.New("").Parse(logTemplate); err != nil {
		fmt.Printf("Parse. err=[%v]\n", err)
		return "", err
	} else if err := template.Execute(funcBody, funcInfo); err != nil {
		fmt.Printf("Execute. err=[%v]\n", err)
		return "", err
	} else {
		fmt.Printf("func:\n%s\n", funcBody.String())
		return funcBody.String(), nil
	}
}
