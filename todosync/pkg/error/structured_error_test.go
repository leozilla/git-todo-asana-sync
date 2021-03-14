package error

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DeepKeyVals_NoStructErr(t *testing.T) {
	err := errors.New("no structured error")
	keyvals := DeepKeyVals(err)
	assert.Equal(t, []interface{}{"error", "no structured error"}, keyvals)
}

func Test_DeepKeyVals_StructErrNoKeyVals(t *testing.T) {
	err := Err("msg", errors.New("inner error"))
	keyvals := DeepKeyVals(err)
	assert.Len(t, keyvals, 3*2)
	assertBaseKeyVals(t, keyvals)
	assertContainsKeyVal(t, "error", "msg: inner error", keyvals)
}

func Test_DeepKeyVals_StructErrWithKeyVals(t *testing.T) {
	err := Err("msg", errors.New("inner error"), "key1", "val1")
	keyvals := DeepKeyVals(err)
	assert.Len(t, keyvals, 4*2)
	assertBaseKeyVals(t, keyvals)
	assertContainsKeyVal(t, "error", "msg: inner error", keyvals)
	assertContainsKeyVal(t, "key1", "val1", keyvals)
}

func Test_DeepKeyVals_NestedStructErrNoKeyVals(t *testing.T) {
	err := Err("outer msg", Err("inner msg", errors.New("inner error")))
	keyvals := DeepKeyVals(err)
	assert.Len(t, keyvals, 3*2)
	assertBaseKeyVals(t, keyvals)
	assertContainsKeyVal(t, "error", "outer msg: inner msg: inner error", keyvals)
}

func Test_DeepKeyVals_NestedStructErrWithKeyVals(t *testing.T) {
	err := Err("outer msg", Err("inner msg", errors.New("inner error"), "key1", "val1"), "key2", "val2")
	keyvals := DeepKeyVals(err)
	assert.Len(t, keyvals, 5*2)
	assertBaseKeyVals(t, keyvals)
	assertContainsKeyVal(t, "error", "outer msg: inner msg: inner error", keyvals)
	assertContainsKeyVal(t, "key1", "val1", keyvals)
	assertContainsKeyVal(t, "key2", "val2", keyvals)
}

func Test_DeepKeyValsWithPairs_StructErrWithKeyVals(t *testing.T) {
	err := Err("msg", errors.New("inner error"), "key1", "val1")
	keyvals := DeepKeyVals(err, "key2", "val2")
	assert.Len(t, keyvals, 5*2)
	assertBaseKeyVals(t, keyvals)
	assertContainsKeyVal(t, "error", "msg: inner error", keyvals)
	assertContainsKeyVal(t, "key1", "val1", keyvals)
	assertContainsKeyVal(t, "key2", "val2", keyvals)
}

func Test_DeepKeyValsWithKeyVals_NestedStructErrWithKeyVals(t *testing.T) {
	err := Err("outer msg", Err("inner msg", errors.New("inner error"), "key1", "val1"), "key2", "val2")
	keyvals := DeepKeyVals(err, "key3", "val3")
	assert.Len(t, keyvals, 6*2)
	assertBaseKeyVals(t, keyvals)
	assertContainsKeyVal(t, "error", "outer msg: inner msg: inner error", keyvals)
	assertContainsKeyVal(t, "key1", "val1", keyvals)
	assertContainsKeyVal(t, "key2", "val2", keyvals)
	assertContainsKeyVal(t, "key3", "val3", keyvals)
}

func Test_Root_OneLevel(t *testing.T) {
	strucerr := Err("msg", errors.New("inner error"), "key1", "val1").(*StructuredError)
	root := strucerr.Root()
	assert.Equal(t, strucerr, root)
}

func Test_Root_MoreThanOneLevel(t *testing.T) {
	inner := Err("inner msg", errors.New("inner error"), "key1", "val1")
	strucerr := Err("outer msg", inner, "key2", "val2").(*StructuredError)
	root := strucerr.Root()
	assert.Equal(t, inner, root)
}

func assertBaseKeyVals(t *testing.T, keyvals []interface{}) {
	assertContainsKeyWithVal(t, "debugID", notEmpty, keyvals)
	assertContainsKeyWithVal(t, "time", notEmpty, keyvals)
}

func assertContainsKeyVal(t *testing.T, key string, val string, actual []interface{}) {
	assertContainsKeyWithVal(t, key, func(actual string) bool {
		return actual == val
	}, actual)
}

func assertContainsKeyWithVal(t *testing.T, key string, valF func(string) bool, actual []interface{}) {
	i := findKey(actual, key)
	assert.GreaterOrEqual(t, i, 0, "key index", i)
	assert.GreaterOrEqual(t, len(actual), i+1, "value index")
	assert.True(t, valF(actual[i+1].(string)), "value")
}

func findKey(actual []interface{}, key string) int {
	for i, k := range actual {
		if key == k {
			return i
		}
	}

	return -1
}

func notEmpty(s string) bool {
	return s != ""
}
