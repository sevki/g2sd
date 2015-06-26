// Copyright 2015 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sql // import "sevki.org/g2sd/sql"

import (
	"log"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"sevki.org/graphql/parse"
	"sevki.org/graphql/query"
	"sevki.org/lib/prettyprint"
)

func d(n *query.Node, db *sqlx.DB, t *testing.T) interface{} {

	if res, err := exec(n, db); err != nil {
		t.Fatal(err)
		return nil
	} else {
		return res
	}

}

func TestConnection(t *testing.T) {

	db, err := sqlx.Open("sqlite3", "./test_db.db")
	if err != nil {
		t.Fatal(err)
	}

	const q = `{
  Products(ProductID: 3) {
    ProductName
  }
}`
	if ast, err := parse.NewQuery([]byte(q)); err != nil {
		t.Error(err.Error())
	} else {
		log.Println(prettyprint.AsJSON(d(ast, db, t)))
	}

	defer db.Close()

}

func TestQuery(t *testing.T) {

	db, err := sqlx.Open("sqlite3", "./test_db.db")
	if err != nil {
		t.Fatal(err)
	}

	const q = `{
  Products(ProductID: 3) {
    ProductName,
    UnitsInStock,
    CategoryId,
    Categories() {
       CategoryName
    }
  }
}`
	if ast, err := parse.NewQuery([]byte(q)); err != nil {
		t.Error(err.Error())
	} else {
		log.Println(prettyprint.AsJSON(d(ast, db, t)))
	}

	defer db.Close()

}

func TestQueryVariable(t *testing.T) {

	db, err := sqlx.Open("sqlite3", "./test_db.db")
	if err != nil {
		t.Fatal(err)
	}

	const q = `{
  Products(ProductID: 3) {
    ProductName,
    UnitsInStock,
    CategoryId,
    Categories(CategoryID: $CategoryID) {
       CategoryName
    }
  }
}`
	if ast, err := parse.NewQuery([]byte(q)); err != nil {
		t.Error(err.Error())
	} else {
		//		log.Println(prettyprint.AsJSON(d(ast, db, t)))
		d(ast, db, t)
	}

	defer db.Close()

}
func TestQueryComplex(t *testing.T) {

	db, err := sqlx.Open("sqlite3", "./test_db.db")
	if err != nil {
		t.Fatal(err)
	}

	const q = `{
  Products(ProductID: 9) {
    ProductName,
    UnitsInStock,
    ProductID,
    OrderDetails(ProductID: $ProductID) {
       OrderID,
       Orders(OrderID: $OrderID) {
         EmployeeID,
         Employees(EmployeeID: $EmployeeID) {
            FirstName
         }
       }
    }
  }
}`
	if ast, err := parse.NewQuery([]byte(q)); err != nil {
		t.Error(err.Error())
	} else {
		log.Println(prettyprint.AsJSON(d(ast, db, t)))
		//d(ast, db, t)
	}

	defer db.Close()

}
