/*
 *    Copyright 2018 INS Ecosystem
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type ContractInterface struct {
	Types    map[string]string
	Methods  map[string][]*ast.FuncDecl
	Contract string
}

var mode string
var outfile string

func main() {
	flag.StringVar(&mode, "mode", "wrapper", "Generation mode: <wrapper|helper>")
	flag.StringVar(&outfile, "o", "-", "output file")
	flag.Parse()

	var output io.WriteCloser
	if outfile == "-" {
		output = os.Stdout
	} else {
		var err error
		output, err = os.OpenFile(outfile, os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer output.Close()
	}

	for _, fn := range flag.Args() {
		w, err := generateForFile(fn)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(output, w)
		if err != nil {
			panic(err)
		}
	}
}

func generateForFile(fn string) (io.Reader, error) {
	fs := token.NewFileSet()

	F, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		return nil, errors.Wrap(err, "Can't open file "+fn)
	}
	defer F.Close()

	buff, err := ioutil.ReadAll(F)
	if err != nil {
		return nil, errors.Wrap(err, "Can't read file "+fn)
	}

	node, err := parser.ParseFile(fs, fn, buff, parser.ParseComments)
	if err != nil {
		return nil, errors.Wrapf(err, "Can't parse %s", fn)
	}
	if node.Name.Name != "main" {
		panic("Contract must be in main package")
	}
	b := bytes.Buffer{}
	ci := getMethods(node, buff)
	if mode == "wrapper" {
		b.WriteString("package " + node.Name.Name + "\n\n")
		b.WriteString(generateWrappers(ci) + "\n")
		b.WriteString(generateExports(ci) + "\n")
	}
	return &b, nil
}

func getMethods(F *ast.File, text []byte) *ContractInterface {
	ci := ContractInterface{
		Types:   make(map[string]string),
		Methods: make(map[string][]*ast.FuncDecl),
	}
	for _, d := range F.Decls {
		switch td := d.(type) {
		case *ast.GenDecl:
			if td.Tok != token.TYPE {
				continue
			}
			typeNode := td.Specs[0].(*ast.TypeSpec)
			if strings.Contains(td.Doc.Text(), "@inscontract") {
				if ci.Contract != "" {
					panic("more than one contract in a file")
				}
				ci.Contract = typeNode.Name.Name
				continue
			}
			ci.Types[typeNode.Name.Name] = string(text[typeNode.Pos()-1 : typeNode.End()])
			continue
		case *ast.FuncDecl:
			if td.Recv.NumFields() == 0 { // not a method
				continue
			}
			r := td.Recv.List[0].Type
			if tr, ok := r.(*ast.StarExpr); ok { // *type
				r = tr.X
			}
			typename := r.(*ast.Ident).Name
			ci.Methods[typename] = append(ci.Methods[typename], td)
		}
	}
	return &ci
}

// nolint
func generateTypes(ci *ContractInterface) string {
	text := ""
	for _, t := range ci.Types {
		text += "type " + t + "\n"
	}

	text += "type " + ci.Contract + " struct { // Contract proxy type\n"
	text += "    address Reference logicrunner.Reference\n"
	text += "}\n\n"

	text += "func (c *" + ci.Contract + ")GetReference"
	// GetReference
	return text
}

func generateWrappers(ci *ContractInterface) string {
	text := `import (
	"github.com/insolar/insolar/logicrunner/goplugin/testplugins/foundation"
	)` + "\n"

	for _, method := range ci.Methods[ci.Contract] {
		text += generateMethodWrapper(method, ci.Contract) + "\n"
	}
	return text
}

func generateMethodWrapper(method *ast.FuncDecl, class string) string {
	text := fmt.Sprintf("func (self *%s) INSWRAPER_%s(cbor foundation.CBORMarshaler, data []byte) ([]byte) {\n",
		class, method.Name.Name)
	text += fmt.Sprintf("\targs := [%d]interface{}{}\n", method.Type.Params.NumFields())

	args := []string{}
	for i, arg := range method.Type.Params.List {
		initializer := ""
		tname := fmt.Sprintf("%v", arg.Type)
		switch tname {
		case "uint", "int", "int8", "uint8", "int32", "uint32", "int64", "uint64":
			initializer = tname + "(0)"
		case "string":
			initializer = `""`
		default:
			initializer = tname + "{}"
		}
		text += fmt.Sprintf("\targs[%d] = %s\n", i, initializer)
		args = append(args, fmt.Sprintf("args[%d].(%s)", i, tname))
	}

	text += "\tcbor.Unmarshal(&args, data)\n"

	rets := []string{}
	for i := range method.Type.Results.List {
		rets = append(rets, fmt.Sprintf("ret%d", i))
	}
	ret := strings.Join(rets, ", ")
	text += fmt.Sprintf("\t%s := self.%s(%s)\n", ret, method.Name.Name, strings.Join(args, ", "))

	text += fmt.Sprintf("\treturn cbor.Marshal([]interface{}{%s})\n", strings.Join(rets, ", "))
	text += "}\n"
	return text
}

/* generated snipped must be something like this

func (hw *HelloWorlder) INSWRAPER_Echo(cbor cborer, data []byte) ([]byte, error) {
	args := [1]interface{}{}
	args[0] = ""
	cbor.Unmarshal(&args, data)
	ret1, ret2 := hw.Echo(args[0].(string))
	return cbor.Marshal([]interface{}{ret1, ret2}), nil
}
*/

func generateExports(ci *ContractInterface) string {
	text := "var INSEXPORT " + ci.Contract + "\n"
	return text
}
