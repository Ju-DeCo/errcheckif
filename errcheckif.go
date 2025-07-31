// errcheckif.go
package errcheckif

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = `checks that errors returned from functions are checked

The errcheckif checker ensures that whenever a function call returns an error,
that error is checked in a subsequent if statement using "err != nil",
"errors.Is", or "errors.As".`

var Analyzer = &analysis.Analyzer{
	Name:     "errcheckif",
	Doc:      doc,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

var errorType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.AssignStmt)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		assignStmt, ok := node.(*ast.AssignStmt)
		if !ok {
			return
		}

		// We only care about `:=` and `=` that come from a function call
		if len(assignStmt.Rhs) != 1 {
			return
		}
		callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr)
		if !ok {
			return
		}

		errIdent := findReturnedError(pass, assignStmt, callExpr)
		if errIdent == nil {
			return
		}

		// Find the code block and the statement's position within it.
		path, _ := astutil.PathEnclosingInterval(findFile(pass, assignStmt), assignStmt.Pos(), assignStmt.End())
		if path == nil {
			return
		}

		if !isErrorHandled(pass, errIdent, path) {
			pass.Reportf(errIdent.Pos(), "error '%s' is not checked with '!= nil', 'errors.Is', or 'errors.As'", errIdent.Name)
		}
	})

	return nil, nil
}

// findReturnedError finds the identifier for a returned error variable in an assignment.
func findReturnedError(pass *analysis.Pass, assign *ast.AssignStmt, call *ast.CallExpr) *ast.Ident {
	sig, ok := pass.TypesInfo.TypeOf(call.Fun).(*types.Signature)
	if !ok {
		return nil
	}

	results := sig.Results()
	if results.Len() == 0 {
		return nil
	}

	for i := 0; i < results.Len(); i++ {
		if types.Implements(results.At(i).Type(), errorType) {
			if i < len(assign.Lhs) {
				if ident, ok := assign.Lhs[i].(*ast.Ident); ok && ident.Name != "_" {
					return ident
				}
			}
		}
	}
	return nil
}

// isErrorHandled checks if the error is handled in the statements following its assignment.
func isErrorHandled(pass *analysis.Pass, errIdent *ast.Ident, path []ast.Node) bool {
	// The path gives us the nodes from the specific `*ast.AssignStmt` up to the root.
	// We need to find the enclosing block and the index of our statement.
	for i, node := range path {
		if block, ok := node.(*ast.BlockStmt); ok {
			// Find which statement in the block is the parent of our assignment
			for stmtIdx, stmt := range block.List {
				// The direct child of the block is at path[i-1]
				if stmt == path[i-1] {
					// Now check all subsequent statements in this block
					for j := stmtIdx + 1; j < len(block.List); j++ {
						subsequentStmt := block.List[j]
						if checkIfStmtHandlesError(pass, subsequentStmt, errIdent) {
							return true
						}
						// If the error variable is reassigned before a check, we consider it unhandled.
						if isIdentifierReassigned(pass, subsequentStmt, errIdent) {
							return false
						}
					}
				}
			}
		}
	}
	return false
}

// checkIfStmtHandlesError checks if a single statement (e.g., an `if` statement) handles the error.
func checkIfStmtHandlesError(pass *analysis.Pass, stmt ast.Node, errIdent *ast.Ident) bool {
	ifStmt, ok := stmt.(*ast.IfStmt)
	if !ok {
		return false
	}
	return checkCondition(pass, ifStmt.Cond, errIdent)
}

// checkCondition recursively checks if an expression (like an `if` condition) is a valid error check.
func checkCondition(pass *analysis.Pass, cond ast.Expr, errIdent *ast.Ident) bool {
	switch c := cond.(type) {
	case *ast.BinaryExpr:
		if c.Op == token.LOR {
			return checkCondition(pass, c.X, errIdent) || checkCondition(pass, c.Y, errIdent)
		}
		if c.Op == token.NEQ {
			// Check for `err != nil` or `nil != err`
			if isIdent(pass, c.X, errIdent) && isNil(pass, c.Y) {
				return true
			}
			if isNil(pass, c.X) && isIdent(pass, c.Y, errIdent) {
				return true
			}
		}

	case *ast.CallExpr:
		// Check for `errors.Is(err, ...)` or `errors.As(err, ...)`
		sel, ok := c.Fun.(*ast.SelectorExpr)
		if !ok {
			return false
		}
		if pkgIdent, ok := sel.X.(*ast.Ident); !ok || pkgIdent.Name != "errors" {
			return false
		}
		if sel.Sel.Name != "Is" && sel.Sel.Name != "As" {
			return false
		}
		if len(c.Args) > 0 && isIdent(pass, c.Args[0], errIdent) {
			return true
		}
	}
	return false
}

// isIdentifierReassigned checks if an identifier is reassigned within a given statement node.
// This implementation is safe against nil pointer panics.
func isIdentifierReassigned(pass *analysis.Pass, stmt ast.Node, errIdent *ast.Ident) bool {
	targetObj := pass.TypesInfo.ObjectOf(errIdent)
	if targetObj == nil {
		return false // Should not happen for a valid errIdent, but a good safeguard.
	}

	reassigned := false
	ast.Inspect(stmt, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true // Not an assignment, continue walking the tree.
		}

		for _, lhs := range assign.Lhs {
			ident, ok := lhs.(*ast.Ident)
			if !ok {
				continue // LHS is not a simple identifier, e.g., s.field
			}

			// Safely get the object for the LHS identifier and compare with our target.
			if pass.TypesInfo.ObjectOf(ident) == targetObj {
				reassigned = true
				return false // Found a reassignment, stop the inspection.
			}
		}
		return true
	})
	return reassigned
}

// --- Helper Functions ---

func isIdent(pass *analysis.Pass, expr ast.Expr, targetIdent *ast.Ident) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && pass.TypesInfo.ObjectOf(ident) == pass.TypesInfo.ObjectOf(targetIdent)
}

func isNil(pass *analysis.Pass, expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && pass.TypesInfo.ObjectOf(ident) == types.Universe.Lookup("nil")
}

func findFile(pass *analysis.Pass, node ast.Node) *ast.File {
	for _, file := range pass.Files {
		if file.Pos() <= node.Pos() && node.End() <= file.End() {
			return file
		}
	}
	return nil
}
