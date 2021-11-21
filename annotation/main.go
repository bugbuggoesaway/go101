package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go101/annotation/plugins"
	"strings"
)

type Annotation struct {
	Type string `json:"type"`
}

type Annotations []*Annotation

func (as Annotations) String() string {
	builder := strings.Builder{}
	builder.WriteString("[")
	for i, annotation := range as {
		if i != 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%+v", annotation))
	}
	builder.WriteString("]")
	return builder.String()
}

func main() {
	fileSet := token.NewFileSet()
	astFile, err := parser.ParseFile(fileSet, "annotation/service/user_service.go", nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("ParseFile. err=[%v]\n", err)
		return
	}
	ast.Inspect(astFile, func(node ast.Node) bool {
		if funcDecl, ok := node.(*ast.FuncDecl); ok {
			if annotations := getAnnotations(funcDecl); len(annotations) != 0 {
				fmt.Printf("node=[%+v], annotations=[%+v]\n", node, annotations)
				for _, annotation := range annotations {
					if plugin, exist := plugins.Get(annotation.Type); exist {
						if err := plugin.Enhance(fileSet, astFile, []ast.Node{node}); err != nil {
							fmt.Printf("Enhance. pluginType=[%s], err=[%v]\n", annotation.Type, err)
						}
					} else {
						fmt.Printf("Plugin not found. pluginType=[%s]\n", annotation.Type)
					}
				}

				buffer := bytes.NewBufferString("")
				if err := format.Node(buffer, fileSet, astFile); err != nil {
					fmt.Printf("format.Node. err=[%v]\n", err)
				} else {
					fmt.Printf("generated file:\n%s", buffer.String())
				}
			}
		}
		return true
	})
}

func getAnnotations(decl *ast.FuncDecl) Annotations {
	if decl != nil {
		if doc := decl.Doc; doc != nil {
			annotations := make([]*Annotation, 0)
			for _, comment := range doc.List {
				if strings.HasPrefix(comment.Text, "//@") {
					annotations = append(annotations, &Annotation{Type: comment.Text[len("//@"):]})
				}
			}
			return annotations
		}
	}
	return nil
}
