package utils

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/henrylee2cn/aster/aster"
	"github.com/rs/xid"
)

// StatementCollection is a collection of all statements in target function as stacks
type StatementCollection struct {
	AssignStack     []ast.Stmt
	ExprStack       []ast.Stmt
	IfStack         []ast.Stmt
	BadStack        []ast.Stmt
	DeclStack       []ast.Stmt
	EmptyStack      []ast.Stmt
	LabeledStack    []ast.Stmt
	SendStack       []ast.Stmt
	IncDecStack     []ast.Stmt
	GoStack         []ast.Stmt
	DeferStack      []ast.Stmt
	ReturnStack     []ast.Stmt
	BranchStack     []ast.Stmt
	BlockStack      []ast.Stmt
	SwitchStack     []ast.Stmt
	TypeSwitchStack []ast.Stmt
	CommStack       []ast.Stmt
	SelectStack     []ast.Stmt
	ForStack        []ast.Stmt
	RangeStack      []ast.Stmt
	Listing         []ast.Stmt
}

// LoadDirs parses the source code of Go files under the directories and loads a new program.
func LoadDirs(dirs ...string) (*aster.Program, error) {
	p := aster.NewProgram()
	for _, dir := range dirs {
		if !filepath.IsAbs(dir) {
			dir, _ = filepath.Abs(dir)
		}
		err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			if err != nil || !f.IsDir() {
				return nil
			}
			p.Import(path)
			return nil
		})
		if err != nil {
			fmt.Println("Error: ", err)
			return nil, err
		}
	}
	return p.Load()
}

// Validate Validates given functions removes empty functions if found exits when no function given
// returns the clean list of functions
func Validate(functions []string) []string {
	out := []string{}
	fn := []string{}
	for i := range functions {
		fn = strings.Split(functions[i], ",")
		for j := range fn {
			if fn[j] == "" {
				continue
			}
			out = append(out, fn[j])
		}
	}
	if len(out) == 0 {
		fmt.Println("Error: no functions given.")
		os.Exit(-1)
	}
	return out
}

// UniqueID generates a unique identifier for variable assignments
func UniqueID() string {
	id := xid.New()
	return id.String()
}

// GetNodeType returns a the given node's type
func GetNodeType(node ast.Node) string {
	val := reflect.ValueOf(node).Elem()
	return val.Type().Name()
}

// FormatNode returs the node as a string
func FormatNode(node ast.Node) string {
	buf := new(bytes.Buffer)
	_ = format.Node(buf, token.NewFileSet(), node)
	return buf.String()
}

// FindFunctions returns a list of *ast.FuncDecl matching given functions in a given folder path
// use verbose flag for printing found functions files.
func FindFunctions(project string, functions []string, verbose bool) []*ast.FuncDecl {
	functions = Validate(functions)
	targetFuncs := []*ast.FuncDecl{}
	for i := range functions {
		fn := findFunction(project, functions[i], verbose)
		targetFuncs = append(targetFuncs, fn)
	}
	return targetFuncs
}

// findFunction finds a function in a given project folder, returns the node of the given function
func findFunction(project string, function string, verbose bool) *ast.FuncDecl {

	var out *ast.FuncDecl
	fset := token.NewFileSet()
	packages, _ := parser.ParseDir(fset, project, nil, parser.AllErrors)

	for i := range packages {

		for _, file := range packages[i].Files {
			ast.Inspect(file, func(n ast.Node) bool {
				fn, ok := n.(*ast.FuncDecl)
				if ok {
					if fn.Name.Name == function {
						if verbose {
							fmt.Printf("[+] Found function `%s` in `%s`..\n\n", fn.Name.Name, fn.Name.Name)
						}
						out = fn
					}
				}
				return true
			})
		}
	}
	return out
}

// ParseFunctions returns a list of collections of statments of the given functions if found.
// verbose prints found file names
// TODO: error handling
func ParseFunctions(project string, functions []string, verbose bool) []StatementCollection {
	targetFuncs := FindFunctions(project, functions, verbose)
	out := []StatementCollection{}
	for i := range targetFuncs {
		out = append(out, parseFunction(targetFuncs[i].Body.List))
	}
	return out

}

func parseFunction(stmts []ast.Stmt) StatementCollection {
	collection := StatementCollection{}
	var element ast.Stmt
	for i := range stmts {
		element = stmts[i]
		//subCollection := StatementCollection{}
		switch GetNodeType(element) {
		case "AssignStmt":
			collection.AssignStack = append(collection.AssignStack, element)
			break
		case "ExprStmt":
			collection.ExprStack = append(collection.ExprStack, element)
			break

		case "IfStmt":
			collection.IfStack = append(collection.IfStack, element)
			break

		case "BadStmt":
			collection.BadStack = append(collection.BadStack, element)
			break

		case "DeclStmt":
			collection.DeclStack = append(collection.DeclStack, element)
			break

		case "EmptyStmt":
			collection.EmptyStack = append(collection.EmptyStack, element)
			break

		case "LabeledStmt":
			collection.LabeledStack = append(collection.LabeledStack, element)
			break

		case "SendStmt":
			collection.SendStack = append(collection.SendStack, element)
			break

		case "IncDecStmt":
			collection.IncDecStack = append(collection.IncDecStack, element)
			break

		case "GoStmt":
			collection.GoStack = append(collection.GoStack, element)
			break

		case "DeferStmt":
			collection.DeferStack = append(collection.DeferStack, element)
			break

		case "ReturnStmt":
			collection.ReturnStack = append(collection.ReturnStack, element)
			break

		case "BranchStmt":
			collection.BranchStack = append(collection.BranchStack, element)
			break

		case "BlockStmt":
			collection.BlockStack = append(collection.BlockStack, element)
			break

		case "SwitchStmt":
			collection.SwitchStack = append(collection.SwitchStack, element)
			break

		case "TypeSwitchStmt":
			collection.TypeSwitchStack = append(collection.TypeSwitchStack, element)
			break

		case "CommClauseStmt":
			collection.CommStack = append(collection.CommStack, element)
			break

		case "SelectStmt":
			collection.SelectStack = append(collection.SelectStack, element)
			break

		case "ForStmt":
			collection.ForStack = append(collection.ForStack, element)
			break

		case "RangeStmt":
			collection.RangeStack = append(collection.RangeStack, element)
			break
		}
		collection.Listing = append(collection.Listing, element)

	}

	return collection
}

// ReturnAssignments returns the assignment statements as a string
func ReturnAssignments(collection StatementCollection) StatementCollection {

	for i := range collection.AssignStack {
		fmt.Printf(FormatNode(collection.AssignStack[i]) + "\n")
	}
	collection = remove(collection)
	return collection
}

// remove assignment statements from the collection listing
func remove(collection StatementCollection) StatementCollection {
	for i := 0; i < len(collection.Listing); i++ {
		for j := 0; j < len(collection.AssignStack); j++ {
			if collection.Listing[i] == collection.AssignStack[j] {
				collection.Listing = append(collection.Listing[:i], collection.Listing[i+1:]...)
			}
		}
	}
	return collection
}
