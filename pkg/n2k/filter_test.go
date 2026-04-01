package n2k

import (
	"testing"

	"github.com/open-ships/n2k/pkg/pgn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterCompilation(t *testing.T) {
	f, err := compileFilter("pgn == 127250")
	require.NoError(t, err)
	assert.NotNil(t, f)
}

func TestFilterPreOnly(t *testing.T) {
	f, err := compileFilter("pgn == 127250 && source == 3")
	require.NoError(t, err)
	assert.True(t, f.preOnly, "expression with only metadata vars should be pre-only")
	assert.False(t, f.hasPost, "should not have post filter")
	assert.NotNil(t, f.preProg, "should have pre program")
	assert.Nil(t, f.postProg, "should not have post program")

	// Passes.
	assert.True(t, f.evalPre(pgn.MessageInfo{PGN: 127250, SourceId: 3}))
	// Fails on PGN.
	assert.False(t, f.evalPre(pgn.MessageInfo{PGN: 99999, SourceId: 3}))
	// Fails on source.
	assert.False(t, f.evalPre(pgn.MessageInfo{PGN: 127250, SourceId: 7}))
}

func TestFilterPostOnly(t *testing.T) {
	f, err := compileFilter("msg.Heading > 3.14")
	require.NoError(t, err)
	assert.False(t, f.preOnly, "should not be pre-only")
	assert.True(t, f.hasPost, "should have post filter")
	assert.Nil(t, f.preProg, "should not have pre program")
	assert.NotNil(t, f.postProg, "should have post program")

	// Passes.
	info := pgn.MessageInfo{PGN: 127250}
	assert.True(t, f.evalPostWithInfo(info, map[string]any{"Heading": 4.0}))
	// Fails.
	assert.False(t, f.evalPostWithInfo(info, map[string]any{"Heading": 1.0}))
}

func TestFilterSplit(t *testing.T) {
	f, err := compileFilter("pgn == 127250 && msg.Heading > 3.14")
	require.NoError(t, err)
	assert.False(t, f.preOnly, "should not be pre-only since msg is referenced")
	assert.True(t, f.hasPost, "should have post filter")
	assert.NotNil(t, f.preProg, "should have pre program for pgn check")
	assert.NotNil(t, f.postProg, "should have post program for msg check")

	info := pgn.MessageInfo{PGN: 127250}
	wrongInfo := pgn.MessageInfo{PGN: 99999}

	// Pre-filter passes, post-filter passes.
	assert.True(t, f.evalPre(info))
	assert.True(t, f.evalPostWithInfo(info, map[string]any{"Heading": 4.0}))

	// Pre-filter fails (wrong PGN) — no need to decode.
	assert.False(t, f.evalPre(wrongInfo))

	// Pre passes but post fails.
	assert.True(t, f.evalPre(info))
	assert.False(t, f.evalPostWithInfo(info, map[string]any{"Heading": 1.0}))
}

func TestFilterOrCannotSplit(t *testing.T) {
	f, err := compileFilter("source == 3 || msg.Heading > 1.0")
	require.NoError(t, err)
	// OR cannot be split: the whole expression references msg, so it all goes to post.
	assert.False(t, f.preOnly)
	assert.True(t, f.hasPost)
	assert.Nil(t, f.preProg, "OR expression cannot be split, no pre program")
	assert.NotNil(t, f.postProg)
}

func TestFilterCaseInsensitive(t *testing.T) {
	f, err := compileFilter("msg.heading > 1.0")
	require.NoError(t, err)
	assert.True(t, f.hasPost)

	info := pgn.MessageInfo{PGN: 127250}
	fields := map[string]any{"heading": 2.0, "Heading": 2.0}
	assert.True(t, f.evalPostWithInfo(info, fields))
}

func TestFilterInvalidExpression(t *testing.T) {
	_, err := compileFilter("this is not valid CEL !!!")
	assert.Error(t, err)
}

func TestFilterDestination(t *testing.T) {
	f, err := compileFilter("destination == 255")
	require.NoError(t, err)
	assert.True(t, f.preOnly)
	assert.True(t, f.evalPre(pgn.MessageInfo{TargetId: 255}))
	assert.False(t, f.evalPre(pgn.MessageInfo{TargetId: 0}))
}

func TestFilterPriority(t *testing.T) {
	f, err := compileFilter("priority < 4")
	require.NoError(t, err)
	assert.True(t, f.preOnly)
	assert.True(t, f.evalPre(pgn.MessageInfo{Priority: 2}))
	assert.False(t, f.evalPre(pgn.MessageInfo{Priority: 5}))
}

func TestStructToFilterMap(t *testing.T) {
	type testStruct struct {
		Info      pgn.MessageInfo
		Heading   *float32
		Deviation *float32
		Sid       *uint8
		Reference string
	}

	heading := float32(1.57)
	sid := uint8(42)
	s := testStruct{
		Info:      pgn.MessageInfo{PGN: 127250},
		Heading:   &heading,
		Deviation: nil, // nil pointer, should be skipped
		Sid:       &sid,
		Reference: "Magnetic",
	}

	m := structToFilterMap(s)

	// Info should be skipped.
	_, hasInfo := m["Info"]
	assert.False(t, hasInfo, "Info field should be skipped")

	// Heading should be present with both cases.
	assert.InDelta(t, 1.57, m["Heading"], 0.01)
	assert.InDelta(t, 1.57, m["heading"], 0.01)

	// Deviation (nil pointer) should be absent.
	_, hasDev := m["Deviation"]
	assert.False(t, hasDev, "nil pointer field should be skipped")

	// Sid should be present as int64.
	assert.Equal(t, int64(42), m["Sid"])
	assert.Equal(t, int64(42), m["sid"])

	// Reference should be present.
	assert.Equal(t, "Magnetic", m["Reference"])
	assert.Equal(t, "Magnetic", m["reference"])
}

func TestStructToFilterMapPointer(t *testing.T) {
	type simple struct {
		Value *float64
	}
	v := 3.14
	m := structToFilterMap(&simple{Value: &v})
	assert.Equal(t, 3.14, m["Value"])
}

func TestStructToFilterMapNonStruct(t *testing.T) {
	m := structToFilterMap(42)
	assert.Nil(t, m)
}
