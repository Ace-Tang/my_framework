package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	err := create()
	if err == nil {
		log.Println("create ok")
	}
}

func test() {
	db, err := sql.Open("mysql", "root:ace@tcp(127.0.0.1:3306)/")
	if err != nil {
		log.Println("open db ", err)
		return
	}

	_, err = db.Exec("use mesos")
	if err != nil {
		log.Println("use mesos ", err)
		return
	}

	//sql := "insert into `task_info` (`task_cpu`, `task_mem`, `id`, `cmd`, `env`, `image`, `hostname`, `name`, `status`, `count`, `slave_id`, `framework_id`) values (?,?,?,?,?,?,?,?,?,?,?,?);"

}

func create() error {
	db, err := sql.Open("mysql", "root:ace@tcp(127.0.0.1:3306)/")
	if err != nil {
		log.Println("open db ", err)
		return err
	}

	//_, err = db.Exec("create database mesos;")
	//if err != nil {
	//	log.Println("create db err")
	//}

	_, err = db.Exec("use mesos")
	if err != nil {
		log.Println("use mesos ", err)
		return err
	}

	sql := "create table slave_info ( hostname varchar(20) not null primary key, attachment varchar(255) not null)"
	_, err = db.Exec(sql)
	if err != nil {
		log.Println("create table slave_info ", err)
	}

	sql = "create table task_info ( task_cpu float not null, task_mem float not null, id varchar(255) primary key, name varchar(20) not null, cmd varchar(255)  not null, env varchar(255)  not null, image varchar(50) not null, slave_id varchar(255)  not null, hostname varchar(20) not null, framework_id varchar(255), status int not null, count int not null)"
	_, err = db.Exec(sql)
	if err != nil {
		log.Println("create table task_info ", err)
	}

	return nil
}
