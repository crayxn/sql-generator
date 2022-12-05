package sql_generator

import (
	"fmt"
	"github.com/crayoon/sql_generator/join"
	"github.com/crayoon/sql_generator/where"
	"strings"
	"time"
)

const (
	TypeSelect = "select"
	TypeInsert = "insert"
	TypeUpdate = "update"
	TypeDelete = "delete"
)

type Query struct {
	table string
	//当前操作
	typ string

	softDelete bool
	//查询
	selects string
	//更新
	update struct {
		fields []string
		values []interface{}
	}
	//新增
	insert struct {
		fields []string
		values []interface{}
	}

	//条件
	wheres *where.Wheres
	//限制
	limit  int64
	offset int64
	//排序
	orderBy []string
	//组
	groupBy []string
	having  *where.Wheres
	//关联
	joins []*join.Joins
}

func NewSql(do func(q *Query), table string) (string, []interface{}) {
	query := Query{
		table:      table,
		selects:    "*",
		softDelete: false,
		wheres:     &where.Wheres{},
		having:     &where.Wheres{},
	}
	do(&query)
	return query.ToSql()
}

func (q *Query) Find() *Query {
	q.limit = 1
	return q
}

func (q *Query) Select(fields string) *Query {
	q.typ = TypeSelect
	q.selects = fields
	return q
}

func (q *Query) Insert(v map[string]interface{}) *Query {
	q.typ = TypeInsert
	for key, value := range v {
		q.insert.fields = append(q.insert.fields, key)
		q.insert.values = append(q.insert.values, value)
	}
	return q
}

func (q *Query) Update(v map[string]interface{}) *Query {
	q.typ = TypeUpdate
	for key, value := range v {
		q.update.fields = append(q.update.fields, key+"=?")
		q.update.values = append(q.update.values, value)
	}
	return q
}

func (q *Query) Delete() *Query {
	q.typ = TypeDelete
	return q
}

func (q *Query) Where(field string, operator string, value interface{}) *Query {
	q.wheres.Where(field, operator, value)
	return q
}

func (q *Query) WhereIn(field string, value []interface{}) *Query {
	q.wheres.WhereIn(field, value)
	return q
}

func (q *Query) WhereBetween(field string, first interface{}, second interface{}) *Query {
	q.wheres.WhereBetween(field, first, second)
	return q
}

func (q *Query) WhereNull(field string) *Query {
	q.wheres.WhereNull(field)
	return q
}
func (q *Query) WhereNoNull(field string) *Query {
	q.wheres.WhereNoNull(field)
	return q
}

func (q *Query) AddWhere(do func(w *where.Wheres)) *Query {
	newWhere := where.Wheres{}
	do(&newWhere)
	if len(newWhere.Raws) > 0 {
		q.wheres.Raws = append(q.wheres.Raws, fmt.Sprintf("( %s )", strings.Join(newWhere.Raws, " and ")))
		q.wheres.Values = append(q.wheres.Values, newWhere.Values...)
	}
	return q
}

func (q *Query) WhereOr(do func(or *where.Wheres)) *Query {
	orWhere := where.Wheres{}
	do(&orWhere)
	if len(orWhere.Raws) > 0 {
		q.wheres.Raws = append(q.wheres.Raws, fmt.Sprintf("( %s )", strings.Join(orWhere.Raws, " or ")))
		q.wheres.Values = append(q.wheres.Values, orWhere.Values...)
	}
	return q
}

func (q *Query) Limit(offset int64, limit int64) *Query {
	q.limit = limit
	q.offset = offset
	return q
}

func (q *Query) UseSoftDelete() *Query {
	q.WhereNull("delete_time")
	q.softDelete = true
	return q
}

func (q *Query) OrderBy(field string, direction string) *Query {
	q.orderBy = append(q.orderBy, fmt.Sprintf("%s %s", field, direction))
	return q
}

func (q *Query) GroupBy(field string) *Query {
	q.orderBy = append(q.groupBy, field)
	return q
}

func (q *Query) Having(do func(having *where.Wheres)) *Query {
	do(q.having)
	return q
}

func (q *Query) Join(item *join.Joins, do func(join *join.Joins)) *Query {
	if do != nil {
		do(item)
	}
	q.joins = append(q.joins, item)
	return q
}

func (q *Query) Count(field string) *Query {
	q.Select(fmt.Sprintf("count(%s) as _COUNT", field))
	return q
}

func (q *Query) ToSql() (sql string, values []interface{}) {
	//type
	switch q.typ {
	case TypeDelete:
		if len(q.wheres.Raws) < 1 {
			panic("删除必须加条件")
		}
		if q.softDelete {
			sql = fmt.Sprintf("update %s set delete_time=? ", q.table)
			values = []interface{}{time.Now().Format("2006-01-02 15:04:05")}
			return
		} else {
			sql = fmt.Sprintf("delete from %s ", q.table)
		}
	case TypeUpdate:
		if len(q.update.fields) < 1 {
			return
		}
		sql = fmt.Sprintf("update %s set %s ", q.table, strings.Join(q.update.fields, ", "))
		values = q.update.values
	case TypeInsert:
		temp, insLen := "?", len(q.insert.fields)
		if insLen < 1 {
			return
		}
		for i := 1; i < insLen; i++ {
			temp += ", ?"
		}
		sql = fmt.Sprintf("insert into %s (%s) values (%s)", q.table, strings.Join(q.insert.fields, ", "), temp)
		values = q.insert.values
		return
	default:
		sql = fmt.Sprintf("select %s from %s ", q.selects, q.table)
	}
	//join
	if len(q.joins) > 0 {
		for _, joinItem := range q.joins {
			if joinItem.Table == "" || joinItem.JoinOn == "" {
				continue
			}
			if joinItem.Typ == "" {
				joinItem.Typ = join.InnerJoin
			}
			sql += fmt.Sprintf("%s join %s on %s", joinItem.Typ, joinItem.Table, joinItem.JoinOn)
			if joinItem.Wheres != nil && len(joinItem.Wheres.Raws) > 0 {
				sql += fmt.Sprintf(" and %s ", strings.Join(joinItem.Wheres.Raws, " and "))
				values = append(values, joinItem.Wheres.Values...)
			}
		}
	}

	//where
	if len(q.wheres.Raws) > 0 {
		sql += fmt.Sprintf("where %s ", strings.Join(q.wheres.Raws, " and "))
		values = append(values, q.wheres.Values...)
	}
	//order by
	if len(q.orderBy) > 0 {
		sql += fmt.Sprintf("order by %s ", strings.Join(q.orderBy, ", "))
		//having
		sql += fmt.Sprintf("having ( %s )", strings.Join(q.having.Raws, " and "))
		values = append(values, q.having.Values...)
	}
	//group by + having
	if len(q.groupBy) > 0 {
		sql += fmt.Sprintf("group by %s ", strings.Join(q.groupBy, ", "))
		//having
		sql += fmt.Sprintf("having ( %s )", strings.Join(q.having.Raws, " and "))
		values = append(values, q.having.Values...)
	}
	//limit
	if q.limit > 0 && q.typ == TypeSelect {
		sql += fmt.Sprintf("limit %d offset %d ", q.limit, q.offset)
	}

	return
}
