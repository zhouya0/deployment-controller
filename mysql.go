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

func InitDB() error {
	path := strings.Join([]string{username, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	DB, _ = sql.Open("mysql", path)
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)

	if err := DB.Ping(); err != nil {
		fmt.Printf("opon database fail: %v\n", err)
		return err
	}
	fmt.Println("connnect success")
	return nil
}

func InserDeployment(dp *appsv1.Deployment) bool {
	tx, err := DB.Begin()
	if err != nil {
		fmt.Println("tx fail")
		return false
	}
	stmt, err := tx.Prepare("INSERT INTO deployment_test (`deployment_name`, `deployment_namespace`, `desired_replicas`, `ready_replicas`) VALUES (?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Prepare fail")
		return false
	}
	res, err := stmt.Exec(dp.Name, dp.Namespace, dp.Spec.Replicas, dp.Status.ReadyReplicas)
	if err != nil {
		fmt.Println("Exec fail")
		return false
	}
	tx.Commit()
	fmt.Println(res.LastInsertId())
	return true

}

func UpdateDeployment(dp *appsv1.Deployment) bool {
	tx, err := DB.Begin()
	if err != nil {
		fmt.Println("tx fail")
	}

	stmt, err := tx.Prepare("UPDATE deployment_test SET desired_replicas = ? , ready_replicas = ? WHERE deployment_name = ?")

	if err != nil {
		fmt.Println("Prepare fail")
		return false
	}

	res, err := stmt.Exec(dp.Spec.Replicas, dp.Status.ReadyReplicas, dp.Name)
	if err != nil {
		fmt.Println("Exec fail")
		return false
	}

	tx.Commit()
	fmt.Println(res.LastInsertId())
	return true
}

func DeleteDeployment(key string) bool {
	strs := strings.Split(key, "/")
	name := strs[1]
	tx, err := DB.Begin()
	if err != nil {
		fmt.Println("tx fail")
	}
	fmt.Println("The key is:", name)
	stmt, err := tx.Prepare("DELETE FROM deployment_test WHERE deployment_name = ?")

	if err != nil {
		fmt.Println("Prepare fail")
		return false
	}

	res, err := stmt.Exec(name)
	if err != nil {
		fmt.Println("Exec fail")
		return false
	}

	tx.Commit()
	fmt.Println(res.LastInsertId())
	return true
}
