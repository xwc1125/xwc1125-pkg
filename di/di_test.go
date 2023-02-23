// Package di
//
// @author: xwc1125
package di

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type A struct {
	Db  *sql.DB `di:"db"`
	Db1 *sql.DB `di:"db"`
	B   *B      `di:"b,prototype"`
	B1  *B      `di:"b,prototype"`
	bb  *B      `di:"b,prototype"`
	bb1 *B      `di:"b,prototype"`
}

func NewA() *A {
	return &A{}
}

// SetBb 对于私有字段，可以使用Set方法反射赋值
func (p *A) SetBb(bb *B) {
	p.bb = bb
}
func (p *A) SetBb1(bb *B) error {
	p.bb1 = bb
	return nil
}

func (p *A) Version() (string, error) {
	rows, err := p.Db.Query("SELECT VERSION() as version")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var version string
	if rows.Next() {
		if err := rows.Scan(&version); err != nil {
			return "", err
		}
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	return version, nil
}

type B struct {
	Name string
}

func NewB() *B {
	return &B{
		Name: time.Now().String(),
	}
}

func TestDi(t *testing.T) {
	container := NewContainer()
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/tchain_ledger?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}
	container.SetSingleton("db", db)
	container.SetPrototype("b", func() (interface{}, error) {
		return NewB(), nil
	})

	a := NewA()
	if err := container.Ensure(a); err != nil {
		fmt.Println(err)
		return
	}
	// 打印指针，确保单例和实例的指针地址
	fmt.Printf("db: %p\ndb1: %p\nb: %p\nb1: %p\n", a.Db, a.Db1, &a.B, &a.B1)
}
