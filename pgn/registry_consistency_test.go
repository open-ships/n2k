package pgn

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

type registrySourceFacts struct {
	decoders   map[string]string
	encoders   map[string]string
	pgnMethods map[string]uint32
}

func TestTypedPgnListMatchesCodeImplementation(t *testing.T) {
	facts := collectRegistrySourceFacts(t)

	registeredIDs := make(map[string]struct{}, len(typedPgnList))
	registeredDecoders := make(map[string]struct{}, len(typedPgnList))
	registeredEncoders := make(map[string]struct{}, len(typedPgnList))

	for _, info := range typedPgnList {
		if info.Id == "" {
			t.Fatalf("typedPgnList contains entry with empty Id for PGN %d", info.PGN)
		}
		if _, exists := registeredIDs[info.Id]; exists {
			t.Fatalf("typedPgnList contains duplicate Id %q", info.Id)
		}
		registeredIDs[info.Id] = struct{}{}

		if info.Decoder == nil {
			t.Errorf("%s has nil Decoder", info.Id)
		} else {
			decoderName := registryFuncName(info.Decoder)
			registeredDecoders[decoderName] = struct{}{}
			if decoderName != "Decode"+info.Id {
				t.Errorf("%s decoder = %s, want Decode%s", info.Id, decoderName, info.Id)
			}
			if _, exists := facts.decoders[decoderName]; !exists {
				t.Errorf("%s references decoder %s, but no matching function exists", info.Id, decoderName)
			}
		}

		if info.Encoder == nil {
			t.Errorf("%s has nil Encoder", info.Id)
		} else {
			encoderName := registryFuncName(info.Encoder)
			registeredEncoders[encoderName] = struct{}{}
			if _, exists := facts.encoders[encoderName]; !exists {
				t.Errorf("%s references encoder %s, but no matching function exists", info.Id, encoderName)
			}
		}

		methodPGN, exists := facts.pgnMethods[info.Id]
		if !exists {
			t.Errorf("%s has no PGNNumber method", info.Id)
		} else if methodPGN != info.PGN {
			t.Errorf("%s PGNNumber() = %d, typedPgnList PGN = %d", info.Id, methodPGN, info.PGN)
		}
	}

	for decoder := range facts.decoders {
		if _, exists := registeredDecoders[decoder]; !exists {
			t.Errorf("decoder %s exists in code but is missing from typedPgnList", decoder)
		}
	}
	for encoder := range facts.encoders {
		if _, exists := registeredEncoders[encoder]; !exists {
			t.Errorf("encoder %s exists in code but is missing from typedPgnList", encoder)
		}
	}
	for structName := range facts.pgnMethods {
		if _, exists := registeredIDs[structName]; !exists {
			t.Errorf("%s has PGNNumber() but is missing from typedPgnList", structName)
		}
	}
}

func TestCatalogOnlyListHasNoTypedImplementation(t *testing.T) {
	for _, info := range catalogOnlyList {
		if info.Decoder != nil {
			t.Errorf("catalog-only entry %s has Decoder set", info.Id)
		}
		if info.Encoder != nil {
			t.Errorf("catalog-only entry %s has Encoder set", info.Id)
		}
	}
}

func collectRegistrySourceFacts(t *testing.T) registrySourceFacts {
	t.Helper()

	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("glob pgn source files: %v", err)
	}

	facts := registrySourceFacts{
		decoders:   make(map[string]string),
		encoders:   make(map[string]string),
		pgnMethods: make(map[string]uint32),
	}

	fset := token.NewFileSet()
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") || file == "registry.go" {
			continue
		}

		parsed, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", file, err)
		}

		for _, decl := range parsed.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			name := fn.Name.Name
			switch {
			case fn.Recv == nil && strings.HasPrefix(name, "Decode"):
				facts.decoders[name] = file
			case fn.Recv == nil && strings.HasPrefix(name, "encode") && strings.HasSuffix(name, "Msg"):
				facts.encoders[name] = file
			case name == "PGNNumber" && fn.Recv != nil:
				structName, ok := receiverStructName(fn)
				if !ok || structName == "UnknownPGN" {
					continue
				}
				pgn, ok := returnedUint32Literal(fn)
				if !ok {
					t.Fatalf("%s PGNNumber() must return a uint32 literal", structName)
				}
				facts.pgnMethods[structName] = pgn
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

func registryFuncName(fn any) string {
	if fn == nil {
		return ""
	}
	fullName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	if idx := strings.LastIndex(fullName, "."); idx >= 0 {
		return fullName[idx+1:]
	}
	return fullName
}
