
# SQLite JSON-RPC Operations
```json

[Sqlite]
Enabled = true
Host = "0.0.0.0"
Port = 8081
ReadTimeout = "2s"
WriteTimeout = "2s"
AuthMethodList = "select,insert,update,delete"
MaxRequestsPerIPAndSecond = 500

```

## Table of Contents

1. [SQLite JSON-RPC Methods](#sqlite-json-rpc-methods)
    - [sqlite_getDbs](#sqlite_getdbs)
    - [sqlite_select](#sqlite_select)
    - [sqlite_insert](#sqlite_insert)
    - [sqlite_update](#sqlite_update)
    - [sqlite_delete](#sqlite_delete)
2. [Example Requests](#example-requests)

## SQLite JSON-RPC Methods

The following JSON-RPC methods can be used to perform database operations:

### sqlite_getDbs

This method retrieves the list of databases managed by the SQLite transaction manager.

#### Request:
```json
{
    "id": 1,
    "jsonrpc": "2.0",
    "method": "sqlite_getDbs",
    "params": []
}
```

#### Response:
```json
{
    "jsonrpc": "2.0",
    "id": 1,
    "result": {
        "seqs_l1tree": [
            "block",
            "gorp_migrations",
            "l1_info_rht",
            "l1_info_root",
            "l1info_initial",
            "l1info_leaf",
            "rollup_exit_rht",
            "rollup_exit_root",
            "verify_batches"
        ],
        "seqs_reorg_l1": [
            "gorp_migrations",
            "tracked_block"
        ],
        "seqs_txmgr": [
            "gorp_migrations",
            "monitored_txs"
        ]
    },
    "error": null
}
```

---

### sqlite_select

This method executes a SQL `SELECT` query on the `monitored_txs` table to retrieve transaction data for a specific transaction ID.

#### Request:
```json
{
    "id": 1,
    "jsonrpc": "2.0",
    "method": "sqlite_select",
    "params": ["seqs_txmgr", "SELECT * FROM monitored_txs WHERE id='0x1c89ee4fbd62bcf844dacb5157df8fdf9036f210a1e3200d610ef1f778f2c32c';"]
}
```

#### Response:
```json
{
    "jsonrpc": "2.0",
    "id": 1,
    "result": [
        {
            "Fields": {
                "blob_gas": 0,
                "blob_gas_price": null,
                "blob_sidecar": null,
                "block_number": 891,
                "created_at": "2024-12-16T02:12:34Z",
                "estimate_gas": 0,
                "from_address": "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
                "gas": 224513,
                "gas_offset": 80000,
                "gas_price": "1",
                "gas_tip_cap": null,
                "history": null,
                "id": "0x1c89ee4fbd62bcf844dacb5157df8fdf9036f210a1e3200d610ef1f778f2c32c",
                "nonce": 117,
                "status": "finalized",
                "to_address": "0x1b173087729c88a47568AF87b17C653039377BA6",
                "tx_data": null,
                "updated_at": "2024-12-16T02:12:54Z",
                "value": "0"
            }
        }
    ],
    "error": null
}
```
---

### sqlite_insert

This method inserts a new record into the `monitored_txs` table. The transaction details such as `from_address`, `to_address`, `gas`, `block_number`, etc., are inserted into the table.

#### Request:
```json
{
    "id": 1,
    "jsonrpc": "2.0",
    "method": "sqlite_insert",
    "params": ["seqs_txmgr", "INSERT INTO monitored_txs (blob_gas, blob_gas_price, block_number, created_at, estimate_gas, from_address, gas, gas_offset, gas_price, gas_tip_cap, id, nonce, status, to_address, updated_at, value) VALUES (0, NULL, 891, '2024-12-16T02:12:34Z', 0, '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266', 224513, 80000, '1', NULL, '0x1c89ee4fbd62bcf844dacb5157df8fdf9036f210a1e3200d610ef1f778f2c32c', 117, 'finalized', '0x1b173087729c88a47568AF87b17C653039377BA6', '2024-12-16T02:12:54Z', '0');
}
```

#### Response:
```json
{
    "jsonrpc": "2.0",
    "id": 1,
    "result": {
        "RowsAffected": 1
    },
    "error": null
}
```

---

### sqlite_update

This method updates an existing record in the `monitored_txs` table. In this example, it updates the `to_address` field for a transaction with the specified ID.

#### Request:
```json
{
    "id": 1,
    "jsonrpc": "2.0",
    "method": "sqlite_update",
    "params": ["seqs_txmgr", "UPDATE monitored_txs SET to_address = '0x2b173087729c88a47568AF87b17C653039377BA6' WHERE id = '0x1c89ee4fbd62bcf844dacb5157df8fdf9036f210a1e3200d610ef1f778f2c32c';"]
}
```

#### Response:
```json
{
    "jsonrpc": "2.0",
    "id": 1,
    "result": {
        "RowsAffected": 1
    },
    "error": null
}
```

---

### sqlite_delete

This method deletes a record from the `monitored_txs` table. In this example, it deletes the transaction with the specified ID.

#### Request:
```json
{
    "id": 1,
    "jsonrpc": "2.0",
    "method": "sqlite_delete",
    "params": ["seqs_txmgr", "DELETE FROM monitored_txs WHERE id = '0x1c89ee4fbd62bcf844dacb5157df8fdf9036f210a1e3200d610ef1f778f2c32c';"]
}
```

#### Response:
```json
{
    "jsonrpc": "2.0",
    "id": 1,
    "result": {
        "RowsAffected": 1
    },
    "error": null
}
```

---
