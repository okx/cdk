package db

import (
	"context"
	dbSql "database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/0xPolygon/cdk-data-availability/rpc"
	jRPC "github.com/0xPolygon/cdk-rpc/rpc"
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
) *jRPC.Server {
	logger := log.WithFields("module", NAME)

	meter := otel.Meter(meterName)
	methodList := strings.Split(cfg.AuthMethodList, ",")
	log.Info(fmt.Sprintf("Sqlite service method auth list: %s", methodList))
	for _, s := range methodList {
		methodList = append(methodList, s)
	}
	log.Info(fmt.Sprintf("Sqlite service dbMaps: %v", dbMaps))
	time.Sleep(10 * time.Second)
	sqlDBs := make(map[string]*dbSql.DB)
	for k, dbPath := range dbMaps {
		log.Info(fmt.Sprintf("Sqlite service: %s, %s", k, dbPath))
		db, err := NewSQLiteDB(dbPath)
		if err != nil {
			log.Fatal(err)
		}
		sqlDBs[k] = db
	}

	services := []jRPC.Service{
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

	return jRPC.NewServer(jRPC.Config{
		Host:         cfg.Host,
		Port:         cfg.Port,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}, services, jRPC.WithLogger(logger.GetSugaredLogger()))
}

type A struct {
	Fields map[string]interface{}
}

func (b *SqliteEndpoints) Select(
	db string,
	sql string,
) (interface{}, rpc.Error) {
	err, dbCon := b.checkAndGetDB(db, sql, METHOD_SELECT)
	if err != nil {
		return zeroHex, rpc.NewRPCError(rpc.DefaultErrorCode, fmt.Sprintf("invalid sql: %s", err.Error()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), b.readTimeout)
	defer cancel()
	rows, err := dbCon.QueryContext(ctx, sql)
	if err != nil {
		if errors.Is(err, dbSql.ErrNoRows) {
			return nil, rpc.NewRPCError(
				rpc.DefaultErrorCode, fmt.Sprintf("No rows"), ErrNotFound)
		}
		return nil, rpc.NewRPCError(rpc.DefaultErrorCode, fmt.Sprintf("failed to query: %s", err.Error()))
	}
	err, result := getResults(rows)

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

func (b *SqliteEndpoints) GetDbList() (interface{}, rpc.Error) {
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

func (b *SqliteEndpoints) checkAndGetDB(db string, sql string, method string) (error, *dbSql.DB) {
	if len(sql) <= LIMIT_SQL_LEN {
		return fmt.Errorf("sql length is too short"), nil
	}

	sqlMethod := strings.ToLower(sql[:6])
	if sqlMethod != method {
		return fmt.Errorf("sql method is not valid"), nil
	}

	found := false
	for _, str := range b.authMethods {
		if str == method {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("sql method is not authorized"), nil
	}

	dbCon, ok := b.sqlDBs[db]
	if !ok {
		return fmt.Errorf("sql db is not valid"), nil
	}

	return nil, dbCon
}

func getResults(rows *dbSql.Rows) (error, []A) {
	var result []A

	for rows.Next() {
		var a A
		a.Fields = make(map[string]interface{})
		err := rows.Scan()
		if err != nil {
			log.Info("Error scanning row:", err)
			return err, nil
		}

		columns, err := rows.Columns()
		if err != nil {
			fmt.Println("Error getting columns:", err)
			return err, nil
		}

		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = &values[i]
		}

		err = rows.Scan(values...)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return err, nil
		}

		for i, col := range columns {
			a.Fields[col] = values[i]
		}

		result = append(result, a)
	}

	return nil, result
}
