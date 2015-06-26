// Copyright 2015 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate stringer -type Type

package sql // import "sevki.org/g2sd/sql"

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"sevki.org/graphql/query"
)

type Select struct {
	Distinct    string
	SelectExprs StrArr
	From        string
	Where       Wheres
	GroupBy     string
}
type StrArr []string
type Wheres []string

func (node *Select) String() string {
	return fmt.Sprintf("select %v from %v %v;",
		node.SelectExprs,
		node.From, node.Where,
	)
}

func New(n *query.Node) (s Select) {
	s.From = string(n.Name)
	for _, e := range n.Edges {
		if e.Params == nil {
			s.SelectExprs = append(s.SelectExprs, string(e.Name))
		}
	}
	for k, v := range n.Params {
		s.Where = append(s.Where, fmt.Sprintf("%s = %v", k, v))
	}
	return

}

func (s StrArr) String() (str string) {
	for i, e := range s {
		if i != 0 {
			str += ", "
		}
		str += e
	}

	return
}

func (s Wheres) String() (str string) {
	for i, e := range s {
		if i != 0 {
			str += ", "
		} else {
			str += "where "
		}
		str += e
	}

	return
}
func exec(n *query.Node, db *sqlx.DB) (interface{}, error) {
	slct := New(n)
	rows, err := db.Queryx(slct.String())

	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		result := make(map[string]interface{})

		cols, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		data, err := rows.SliceScan()

		for i, col := range cols {
			switch data[i].(type) {
			case []uint8:
				result[col] = string(data[i].([]uint8))
				break
			default:
				result[col] = data[i]
				break
			}
		}
		for _, e := range n.Edges {
			if e.Params != nil {
				t, _ := query.ApplyContext(&e, result)
				result[string(e.Name)], err = exec(t, db)
			}
		}

		results = append(results, result)
	}
	rows.Close()

	if len(results) == 1 {
		return results[0], nil
	} else {
		return results, nil
	}
}
