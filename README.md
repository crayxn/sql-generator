# SQL Generator
前言：使用go-zero的时候只能使用写原生语句查询，故写了这个比较简单的sql语句生成器

## 基础使用
```go
// 基础查询
// select * from table_name where id = ? limit 1 offset 0
sql, vars := NewSql(func(q *Query) {
    q.Where("id", "=", 1)
	q.Find()
}, "table_name")

// 更新
// update table_name set status=? where id = ?
sql, vars := NewSql(func(q *Query) {
    q.Where("id", "=", 1)
    q.Update(map[string]interface{}{
        "status": 2,
    })
}, "table_name")

// 新增
// insert into table_name (status, title) values (?, ?)
sql, vars := NewSql(func(q *Query) {
    q.Insert(map[string]interface{}{
        "status": 2,
        "title":  "test",
    })
}, "table_name")

//删除
//delete from table_name where id = ?
sql, vars := NewSql(func(q *Query) {
    q.Where("id", "=", 1)
    q.Delete()
}, "table_name")

//支持链式
sql, vars := NewSql(func(q *Query) {
    q.Where("id", "=", 1).Where("status", ">", 1).Find()
}, "table_name")
```
## 条件查询
支持基本的 =、<、>、like、in、between、or、is null、is not null
#### func (w *Wheres) Where(field string, operator string, value interface{}) *Wheres
#### func (w *Wheres) WhereIn(field string, value []interface{}) *Wheres
#### func (w *Wheres) WhereBetween(field string, first interface{}, second interface{}) *Wheres
#### func (w *Wheres) WhereNull(field string) *Wheres
#### func (w *Wheres) WhereNoNull(field string) *Wheres

```go
q.Where("id", "=", 1)
// id = ?
q.Where("count", ">", 0)
// count > ?
q.Where("title", "like", "%test%")          //模糊
// title like ?
q.WhereIn("status", []interface{}{1, 2, 3}) //In
// status in (?,?,?)
q.WhereBetween("status", 1, 3)              // between
// status between (?,?)
q.WhereOr(func(or *where.Wheres) {
or.Where("id","=",1)
or.WhereIn("status", []interface{}{ 1, 3 })
})
// (id = ? or status in (?,?))
```
## Limit
#### func (q *Query) Limit(offset int64, limit int64) *Query
```go 
sql, vars := NewSql(func(q *Query) {
    q.Limit(0, 10)	
    // select * from table_name limit 10 offest 0
}, "table_name")
```
## OrderBy
#### func (q *Query) OrderBy(field string, direction string) *Query
```go
sql, vars := NewSql(func(q *Query) {
    q.OrderBy("id", "desc") 
    //select * from table_name order by id desc
    q.OrderBy("status", "asc") 
    //select * from table_name order by id desc, status asc
}, "table_name")
```
## GroupBy / Having
#### func (q *Query) GroupBy(field string) *Query
#### func (q *Query) Having(do func(having *where.Wheres)) *Query
```go
sql, vars := NewSql(func(q *Query) {
    q.GroupBy("status").Having(func(having *where.Wheres) {
        having.Where("count(id)", ">", 5)
    })
// select * from table_name order by status having ( count(id) > ? )
}, "table_name")
```
## Join
#### func (q *Query) Join(item *join.Joins, do func(join *join.Joins)) *Query
```go
sql, vars := NewSql(func(q *Query) {
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

// select q1.id, q2.id from table_name as t1 inner join table_name2 as t2 on t1.id = t2.table_name_id and t2.status = ? limit 10 offset 0

```
## 搭配使用（go-zero）
```go
productList := make([]*Product, 0)
searchSql, vars := query.NewSql(func(q *query.Query) {
    q.UseSoftDelete()
    // q.Select(productRows)
    q.Limit(0, 10)
    q.Where("status", "=", 1)
}, m.table)
// select * from product where status = ? and delete_time is null limit 10 offset 0
if err := m.QueryRowsNoCache(&productList, searchSql, vars...); err != nil {
    return 0, nil, err
}
```

## 性能
因为实在简单，没有多少东西 10w 句sql生成量只需大约 100ms
