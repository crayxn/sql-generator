package sql_generator

import (
	"github.com/crayoon/sql_generator/join"
	"github.com/crayoon/sql_generator/where"
	"reflect"
	"testing"
)

func TestQuery_ToSql(t *testing.T) {
	t.Run("test-normal-sql", func(t *testing.T) {
		sql, params := NewSql(func(q *Query) {
			q.Where("id", "=", 1)
			q.Delete()
		}, "table_name")
		if !reflect.DeepEqual(sql, "select * from table_name where id = ? limit 1 offset 0 ") {
			t.Errorf("fail sql = %v", sql)
		}
		if !reflect.DeepEqual(params, []interface{}{
			1,
		}) {
			t.Errorf("fail params = %v", params)
		}
	})
	t.Run("test-join-sql", func(t *testing.T) {
		sql, params := NewSql(func(q *Query) {
			q.Join(&join.Joins{
				Typ:   join.InnerJoin,
				Table: "table_name2 as t2",
			}, func(join *join.Joins) {
				join.On("t1.id", "=", "t2.table_name_id")
				join.Where(func(w *where.Wheres) {
					w.Where("t2.status", "=", 1)
				})
			})
			q.Limit(0, 10)
			q.Select("q1.id, q2.id")
		}, "table_name as t1")
		if !reflect.DeepEqual(sql, "select q1.id, q2.id from table_name as t1 inner join table_name2 as t2 on t1.id = t2.table_name_id and t2.status = ? limit 10 offset 0 ") {
			t.Errorf("fail sql = %v", sql)
		}
		if !reflect.DeepEqual(params, []interface{}{
			1,
		}) {
			t.Errorf("fail params = %v", params)
		}
	})
}
