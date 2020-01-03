package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	appsv1 "k8s.io/api/apps/v1"
	"strings"
)

const (
	username = "root"
	password = "password"
	ip       = "127.0.0.1"
	port     = "30944"
	dbName   = "RUNOOB"
)

var DB *sql.DB

func InitDB() {
	path := strings.Join([]string{username, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	DB, _ = sql.Open("mysql", path)
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)

	if err := DB.Ping(); err != nil {
		fmt.Printf("opon database fail: %v\n", err)
		return
	}
	fmt.Println("connnect success")
}

func InserDeployment(dp *appsv1.Deployment) bool {
	tx, err := DB.Begin()
	if err != nil {
		fmt.Println("tx fail")
		return false
	}
	stmt, err := tx.Prepare("INSERT INTO deployment (`deployment_name`, `deployment_status`) VALUES (?, ?)")
	if err != nil {
		fmt.Println("Prepare fail")
		return false
	}
	res, err := stmt.Exec(dp.Name, "Ready")
	if err != nil {
		fmt.Println("Exec fail")
		return false
	}
	tx.Commit()
	fmt.Println(res.LastInsertId())
	return true

}
