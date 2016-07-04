/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package storage

import (
	"database/sql"
	"sync"

	"github.com/golang/glog"
)

var (
	/* db */
	appDb     *sql.DB
	dbLock    sync.RWMutex
	dbConnMap map[string]*sql.DB
)

func GetDbConn(connName, dsn string, maxIdle int) (c *sql.DB, e error) {
	dbLock.Lock()
	defer dbLock.Unlock()

	var err error
	var dbConn *sql.DB
	dbConn = dbConnMap[connName]
	if dbConn == nil {
		dbConn, err = makeDbConn(dsn, maxIdle)
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
func makeDbConn(dsn string, maxIdle int) (conn *sql.DB, err error) {
	conn, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	conn.SetMaxIdleConns(maxIdle)
	err = conn.Ping()

	return conn, err
}

func closeDbConn(conn *sql.DB) {
	if conn != nil {
		conn.Close()
	}
}

func dbStart(dsn string, maxIdle int) {
	var err error
	appDb, err = makeDbConn(dsn, maxIdle)
	if appDb == nil || err != nil {
		glog.Fatal("dbInit, get db conn fail", err)
	}

	dbConnMap = make(map[string]*sql.DB)
	glog.Info("dbInit ok")
}
