package main

import (
	"fmt"

	cdkcommon "github.com/0xPolygon/cdk/common"
	"github.com/0xPolygon/cdk/config"
	sqlite "github.com/0xPolygon/cdk/db"
	"github.com/0xPolygon/cdk/db/types"
	"github.com/0xPolygon/cdk/log"
)

func runSqliteServiceIfNeeded(
	components []string,
	cfg config.Config,
) {
	dbPath := make(map[string]string)
	if isNeeded([]string{
		cdkcommon.AGGREGATOR},
		components) {
		dbPath[types.AggTxMgr] = cfg.Aggregator.EthTxManager.StoragePath
		dbPath[types.AggSync] = cfg.Aggregator.Synchronizer.SQLDB.DataSource
		dbPath[types.AggReorgL1] = cfg.ReorgDetectorL1.DBPath
	} else if isNeeded([]string{
		cdkcommon.SEQUENCE_SENDER},
		components) {
		dbPath[types.SeqsTxMgr] = cfg.SequenceSender.EthTxManager.StoragePath
		dbPath[types.SeqsL1Tree] = cfg.L1InfoTreeSync.DBPath
		dbPath[types.SeqsReorgL1] = cfg.ReorgDetectorL1.DBPath
	} else {
		log.Warn("No need to start sqlite service")
		return
	}

	allDBPath := ""
	for dbPathKey, dbPathValue := range dbPath {
		allDBPath += fmt.Sprintf("%s=%s \n", dbPathKey, dbPathValue)
	}

	server := sqlite.CreateSqliteService(cfg.Sqlite, dbPath)
	log.Info(fmt.Sprintf("Starting sqlite service on %s:%d,max:%v,\n%v",
		cfg.Sqlite.Host, cfg.Sqlite.Port, cfg.Sqlite.MaxRequestsPerIPAndSecond, allDBPath))
	go func() {
		if err := server.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	return
}
