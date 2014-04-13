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
	start_time := time.Now().Unix()
	lock = make(chan bool)
	domains = make(map[string]string, max_routines)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", db_user, db_pass, db_host, db_port, db_name)

	//
	var temp_domain string
	var offset, cur_routines, update_num int
	var routine_count int
	var is_finish bool = false
	for {
		conn, err := sql.Open("mysql", dsn)
		if err != nil {
			fmt.Println("mysql query error #1: ", err)
			return
		}
		sqlstr := fmt.Sprintf("select domain from %s order by domain asc limit %d,%d", domain_table, offset, max_rows)
		if app_debug {
			fmt.Println(sqlstr)
		}
		rows, err := conn.Query(sqlstr)
		if err != nil {
			fmt.Println("mysql query error #2: ", err)
			continue
		} else if false == rows.Next() {
			end_time := time.Now().Unix()
			fmt.Println("job done")
			fmt.Printf("use time: %d seconds,update_num=%d,routine_count=%d\n", end_time-start_time, update_num, routine_count)
			return
		}

		//
		offset += max_rows
		for i := 0; i < max_rows; i++ {
			if is_finish {
				break
			}
			if cur_routines < max_routines {
				rows.Scan(&temp_domain)
				if app_debug {
					fmt.Printf("read domain: %s\n", temp_domain)
				}

				go getHost(temp_domain)
				routine_count++
				cur_routines++

				if false == rows.Next() {
					if cur_routines < max_routines {
						is_finish = true
					}
					goto RUNSQL
				}
				continue
			}

		RUNSQL:
			for j := 0; j < cur_routines; j++ {
				temp := <-lock
				if app_debug {
					fmt.Printf("read routine result:%v\n", temp)
				}

			}
			cur_routines = 0

			for k, v := range domains {
				if len(k) < 1 {
					delete(domains, k)
					continue
				}

				timestamp := time.Now().Unix()
				sqlstr = fmt.Sprintf("update %s set ip='%s',refresh_time=%d where domain='%s'", domain_table, v, timestamp, k)
				if app_debug {
					fmt.Println(sqlstr)
				}
				_, err := conn.Exec(sqlstr)
				update_num++
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
