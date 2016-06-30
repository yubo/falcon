/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package storage

import (
	"database/sql"
	"log"
)

func GetDbConn(connName string) (c *sql.DB, e error) {
	dbLock.Lock()
	defer dbLock.Unlock()

	var err error
	var dbConn *sql.DB
	dbConn = dbConnMap[connName]
	if dbConn == nil {
		dbConn, err = makeDbConn()
		if dbConn == nil || err != nil {
			closeDbConn(dbConn)
			return nil, err
		}
		dbConnMap[connName] = dbConn
	}

	err = dbConn.Ping()
	if err != nil {
		closeDbConn(dbConn)
		delete(dbConnMap, connName)
		return nil, err
	}

	return dbConn, err
}

// 创建一个新的mysql连接
func makeDbConn() (conn *sql.DB, err error) {
	conn, err = sql.Open("mysql", config().Dsn)
	if err != nil {
		return nil, err
	}

	conn.SetMaxIdleConns(config().DbMaxIdle)
	err = conn.Ping()

	return conn, err
}

func closeDbConn(conn *sql.DB) {
	if conn != nil {
		conn.Close()
	}
}

func dbInit() {
	var err error
	DB, err = makeDbConn()
	if DB == nil || err != nil {
		log.Fatalln("dbInit, get db conn fail", err)
	}

	dbConnMap = make(map[string]*sql.DB)
	log.Println("dbInit ok")
}
