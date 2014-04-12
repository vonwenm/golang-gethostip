package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	lock    chan bool         //lock of getHost()
	domains map[string]string //result of getHost()
)

func main() {
	lock = make(chan bool)
	domains = make(map[string]string, max_routines)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", db_user, db_pass, db_host, db_port, db_name)

	//
	var temp_offset int = 0
	var temp_routines int = 0
	var temp_domain string
	for {
		//
		conn, err := sql.Open("mysql", dsn)
		if err != nil {
			fmt.Println(err)
			return
		}

		sqlstr := fmt.Sprintf("select domain from %s limit %d,%d", src_table_name, temp_offset, max_query_rows)
		fmt.Println(sqlstr)
		rows, err := conn.Query(sqlstr)
		if err != nil {
			fmt.Println("mysql query error #2: ", err)
			continue
		} else if false == rows.Next() {
			fmt.Println("job done")
			return
		}

		temp_offset += max_query_rows
		for i := 0; i < max_query_rows; i++ {
			if temp_routines < max_routines {
				rows.Scan(&temp_domain)
				go getHost(temp_domain)
				temp_routines++

				if false == rows.Next() {
					goto RELEASE
				}
				continue
			}

		RELEASE:
			for j := 0; j < temp_routines; j++ {
				<-lock
			}
			temp_routines = 0

			for k, v := range domains {
				if len(v) < 1 {
					continue
				}
				timestamp := time.Now().Unix()
				sqlstr = fmt.Sprintf("update %s set ip='%s',refresh_time=%d where domain='%s'", src_table_name, v, timestamp, k)
				fmt.Println(sqlstr)
				_, err := conn.Exec(sqlstr)
				if err != nil {
					fmt.Println("mysql query error #3:", err, sqlstr)
				}
				delete(domains, k)
			}
		}
		rows.Close()
		conn.Close()
	}
}
