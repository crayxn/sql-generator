package join

import (
	"fmt"
	"github.com/crayoon/sql_generator/where"
)

const (
	InnerJoin = "inner"
	LeftJoin  = "left"
	RightJoin = "right"
)

type Joins struct {
	Typ    string
	Table  string
	Wheres *where.Wheres
	JoinOn string
}

func (join *Joins) On(first string, operator string, second string) *Joins {
	join.JoinOn = fmt.Sprintf("%s %s %s", first, operator, second)
	return join
}

func (join *Joins) Where(do func(w *where.Wheres)) *Joins {
	newWhere := where.Wheres{}
	do(&newWhere)
	join.Wheres = &newWhere
	return join
}
