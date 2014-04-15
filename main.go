package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	lock    chan bool
	domains map[string]string
	indexes []string
)

func main() {
	start_time := time.Now().Unix()
	lock = make(chan bool)
	indexes = make([]string, max_rows)
	domains = make(map[string]string, max_rows)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", db_user, db_pass, db_host, db_port, db_name)

	//init
	var offset, update_num, routine_num int
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("mysql query error #1: ", err)
		return
	}

	//do refresh
	for {
		sqlstr := fmt.Sprintf("select domain from %s order by domain asc limit %d,%d", domain_table, offset, max_rows)
		if app_debug {
			fmt.Println(sqlstr)
		}
		rows, err := conn.Query(sqlstr)
		if err != nil {
			fmt.Println("mysql query error #2: ", err)
			return
		} else if false == rows.Next() {
			end_time := time.Now().Unix()
			fmt.Println("job done")
			fmt.Printf("use time: %d seconds,update_num=%d,routine_num=%d\n", end_time-start_time, update_num, routine_num)
			return
		}

		//get domains from db
		query_rows := 0
		temp_domain := ""
		offset += max_rows
		for query_rows = 0; query_rows < max_rows; query_rows++ {
			rows.Scan(&temp_domain)
			indexes[query_rows] = temp_domain
			domains[temp_domain] = "0.0.0.0"
			if app_debug {
				fmt.Printf("read domain: %s\n", temp_domain)
			}
			if rows.Next() == false {
				query_rows++
				break
			}
		}
		rows.Close()

		cur_routines := 0
		cur_domains := make([]int, max_routines)
		for i := 0; i < query_rows; i++ {
			temp_domain = indexes[i]
			cur_domains[cur_routines] = i
			go getHost(temp_domain)
			//if app_debug {
			//	fmt.Printf("get domain: %s\n", temp_domain)
			//}
			cur_routines++
			routine_num++
			if max_rows != query_rows+1 {
				goto RUNSQL
			}
			if cur_routines < max_routines {
				continue
			}

			//
		RUNSQL:
			for j := 0; j < cur_routines; j++ {
				<-lock
				//temp := <-lock
				//if app_debug {
				//	fmt.Printf("read routine result:%v\n", temp)
				//}
			}

			for j := 0; j < cur_routines; j++ {
				timestamp := time.Now().Unix()
				sqlstr = fmt.Sprintf("update %s set ip='%s',refresh_time=%d where domain='%s'", domain_table, domains[indexes[j]], timestamp, indexes[j])
				if app_debug {
					fmt.Println(sqlstr)
				}
				_, err := conn.Exec(sqlstr)
				update_num++
				if err != nil {
					fmt.Println("mysql query error #3:", err, sqlstr)
				}
			}
			cur_routines = 0
		}
	}
}
