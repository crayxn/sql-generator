package where

import "fmt"

type Wheres struct {
	Raws   []string
	Values []interface{}
}

func (w *Wheres) Where(field string, operator string, value interface{}) *Wheres {
	if value != "" {
		w.Raws = append(w.Raws, fmt.Sprintf("%s %s ?", field, operator))
		w.Values = append(w.Values, value)
	}
	return w
}

func (w *Wheres) WhereIn(field string, value []interface{}) *Wheres {
	if len(value) > 0 {
		temp := "?"
		for i := 1; i < len(value); i++ {
			temp += ",?"
		}
		w.Raws = append(w.Raws, fmt.Sprintf("%s in (%s)", field, temp))
		w.Values = append(w.Values, value...)
	}
	return w
}

func (w *Wheres) WhereBetween(field string, first interface{}, second interface{}) *Wheres {
	w.Raws = append(w.Raws, fmt.Sprintf("%s between ? and ?", field))
	w.Values = append(w.Values, first, second)
	return w
}

func (w *Wheres) WhereNull(field string) *Wheres {
	w.Raws = append(w.Raws, fmt.Sprintf("%s is null", field))
	return w
}

func (w *Wheres) WhereNoNull(field string) *Wheres {
	w.Raws = append(w.Raws, fmt.Sprintf("%s is not null", field))
	return w
}
