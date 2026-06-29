package pgn

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

type pgnStructFacts struct {
	structs map[string]string
	methods map[string]map[string]struct{}
	pgns    map[string]uint32
}

func TestPGNMetadataMatchesGeneratedStructs(t *testing.T) {
	facts := collectPGNStructFacts(t)

	metadataIDs := make(map[string]uint32)
	for _, infos := range PgnInfoLookup {
		for _, info := range infos {
			if info.Id == "" {
				t.Fatalf("PgnInfoLookup contains entry with empty Id for PGN %d", info.PGN)
			}
			if existing, exists := metadataIDs[info.Id]; exists {
				t.Fatalf("PgnInfoLookup contains duplicate Id %q for PGNs %d and %d", info.Id, existing, info.PGN)
			}
			metadataIDs[info.Id] = info.PGN

			if _, exists := facts.structs[info.Id]; !exists {
				t.Errorf("%s is present in PgnInfoLookup but has no PGN struct", info.Id)
				continue
			}
			methodPGN, exists := facts.pgns[info.Id]
			if !exists {
				t.Errorf("%s has no PGNNumber method", info.Id)
			} else if methodPGN != info.PGN {
				t.Errorf("%s PGNNumber() = %d, PgnInfoLookup PGN = %d", info.Id, methodPGN, info.PGN)
			}
		}
	}

	requiredMethods := []string{"PGNNumber", "MessageInfo", "SetMessageInfo", "DecodePayload", "EncodePayload"}
	for structName := range facts.pgns {
		pgn, exists := metadataIDs[structName]
		if !exists {
			t.Errorf("%s has a PGN struct but is missing from PgnInfoLookup", structName)
			continue
		}
		for _, method := range requiredMethods {
			if _, exists := facts.methods[structName][method]; !exists {
				t.Errorf("%s is missing %s method", structName, method)
			}
		}
		if facts.pgns[structName] != pgn {
			t.Errorf("%s PGN mismatch: method=%d metadata=%d", structName, facts.pgns[structName], pgn)
		}
	}
}

func collectPGNStructFacts(t *testing.T) pgnStructFacts {
	t.Helper()

	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("glob PGN source files: %v", err)
	}
	sort.Strings(files)

	facts := pgnStructFacts{
		structs: make(map[string]string),
		methods: make(map[string]map[string]struct{}),
		pgns:    make(map[string]uint32),
	}

	fset := token.NewFileSet()
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		parsed, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", file, err)
		}

		for _, decl := range parsed.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					if _, ok := typeSpec.Type.(*ast.StructType); !ok {
						continue
					}
					facts.structs[typeSpec.Name.Name] = file
				}
			case *ast.FuncDecl:
				if d.Recv == nil {
					continue
				}
				structName, ok := receiverStructName(d)
				if !ok {
					continue
				}
				if facts.methods[structName] == nil {
					facts.methods[structName] = make(map[string]struct{})
				}
				facts.methods[structName][d.Name.Name] = struct{}{}
				if d.Name.Name == "PGNNumber" {
					pgn, ok := returnedUint32Literal(d)
					if !ok {
						continue
					}
					facts.pgns[structName] = pgn
				}
			}
		}
	}

	return facts
}

func receiverStructName(fn *ast.FuncDecl) (string, bool) {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return "", false
	}
	switch recv := fn.Recv.List[0].Type.(type) {
	case *ast.StarExpr:
		if ident, ok := recv.X.(*ast.Ident); ok {
			return ident.Name, true
		}
	case *ast.Ident:
		return recv.Name, true
	}
	return "", false
}

func returnedUint32Literal(fn *ast.FuncDecl) (uint32, bool) {
	if fn.Body == nil || len(fn.Body.List) != 1 {
		return 0, false
	}
	ret, ok := fn.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 1 {
		return 0, false
	}
	lit, ok := ret.Results[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.INT {
		return 0, false
	}
	value, err := strconv.ParseUint(lit.Value, 0, 32)
	if err != nil {
		return 0, false
	}
	return uint32(value), true
}
