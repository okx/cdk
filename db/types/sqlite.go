package types

type SqliteData struct {
	RowsAffected int64
}

const (
	NAME = "sqlite"

	MethodSelect = "select"
	MethodInsert = "insert"
	MethodUpdate = "update"
	MethodDelete = "delete"

	LimitSqlLen = 6

	ZeroHex = "0x0"

	SeqsL1Tree  = "seqs_l1tree"
	SeqsTxMgr   = "seqs_txmgr"
	SeqsReorgL1 = "seqs_reorg_l1"

	AggSync    = "agg_sync"
	AggTxMgr   = "agg_txmgr"
	AggReorgL1 = "agg_reorg_l1"
)

type QueryData struct {
	Fields map[string]interface{}
}
