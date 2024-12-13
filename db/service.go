package db

import (
	"context"
	dbSql "database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/0xPolygon/cdk-rpc/rpc"
	"github.com/0xPolygon/cdk/db/types"
	"github.com/0xPolygon/cdk/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	NAME      = "sqlite"
	meterName = "github.com/0xPolygon/cdk/sqlite/service"

	METHOD_SELECT = "select"
	METHOD_UPDATE = "update"
	METHOD_DELETE = "delete"

	LIMIT_SQL_LEN = 6

	zeroHex = "0x0"

	SEQS_L1_TREE  = "seqs_l1tree"
	SEQS_TX_MGR   = "seqs_txmgr"
	SEQS_REORG_L1 = "seqs_reorg_l1"

	AGG_SYNC     = "agg_sync"
	AGG_TX_MGR   = "agg_txmgr"
	AGG_REORG_L1 = "agg_reorg_l1"
)

type SqliteEndpoints struct {
	logger       *log.Logger
	meter        metric.Meter
	readTimeout  time.Duration
	writeTimeout time.Duration

	authMethods []string

	dbMaps map[string]string
	sqlDBs map[string]*dbSql.DB
}

func CreateSqliteService(
	cfg Config,
	dbMaps map[string]string,
) *rpc.Server {
	logger := log.WithFields("module", NAME)

	meter := otel.Meter(meterName)
	methodList := strings.Split(cfg.AuthMethodList, ",")
	log.Info(fmt.Sprintf("Sqlite service method auth list: %s", methodList))
	for _, s := range methodList {
		methodList = append(methodList, s)
	}
	log.Info(fmt.Sprintf("Sqlite service dbMaps: %v", dbMaps))
	sqlDBs := make(map[string]*dbSql.DB)
	for k, dbPath := range dbMaps {
		log.Info(fmt.Sprintf("Sqlite service: %s, %s", k, dbPath))
		db, err := NewSQLiteDB(dbPath)
		if err != nil {
			log.Fatal(err)
		}
		sqlDBs[k] = db
	}
	log.Info(fmt.Sprintf("Sqlite service sqlDBs sucess %v", sqlDBs))

	services := []rpc.Service{
		{
			Name: NAME,
			Service: &SqliteEndpoints{
				logger:       logger,
				meter:        meter,
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

type A struct {
	Fields map[string]interface{}
}

func (b *SqliteEndpoints) Select(
	db string,
	sql string,
) (interface{}, rpc.Error) {
	log.Info(fmt.Sprintf("zjg----1"))
	err, dbCon := b.checkAndGetDB(db, sql, METHOD_SELECT)
	if err != nil {
		log.Info(fmt.Sprintf("zjg----1-2---error"))
		return zeroHex, rpc.NewRPCError(rpc.DefaultErrorCode, fmt.Sprintf("check params invalid: %s", err.Error()))
	}
	log.Info(fmt.Sprintf("zjg----2"))
	ctx, cancel := context.WithTimeout(context.Background(), b.readTimeout)
	defer cancel()
	log.Info(fmt.Sprintf("zjg----2--1"))
	rows, err := dbCon.QueryContext(ctx, sql)
	if err != nil {
		if errors.Is(err, dbSql.ErrNoRows) {
			return nil, rpc.NewRPCError(
				rpc.DefaultErrorCode, fmt.Sprintf("No rows"), ErrNotFound)
		}
		return nil, rpc.NewRPCError(rpc.DefaultErrorCode, fmt.Sprintf("failed to query: %s", err.Error()))
	}
	log.Info(fmt.Sprintf("zjg----3"))

	err, result := getResults(rows)
	if err != nil {
		return nil, rpc.NewRPCError(rpc.DefaultErrorCode, fmt.Sprintf("failed to get results: %s", err.Error()))
	}
	log.Info(fmt.Sprintf("zjg----14"))

	return result, nil
}

func (b *SqliteEndpoints) Update(
	db string,
	sql string,
) (interface{}, rpc.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), b.readTimeout)
	defer cancel()
	c, merr := b.meter.Int64Counter("claim_proof")
	if merr != nil {
		b.logger.Warnf("failed to create claim_proof counter: %s", merr)
	}
	c.Add(ctx, 1)

	return types.SqliteData{
		ProofLocalExitRoot:  "ProofLocalExitRoot",
		ProofRollupExitRoot: "ProofRollupExitRoot",
		L1InfoTreeLeaf:      "L1InfoTreeLeaf",
	}, nil
}

func (b *SqliteEndpoints) Delete(
	db string,
	sql string,
) (interface{}, rpc.Error) {
	log.Info(fmt.Sprintf("Sqlite service Delete: %s, %s", db, sql))

	return types.SqliteData{
		ProofLocalExitRoot:  "ProofLocalExitRoot",
		ProofRollupExitRoot: "ProofRollupExitRoot",
		L1InfoTreeLeaf:      "L1InfoTreeLeaf",
	}, nil
}

func (b *SqliteEndpoints) GetDbs() (interface{}, rpc.Error) {
	//var dbList []string
	dbList := make(map[string][]string)
	for k, dbPath := range b.dbMaps {
		log.Info(fmt.Sprintf("Sqlite service: %s, %s", k, dbPath))

		sql := "SELECT name FROM sqlite_master WHERE type = 'table' ORDER BY name;"
		err, dbCon := b.checkAndGetDB(k, sql, METHOD_SELECT)
		if err != nil {
			log.Info(fmt.Sprintf("zjg----1-2---error"))
			return zeroHex, rpc.NewRPCError(rpc.DefaultErrorCode, fmt.Sprintf("%s", err.Error()))
		}
		log.Info(fmt.Sprintf("zjg----2"))
		ctx, cancel := context.WithTimeout(context.Background(), b.readTimeout)
		defer cancel()
		log.Info(fmt.Sprintf("zjg----2--1"))
		rows, err := dbCon.QueryContext(ctx, sql)
		if err != nil {
			if errors.Is(err, dbSql.ErrNoRows) {
				return nil, rpc.NewRPCError(
					rpc.DefaultErrorCode, fmt.Sprintf("No rows"), ErrNotFound)
			}
			return nil, rpc.NewRPCError(rpc.DefaultErrorCode, fmt.Sprintf("failed to query: %s", err.Error()))
		}
		log.Info(fmt.Sprintf("zjg----3"))

		err, result := getTables(rows)
		if err != nil {
			return nil, rpc.NewRPCError(rpc.DefaultErrorCode, fmt.Sprintf("failed to get results: %s", err.Error()))
		}

		dbList[k] = result
	}

	return dbList, nil
}

func (b *SqliteEndpoints) checkAndGetDB(db string, sql string, method string) (error, *dbSql.DB) {
	log.Info("zjg, chenckAndGetDB: -----------1")
	if len(sql) <= LIMIT_SQL_LEN {
		log.Info("zjg, chenckAndGetDB: -----------1--1")
		return fmt.Errorf("sql length is too short"), nil
	}
	log.Info("zjg, chenckAndGetDB: -----------2")

	sqlMethod := strings.ToLower(sql[:6])
	if sqlMethod != method {
		return fmt.Errorf("sql method is not valid"), nil
	}
	log.Info("zjg, chenckAndGetDB: -----------3")

	found := false
	for _, str := range b.authMethods {
		log.Info("zjg, chenckAndGetDB: -----------4")

		if str == method {
			found = true
			break
		}
	}
	if !found {
		log.Info("zjg, chenckAndGetDB: -----------5")
		return fmt.Errorf("sql method is not authorized"), nil
	}
	log.Info("zjg, chenckAndGetDB: -----------6")

	dbCon, ok := b.sqlDBs[db]
	log.Info("zjg, chenckAndGetDB: -----------7")
	if !ok {
		return fmt.Errorf("sql db is not valid"), nil
	}
	log.Info("zjg, chenckAndGetDB: -----------8")
	return nil, dbCon
}

func getResults(rows *dbSql.Rows) (error, []A) {
	var result []A

	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		log.Fatalf("Failed to get columns: %v", err)
	}
	for rows.Next() {
		record := A{Fields: make(map[string]interface{})}

		// 创建值和指针的切片
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// 扫描结果
		if err := rows.Scan(valuePtrs...); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}

		// 将每列数据存入 map
		for i, colName := range columns {
			record.Fields[colName] = values[i]
		}

		// 输出结果或处理
		fmt.Printf("Record: %+v\n", record.Fields)
		result = append(result, record)
	}

	return nil, result
}

func getTables(rows *dbSql.Rows) (error, []string) {
	var result []string

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return err, nil
		}
		log.Info(fmt.Sprintf("Table name: %s", tableName))
		result = append(result, tableName)
	}

	return nil, result
}
