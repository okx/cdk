package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/0xPolygon/cdk-rpc/rpc"
	"github.com/0xPolygon/cdk/db/types"
	"github.com/0xPolygon/cdk/log"
)

type SqliteEndpoints struct {
	logger       *log.Logger
	readTimeout  time.Duration
	writeTimeout time.Duration

	authMethods []string

	dbMaps map[string]string
	sqlDBs map[string]*sql.DB
}

func CreateSqliteService(
	cfg Config,
	dbMaps map[string]string,
) *rpc.Server {
	logger := log.WithFields("module", types.NAME)

	methodList := strings.Split(cfg.AuthMethodList, ",")
	log.Info(fmt.Sprintf("Sqlite service method auth list: %s", methodList))
	for _, s := range methodList {
		methodList = append(methodList, s)
	}
	sqlDBs := make(map[string]*sql.DB)
	for k, dbPath := range dbMaps {
		log.Info(fmt.Sprintf("Sqlite service load db: %s, %s", k, dbPath))
		db, err := NewSQLiteDB(dbPath)
		if err != nil {
			log.Fatal(err)
		}
		sqlDBs[k] = db
	}

	services := []rpc.Service{
		{
			Name: types.NAME,
			Service: &SqliteEndpoints{
				logger:       logger,
				readTimeout:  cfg.ReadTimeout.Duration,
				writeTimeout: cfg.WriteTimeout.Duration,
				authMethods:  methodList,
				dbMaps:       dbMaps,
				sqlDBs:       sqlDBs,
			},
		},
	}

	return rpc.NewServer(rpc.Config{
		Host:                      cfg.Host,
		Port:                      cfg.Port,
		ReadTimeout:               cfg.ReadTimeout,
		WriteTimeout:              cfg.WriteTimeout,
		MaxRequestsPerIPAndSecond: cfg.MaxRequestsPerIPAndSecond,
	}, services, rpc.WithLogger(logger.GetSugaredLogger()))
}

func (b *SqliteEndpoints) Select(
	dbName string,
	sqlCmd string,
) (interface{}, rpc.Error) {
	dbCon, err := b.checkAndGetDB(dbName, sqlCmd, types.MethodSelect)
	if err != nil {
		return types.ZeroHex, rpc.NewRPCError(rpc.DefaultErrorCode, fmt.Sprintf("check params invalid: %s", err.Error()))
	}
	ctx, cancel := context.WithTimeout(context.Background(), b.readTimeout)
	defer cancel()
	rows, err := dbCon.QueryContext(ctx, sqlCmd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, rpc.NewRPCError(
				rpc.DefaultErrorCode, fmt.Sprintf("No rows"), ErrNotFound)
		}
		return nil, rpc.NewRPCError(rpc.DefaultErrorCode, fmt.Sprintf("failed to query: %s", err.Error()))
	}

	result, err := getResults(rows)
	if err != nil {
		return nil, rpc.NewRPCError(rpc.DefaultErrorCode, fmt.Sprintf("failed to get results: %s", err.Error()))
	}

	return result, nil
}

func (b *SqliteEndpoints) Insert(
	dbName string,
	sqlCmd string,
) (interface{}, rpc.Error) {
	return b.alterMethod(types.MethodInsert, dbName, sqlCmd)
}

func (b *SqliteEndpoints) Delete(
	dbName string,
	sqlCmd string,
) (interface{}, rpc.Error) {
	return b.alterMethod(types.MethodDelete, dbName, sqlCmd)
}

func (b *SqliteEndpoints) Update(
	dbName string,
	sqlCmd string,
) (interface{}, rpc.Error) {
	return b.alterMethod(types.MethodUpdate, dbName, sqlCmd)
}

func (b *SqliteEndpoints) alterMethod(
	method string,
	db string,
	sql string,
) (interface{}, rpc.Error) {
	log.Info(fmt.Sprintf("Sqlite: %s, %s", db, sql))
	dbCon, err := b.checkAndGetDB(db, sql, method)
	if err != nil {
		return types.ZeroHex, rpc.NewRPCError(rpc.DefaultErrorCode, fmt.Sprintf("check params invalid: %s", err.Error()))
	}
	ctx, cancel := context.WithTimeout(context.Background(), b.readTimeout)
	defer cancel()

	tx, err := NewTx(ctx, dbCon)
	if err != nil {
		log.Error(fmt.Sprintf("failed to create tx: %s", err))
		return types.ZeroHex, rpc.NewRPCError(rpc.DefaultErrorCode, fmt.Sprintf("failed to create tx: %s", err))
	}
	shouldRollback := true
	defer func() {
		if shouldRollback {
			if errRollback := tx.Rollback(); errRollback != nil {
				log.Errorf("error while rolling back tx %v", errRollback)
			}
		}
	}()

	ct, err := tx.Exec(sql)
	if err != nil {
		log.Error(fmt.Sprintf("failed to exec: %s", err))
		return types.ZeroHex, rpc.NewRPCError(rpc.DefaultErrorCode, fmt.Sprintf("failed to exec: %s", err))
	}
	count, err := ct.RowsAffected()
	if err != nil {
		log.Error(fmt.Sprintf("failed to get rows affected: %s", err))
		return types.ZeroHex, rpc.NewRPCError(rpc.DefaultErrorCode, fmt.Sprintf("failed to get rows affected: %s", err))
	}

	if err := tx.Commit(); err != nil {
		log.Error(fmt.Sprintf("failed to commit: %s", err))
		return types.ZeroHex, rpc.NewRPCError(rpc.DefaultErrorCode, fmt.Sprintf("failed to commit: %s", err))
	}
	shouldRollback = false

	return types.SqliteData{
		RowsAffected: count,
	}, nil
}

func (b *SqliteEndpoints) GetDbs() (interface{}, rpc.Error) {
	dbList := make(map[string][]string)
	ctx, cancel := context.WithTimeout(context.Background(), b.readTimeout)
	defer cancel()
	for k, dbPath := range b.dbMaps {
		log.Info(fmt.Sprintf("Sqlite service: %s, %s", k, dbPath))
		sqlCmd := "SELECT name FROM sqlite_master WHERE type = 'table' ORDER BY name;"
		dbCon, err := b.checkAndGetDB(k, sqlCmd, types.MethodSelect)
		if err != nil {
			return types.ZeroHex, rpc.NewRPCError(rpc.DefaultErrorCode, fmt.Sprintf("%s", err.Error()))
		}
		rows, err := dbCon.QueryContext(ctx, sqlCmd)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, rpc.NewRPCError(
					rpc.DefaultErrorCode, fmt.Sprintf("No rows"), ErrNotFound)
			}
			return nil, rpc.NewRPCError(rpc.DefaultErrorCode, fmt.Sprintf("failed to query: %s", err.Error()))
		}

		result, err := getTables(rows)
		if err != nil {
			return nil, rpc.NewRPCError(rpc.DefaultErrorCode, fmt.Sprintf("failed to get results: %s", err.Error()))
		}
		dbList[k] = result
	}

	return dbList, nil
}

func (b *SqliteEndpoints) checkAndGetDB(db string, sql string, method string) (*sql.DB, error) {
	log.Info(fmt.Sprintf("Sqlite check db:%v, sql:%v, method:%v", db, sql, method))
	if len(sql) <= types.LimitSQLLen {
		return nil, fmt.Errorf("sql length is too short")
	}
	sqlMethod := strings.ToLower(sql[:6])
	if sqlMethod != method {
		return nil, fmt.Errorf("sql method is not valid")
	}
	found := false
	for _, str := range b.authMethods {
		if str == method {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("sql method is not authorized")
	}
	dbCon, ok := b.sqlDBs[db]
	if !ok {
		return nil, fmt.Errorf("sql db is not valid")
	}
	return dbCon, nil
}

func getResults(rows *sql.Rows) ([]types.QueryData, error) {
	var result []types.QueryData
	columns, err := rows.Columns()
	if err != nil {
		log.Error(fmt.Sprintf("Failed to get columns: %v", err))
		return nil, err
	}
	for rows.Next() {
		record := types.QueryData{Fields: make(map[string]interface{})}
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			log.Error("Failed to scan row: %v", err)
			return nil, err
		}
		for i, colName := range columns {
			record.Fields[colName] = values[i]
		}
		result = append(result, record)
	}
	return result, nil
}

func getTables(rows *sql.Rows) ([]string, error) {
	var result []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		log.Info(fmt.Sprintf("Table name: %s", tableName))
		result = append(result, tableName)
	}
	return result, nil
}
