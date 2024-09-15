package gogoast_test

import (
	"bytes"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strings"
	"testing"
)

func TestAst(t *testing.T) {
	testCases := map[string]struct {
		src []string
		out []string
	}{
		"just a package": {
			src: []string{
				`package test`,
			},
			out: []string{
				`File`,
				`  Ident test`,
			},
		},
		"empty func, no args, no returns": {
			src: []string{
				`package test`,
				`func f() {`,
				`}`,
			},
			out: []string{
				`File`,
				`  Ident test`,
				`  FuncDecl`,
				`    Ident f`,
				`    FuncType`,
				`      FieldList`,
				`    BlockStmt`,
			},
		},
		"empty func, two int/string args, two bool/error returns": {
			src: []string{
				`package test`,
				`func f(i int, s string) (bool, error) {`,
				`}`,
			},
			out: []string{
				`File`,
				`  Ident test`,
				`  FuncDecl`,
				`    Ident f`,
				`    FuncType`,
				`      FieldList`,
				`        Field`,
				`          Ident i`,
				`          Ident int`,
				`        Field`,
				`          Ident s`,
				`          Ident string`,
				`      FieldList`,
				`        Field`,
				`          Ident bool`,
				`        Field`,
				`          Ident error`,
				`    BlockStmt`,
			},
		},
		"func with a println call": {
			src: []string{
				`package test`,
				`import "fmt"`,
				`func f() {`,
				`	fmt.Println("hello")`,
				`}`,
			},
			out: []string{
				`File`,
				`  Ident test`,
				`  GenDecl`,
				`    ImportSpec`,
				`      BasicLit STRING "fmt"`,
				`  FuncDecl`,
				`    Ident f`,
				`    FuncType`,
				`      FieldList`,
				`    BlockStmt`,
				`      ExprStmt`,
				`        CallExpr`,
				`          SelectorExpr`,
				`            Ident fmt`,
				`            Ident Println`,
				`          BasicLit STRING "hello"`,
			},
		},
		"if true with 2 returns": {
			src: []string{
				`package test`,
				`func f() string {`,
				`	if true {`,
				`		return "a"`,
				`	} else {`,
				`		return "b"`,
				`	}`,
				`}`,
			},
			out: []string{
				`File`,
				`  Ident test`,
				`  FuncDecl`,
				`    Ident f`,
				`    FuncType`,
				`      FieldList`,
				`      FieldList`,
				`        Field`,
				`          Ident string`,
				`    BlockStmt`,
				`      IfStmt`,
				`        Ident true`,
				`        BlockStmt`,
				`          ReturnStmt`,
				`            BasicLit STRING "a"`,
				`        BlockStmt`,
				`          ReturnStmt`,
				`            BasicLit STRING "b"`,
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			buffer := bytes.Buffer{}
			_, err := io.Copy(&buffer, strings.NewReader(strings.Join(testCase.src, "\n")))
			fileSet := token.NewFileSet()
			parsedFile, err := parser.ParseFile(fileSet, "", buffer.String(), 0)
			if err != nil {
				t.Error(err.Error())
			}
			indentLevel := 0
			var out []string
			ast.Inspect(parsedFile, func(node ast.Node) bool {
				if node == nil {
					indentLevel--
					return true
				}
				stringID := "Undefined"
				switch tNode := node.(type) {
				case *ast.File:
					stringID = "File"
				case *ast.Ident:
					stringID = fmt.Sprintf("Ident %s", tNode.Name)
				case *ast.FuncDecl:
					stringID = "FuncDecl"
				case *ast.FuncType:
					stringID = "FuncType"
				case *ast.FieldList:
					stringID = "FieldList"
				case *ast.BlockStmt:
					stringID = "BlockStmt"
				case *ast.Field:
					stringID = "Field"
				case *ast.ExprStmt:
					stringID = "ExprStmt"
				case *ast.CallExpr:
					stringID = "CallExpr"
				case *ast.SelectorExpr:
					stringID = "SelectorExpr"
				case *ast.BasicLit:
					stringID = fmt.Sprintf("BasicLit %s %s", tNode.Kind.String(), tNode.Value)
				case *ast.GenDecl:
					stringID = "GenDecl"
				case *ast.ImportSpec:
					stringID = "ImportSpec"
				case *ast.IfStmt:
					stringID = "IfStmt"
				case *ast.ReturnStmt:
					stringID = "ReturnStmt"
				default:
					stringID = fmt.Sprintf("Unknown %#v", tNode)
				}
				out = append(out, fmt.Sprintf("%s%s", strings.Repeat("  ", indentLevel), stringID))
				indentLevel++
				return true
			})
			if diff := cmp.Diff(testCase.out, out); diff != "" {
				t.Errorf("missmatch mismatch (-want +got):\\n%s", diff)
			}
		})
	}
}
