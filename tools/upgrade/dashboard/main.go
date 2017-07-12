/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yubo/falcon/ctrl/api/models"
	"log"
)

const (
	src_dsn = "root:123456@tcp(localhost:3306)/mi_dashboard"
	dst_dsn = "root:123456@tcp(localhost:3306)/falcon"
)

var (
	src, dst *sql.DB
	err      error
)

func main() {

	src, err = sql.Open("mysql", src_dsn)
	if err != nil {
		goto err_out
	}

	dst, err = sql.Open("mysql", dst_dsn)
	if err != nil {
		goto err_out
	}

	if err = do_dashboard_graph(); err != nil {
		goto err_out
	}
	if err = do_dashboard_screen(); err != nil {
		goto err_out
	}
	if err = do_tmp_graph(); err != nil {
		goto err_out
	}
	return
err_out:
	log.Println(err)
}

func do_dashboard_graph() error {
	var o models.DashboardGraph
	do_sql(dst, "TRUNCATE TABLE dashboard_graph")
	set_auto_increment("dashboard_graph")

	rows, err := src.Query("select id, title, hosts, counters, screen_id, timespan, graph_type, method, position, falcon_tags from dashboard_graph")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&o.Id,
			&o.Title,
			&o.Hosts,
			&o.Counters,
			&o.ScreenId,
			&o.Timespan,
			&o.GraphType,
			&o.Method,
			&o.Position,
			&o.FalconTags)
		if err != nil {
			return err
		}
		err = do_sql(dst, "insert into dashboard_graph(id, title, hosts, counters, screen_id, timespan, graph_type, method, position, falcon_tags) values(?,?,?,?,?,?,?,?,?,?)",
			o.Id,
			o.Title,
			o.Hosts,
			o.Counters,
			o.ScreenId,
			o.Timespan,
			o.GraphType,
			o.Method,
			o.Position,
			o.FalconTags,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func do_dashboard_screen() error {
	var o models.DashboardScreen
	do_sql(dst, "TRUNCATE TABLE dashboard_screen")
	set_auto_increment("dashboard_screen")

	rows, err := src.Query("select id, pid, name from dashboard_screen")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&o.Id,
			&o.Pid,
			&o.Name,
		)
		if err != nil {
			return err
		}
		err = do_sql(dst, "insert into dashboard_screen(id, pid, name) values(?,?,?)",
			o.Id,
			o.Pid,
			o.Name,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func do_tmp_graph() error {
	var o models.TmpGraph
	do_sql(dst, "TRUNCATE TABLE tmp_graph")
	set_auto_increment("tmp_graph")

	rows, err := src.Query("select id, endpoints, counters, ck name from tmp_graph")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&o.Id,
			&o.Endpoints,
			&o.Counters,
			&o.Ck,
		)
		if err != nil {
			return err
		}
		err = do_sql(dst, "insert into tmp_graph(id, endpoints, counters, ck) values(?,?,?,?)",
			o.Id,
			o.Endpoints,
			o.Counters,
			o.Ck,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func do_sql(db *sql.DB, s string, args ...interface{}) error {
	stmt, err := db.Prepare(s)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	return err
}

func set_auto_increment(table string) (err error) {
	var i int
	err = src.QueryRow("select AUTO_INCREMENT from information_schema.tables  where table_name = ?", table).Scan(&i)
	if err != nil {
		return
	}

	log.Printf("ALTER TABLE "+table+" = %d\n", i)
	return do_sql(dst, "ALTER TABLE "+table+" = ?", i)
}
