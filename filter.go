package n2k

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/cel-go/cel"
	celast "github.com/google/cel-go/common/ast"
	"github.com/open-ships/n2k/pgn"
)

// metaVars are the CEL variable names that correspond to message header metadata.
var metaVars = map[string]bool{
	"pgn":         true,
	"source":      true,
	"priority":    true,
	"destination": true,
}

// filter holds a compiled CEL expression that has been partitioned into
// pre-decode (metadata-only) and post-decode (struct field) stages.
type filter struct {
	preOnly  bool        // true if expression only references metadata
	hasPost  bool        // true if expression references msg.* fields
	preProg  cel.Program // evaluates metadata-only predicates (nil if none)
	postProg cel.Program // evaluates full expression including msg fields (nil if preOnly)
}

// newFullEnv creates a CEL environment with all variables (metadata + msg).
func newFullEnv() (*cel.Env, error) {
	return cel.NewEnv(
		cel.Variable("pgn", cel.IntType),
		cel.Variable("source", cel.IntType),
		cel.Variable("priority", cel.IntType),
		cel.Variable("destination", cel.IntType),
		cel.Variable("msg", cel.DynType),
	)
}

// newMetaEnv creates a CEL environment with only metadata variables.
func newMetaEnv() (*cel.Env, error) {
	return cel.NewEnv(
		cel.Variable("pgn", cel.IntType),
		cel.Variable("source", cel.IntType),
		cel.Variable("priority", cel.IntType),
		cel.Variable("destination", cel.IntType),
	)
}

// compileFilter compiles a CEL expression and partitions it into pre/post stages.
// The expression is split at top-level AND (&&) boundaries: conjuncts that only
// reference metadata variables become the pre-filter (evaluated before decoding),
// while conjuncts that reference msg.* become the post-filter.
func compileFilter(expr string) (*filter, error) {
	fullEnv, err := newFullEnv()
	if err != nil {
		return nil, fmt.Errorf("creating CEL environment: %w", err)
	}

	ast, iss := fullEnv.Compile(expr)
	if iss != nil && iss.Err() != nil {
		return nil, fmt.Errorf("compiling filter expression: %w", iss.Err())
	}

	// Collect top-level AND conjuncts.
	conjuncts := flattenAnd(ast.NativeRep().Expr())

	var metaExprs, msgExprs []string
	for _, c := range conjuncts {
		s, err := cel.ExprToString(c, ast.NativeRep().SourceInfo())
		if err != nil {
			return nil, fmt.Errorf("unparsing conjunct: %w", err)
		}
		if referencesMsg(c) {
			msgExprs = append(msgExprs, s)
		} else {
			metaExprs = append(metaExprs, s)
		}
	}

	f := &filter{}

	// Build pre-filter program from metadata-only conjuncts.
	if len(metaExprs) > 0 {
		metaEnv, err := newMetaEnv()
		if err != nil {
			return nil, fmt.Errorf("creating meta CEL environment: %w", err)
		}
		preExpr := strings.Join(metaExprs, " && ")
		preAst, iss := metaEnv.Compile(preExpr)
		if iss != nil && iss.Err() != nil {
			return nil, fmt.Errorf("compiling pre-filter: %w", iss.Err())
		}
		f.preProg, err = metaEnv.Program(preAst)
		if err != nil {
			return nil, fmt.Errorf("creating pre-filter program: %w", err)
		}
	}

	// Build post-filter program if any conjuncts reference msg.
	if len(msgExprs) > 0 {
		f.hasPost = true
		// Post-filter evaluates the full expression (meta + msg conjuncts) so that
		// short-circuit evaluation still works correctly.
		postProg, err := fullEnv.Program(ast)
		if err != nil {
			return nil, fmt.Errorf("creating post-filter program: %w", err)
		}
		f.postProg = postProg
	} else {
		// Expression is metadata-only.
		f.preOnly = true
	}

	return f, nil
}

// flattenAnd recursively collects conjuncts from top-level && operators.
func flattenAnd(e celast.Expr) []celast.Expr {
	if e.Kind() == celast.CallKind {
		call := e.AsCall()
		if call.FunctionName() == "_&&_" && !call.IsMemberFunction() {
			args := call.Args()
			if len(args) == 2 {
				var result []celast.Expr
				result = append(result, flattenAnd(args[0])...)
				result = append(result, flattenAnd(args[1])...)
				return result
			}
		}
	}
	return []celast.Expr{e}
}

// referencesMsg returns true if the expression references the "msg" variable.
func referencesMsg(e celast.Expr) bool {
	switch e.Kind() {
	case celast.IdentKind:
		return e.AsIdent() == "msg"
	case celast.SelectKind:
		return referencesMsg(e.AsSelect().Operand())
	case celast.CallKind:
		call := e.AsCall()
		if call.IsMemberFunction() {
			if referencesMsg(call.Target()) {
				return true
			}
		}
		for _, arg := range call.Args() {
			if referencesMsg(arg) {
				return true
			}
		}
	case celast.ListKind:
		list := e.AsList()
		for i := 0; i < list.Size(); i++ {
			if referencesMsg(list.Elements()[i]) {
				return true
			}
		}
	case celast.MapKind:
		m := e.AsMap()
		for i := 0; i < m.Size(); i++ {
			entry := m.Entries()[i].AsMapEntry()
			if referencesMsg(entry.Key()) || referencesMsg(entry.Value()) {
				return true
			}
		}
	case celast.ComprehensionKind:
		comp := e.AsComprehension()
		return referencesMsg(comp.IterRange()) ||
			referencesMsg(comp.AccuInit()) ||
			referencesMsg(comp.LoopCondition()) ||
			referencesMsg(comp.LoopStep()) ||
			referencesMsg(comp.Result())
	}
	return false
}

// evalPre evaluates the pre-filter against message metadata.
// Returns true if the message should proceed to decoding.
func (f *filter) evalPre(info pgn.MessageInfo) bool {
	if f.preProg == nil {
		return true
	}
	activation := map[string]any{
		"pgn":         int64(info.PGN),
		"source":      int64(info.SourceId),
		"priority":    int64(info.Priority),
		"destination": int64(info.TargetId),
	}
	out, _, err := f.preProg.Eval(activation)
	if err != nil {
		return false
	}
	b, ok := out.Value().(bool)
	return ok && b
}

// evalPostWithInfo evaluates the post-filter against decoded struct fields
// combined with message metadata. Returns true if the message passes.
func (f *filter) evalPostWithInfo(info pgn.MessageInfo, fields map[string]any) bool {
	if f.postProg == nil {
		return true
	}
	activation := map[string]any{
		"pgn":         int64(info.PGN),
		"source":      int64(info.SourceId),
		"priority":    int64(info.Priority),
		"destination": int64(info.TargetId),
		"msg":         fields,
	}
	out, _, err := f.postProg.Eval(activation)
	if err != nil {
		return false
	}
	b, ok := out.Value().(bool)
	return ok && b
}

// structToFilterMap converts a decoded PGN struct to a map[string]any for CEL evaluation.
// Both the original Go field name and a lowercase version are included as keys.
// The embedded Info field is skipped. Nil pointers are skipped. Numeric types are
// converted to int64 or float64 for CEL compatibility.
func structToFilterMap(v any) map[string]any {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil
	}

	rt := rv.Type()
	m := make(map[string]any, rt.NumField()*2)

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		// Skip unexported fields.
		if !field.IsExported() {
			continue
		}
		// Skip the embedded Info field.
		if field.Name == "Info" && field.Type == reflect.TypeOf(pgn.MessageInfo{}) {
			continue
		}

		fv := rv.Field(i)

		// Dereference pointers; skip nil.
		if fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				continue
			}
			fv = fv.Elem()
		}

		val := convertNumeric(fv)

		m[field.Name] = val
		lower := strings.ToLower(field.Name)
		if lower != field.Name {
			m[lower] = val
		}
	}

	return m
}

// convertNumeric converts numeric reflect.Values to int64 or float64 for CEL.
func convertNumeric(v reflect.Value) any {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return v.Float()
	default:
		return v.Interface()
	}
}
