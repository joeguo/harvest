package main

import (
	"testing"
	"fmt"
	_"database/sql"
)

func aTestMySQL(t *testing.T) {
	mysql, _ := NewMySQL("harvest/root/apples")
	defer mysql.Close()
	ds,err:=mysql.Uncrawled(30)
	if err!=nil{
		fmt.Println(err)
		return
	}
	for _,d:=range(ds){
		fmt.Println(d)
	}

}

