package goutils

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/georgysavva/scany/sqlscan"
)

// NewMySQLConn 打开 MySQL 连接。
func NewMySQLConn(host string, port int, username, passwd, database, charset string) (*sql.DB, error) {
	return sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True", username, passwd, host, port, database, charset))
}

type DbSchema struct {
	CatalogName             string         `db:"CATALOG_NAME"`
	SchemaName              string         `db:"SCHEMA_NAME"`
	DefaultCharacterSetName string         `db:"DEFAULT_CHARACTER_SET_NAME"`
	DefaultCollationName    string         `db:"DEFAULT_COLLATION_NAME"`
	SqlPath                 sql.NullString `db:"SQL_PATH"`
}

type DbTableSchema struct {
	TableCatalog   string         `db:"TABLE_CATALOG"`
	TableSchema    string         `db:"TABLE_SCHEMA"`
	TableName      string         `db:"TABLE_NAME"`
	TableType      string         `db:"TABLE_TYPE"`
	Engine         sql.NullString `db:"ENGINE"`
	Version        sql.NullInt64  `db:"VERSION"`
	RowFormat      sql.NullString `db:"ROW_FORMAT"`
	TableRows      sql.NullInt64  `db:"TABLE_ROWS"`
	AvgRowLength   sql.NullInt64  `db:"AVG_ROW_LENGTH"`
	DataLength     sql.NullInt64  `db:"DATA_LENGTH"`
	MaxDataLength  sql.NullInt64  `db:"MAX_DATA_LENGTH"`
	IndexLength    sql.NullInt64  `db:"INDEX_LENGTH"`
	DataFree       sql.NullInt64  `db:"DATA_FREE"`
	AutoIncrement  sql.NullInt64  `db:"AUTO_INCREMENT"`
	CreateTime     sql.NullString `db:"CREATE_TIME"`
	UpdateTime     sql.NullString `db:"UPDATE_TIME"`
	CheckTime      sql.NullString `db:"CHECK_TIME"`
	TableCollation sql.NullString `db:"TABLE_COLLATION"`
	Checksum       sql.NullInt64  `db:"CHECKSUM"`
	CreateOptions  sql.NullString `db:"CREATE_OPTIONS"`
	TableComment   string         `db:"TABLE_COMMENT"`
}

type DbColumnSchema struct {
	TableCatalog           string         `db:"TABLE_CATALOG"`
	TableSchema            string         `db:"TABLE_SCHEMA"`
	TableName              string         `db:"TABLE_NAME"`
	ColumnName             string         `db:"COLUMN_NAME"`
	OrdinalPosition        uint64         `db:"ORDINAL_POSITION"`
	ColumnDefault          sql.NullString `db:"COLUMN_DEFAULT"`
	IsNullable             string         `db:"IS_NULLABLE"`
	DataType               string         `db:"DATA_TYPE"`
	CharacterMaximumLength sql.NullInt64  `db:"CHARACTER_MAXIMUM_LENGTH"`
	CharacterOctetLength   sql.NullInt64  `db:"CHARACTER_OCTET_LENGTH"`
	NumericPrecision       sql.NullInt64  `db:"NUMERIC_PRECISION"`
	NumericScale           sql.NullInt64  `db:"NUMERIC_SCALE"`
	DatetimePrecision      sql.NullInt64  `db:"DATETIME_PRECISION"`
	CharacterSetName       sql.NullString `db:"CHARACTER_SET_NAME"`
	CollationName          sql.NullString `db:"COLLATION_NAME"`
	ColumnType             string         `db:"COLUMN_TYPE"`
	ColumnKey              string         `db:"COLUMN_KEY"`
	Extra                  string         `db:"EXTRA"`
	Privileges             string         `db:"PRIVILEGES"`
	ColumnComment          string         `db:"COLUMN_COMMENT"`
	GenerationExpression   string         `db:"GENERATION_EXPRESSION"`
}

func GetDatabaseSchemas(db *sql.DB, database string) ([]*DbSchema, error) {
	ctx := context.Background()

	var ds []*DbSchema
	var err error

	if database == "" {
		err = sqlscan.Select(ctx, db, &ds, "SELECT `CATALOG_NAME`, `SCHEMA_NAME`, `DEFAULT_CHARACTER_SET_NAME`, `DEFAULT_COLLATION_NAME`, `SQL_PATH` FROM `information_schema`.`SCHEMATA` WHERE `SCHEMA_NAME` NOT IN ('information_schema', 'mysql', 'sys')")
	} else {
		err = sqlscan.Select(ctx, db, &ds, "SELECT `CATALOG_NAME`, `SCHEMA_NAME`, `DEFAULT_CHARACTER_SET_NAME`, `DEFAULT_COLLATION_NAME`, `SQL_PATH` FROM `information_schema`.`SCHEMATA` WHERE `SCHEMA_NAME` = ?", database)
	}

	if err != nil {
		return nil, err
	}

	return ds, nil
}

func GetTableSchema(db *sql.DB, database, tableName string) (*DbTableSchema, error) {
	ctx := context.Background()

	ds := &DbTableSchema{}

	err := sqlscan.Get(ctx, db, ds, "SELECT `TABLE_CATALOG`, `TABLE_SCHEMA`, `TABLE_NAME`, `TABLE_TYPE`, `ENGINE`, `VERSION`, `ROW_FORMAT`, `TABLE_ROWS`, `AVG_ROW_LENGTH`, `DATA_LENGTH`, `MAX_DATA_LENGTH`, `INDEX_LENGTH`, `DATA_FREE`, `AUTO_INCREMENT`, `CREATE_TIME`, `UPDATE_TIME`, `CHECK_TIME`, `TABLE_COLLATION`, `CHECKSUM`, `CREATE_OPTIONS`, `TABLE_COMMENT` FROM `information_schema`.`TABLES` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?", database, tableName)

	if err != nil {
		return nil, err
	}

	return ds, nil
}

func GetTableSchemas(db *sql.DB, database string) ([]*DbTableSchema, error) {
	ctx := context.Background()

	var ds []*DbTableSchema

	err := sqlscan.Select(ctx, db, &ds, "SELECT `TABLE_CATALOG`, `TABLE_SCHEMA`, `TABLE_NAME`, `TABLE_TYPE`, `ENGINE`, `VERSION`, `ROW_FORMAT`, `TABLE_ROWS`, `AVG_ROW_LENGTH`, `DATA_LENGTH`, `MAX_DATA_LENGTH`, `INDEX_LENGTH`, `DATA_FREE`, `AUTO_INCREMENT`, `CREATE_TIME`, `UPDATE_TIME`, `CHECK_TIME`, `TABLE_COLLATION`, `CHECKSUM`, `CREATE_OPTIONS`, `TABLE_COMMENT` FROM `information_schema`.`TABLES` WHERE `TABLE_SCHEMA` = ?", database)

	if err != nil {
		return nil, err
	}

	return ds, nil
}

func GetColumnSchemas(db *sql.DB, database, tableName string) ([]*DbColumnSchema, error) {
	ctx := context.Background()

	var ds []*DbColumnSchema

	err := sqlscan.Select(ctx, db, &ds, "SELECT `TABLE_CATALOG`, `TABLE_SCHEMA`, `TABLE_NAME`, `COLUMN_NAME`, `ORDINAL_POSITION`, `COLUMN_DEFAULT`, `IS_NULLABLE`, `DATA_TYPE`, `CHARACTER_MAXIMUM_LENGTH`, `CHARACTER_OCTET_LENGTH`, `NUMERIC_PRECISION`, `NUMERIC_SCALE`, `DATETIME_PRECISION`, `CHARACTER_SET_NAME`, `COLLATION_NAME`, `COLUMN_TYPE`, `COLUMN_KEY`, `EXTRA`, `PRIVILEGES`, `COLUMN_COMMENT`, `GENERATION_EXPRESSION` FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ? ORDER BY `ORDINAL_POSITION` ASC", database, tableName)

	if err != nil {
		return nil, err
	}

	return ds, nil
}
