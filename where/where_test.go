package where

import (
	"reflect"
	"testing"
)

func TestWheres_Where(t *testing.T) {
	t.Run("test-where", func(t *testing.T) {
		w := &Wheres{
			[]string{},
			[]interface{}{},
		}
		w.Where("field", "=", 1)
		if got := w.Raws; !reflect.DeepEqual(got, []string{
			"`field` = ?",
		}) {
			t.Errorf("Fail got = %v", got)
		}
	})
}
