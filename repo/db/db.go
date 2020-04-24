package db

import (
	"database/sql"
	"path"
	"strings"
	"sync"
    "fmt"

	"github.com/golang/protobuf/jsonpb"
	_ "github.com/mutecomm/go-sqlcipher"
	"github.com/SJTU-OpenNetwork/hon-textile-switch/repo"
)


var pbMarshaler = jsonpb.Marshaler{
	OrigName: true,
}

var pbUnmarshaler = jsonpb.Unmarshaler{
	AllowUnknownFields: true,
}

type SQLiteDatastore struct {
	streammetas		   repo.StreamMetaStore
	streamblocks	   repo.StreamBlockStore
	db                 *sql.DB
	lock               *sync.Mutex
}



func Create(repoPath, pin string) (*SQLiteDatastore, error) {
	dbPath := path.Join(repoPath, "datastore", "mainnet.db")
	fmt.Printf("Create datastore to %s\n", dbPath)
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	conn.SetMaxIdleConns(2)
	conn.SetMaxOpenConns(4)
	if pin != "" {
		p := "pragma key='" + strings.Replace(pin, "'", "''", -1) + "';"
		if _, err := conn.Exec(p); err != nil {
			return nil, err
		}
	}
	lock := new(sync.Mutex)
    //videoChunkLock := new(sync.Mutex)
    //videoLock := new(sync.Mutex)
	return &SQLiteDatastore{
		streammetas:		NewStreamMetaStore(conn, lock),
		streamblocks:       NewStreamBlockStore(conn, lock),
		db:                 conn,
		lock:               lock,
	}, nil
}

func (d *SQLiteDatastore) Ping() error {
	return d.db.Ping()
}

func (d *SQLiteDatastore) Close() {
	_ = d.db.Close()
}

func (d *SQLiteDatastore) StreamMetas() repo.StreamMetaStore {
	return d.streammetas
}

func (d *SQLiteDatastore) StreamBlocks() repo.StreamBlockStore {
	return d.streamblocks
}

func (d *SQLiteDatastore) Copy(dbPath string, pin string) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	var cp string
	stmt := "select name from sqlite_master where type='table'"
	rows, err := d.db.Query(stmt)
	if err != nil {
		return err
	}
	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		tables = append(tables, name)
	}
	if pin == "" {
		cp = `attach database '` + dbPath + `' as plaintext key '';`
		for _, name := range tables {
			cp = cp + "insert into plaintext." + name + " select * from main." + name + ";"
		}
	} else {
		cp = `attach database '` + dbPath + `' as encrypted key '` + pin + `';`
		for _, name := range tables {
			cp = cp + "insert into encrypted." + name + " select * from main." + name + ";"
		}
	}
	_, err = d.db.Exec(cp)
	if err != nil {
		return err
	}
	return nil
}

func (d *SQLiteDatastore) InitTables(pin string) error {
	return initDatabaseTables(d.db, pin)
}

func ConflictError(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}

func initDatabaseTables(db *sql.DB, pin string) error {
	var sqlStmt string
	if pin != "" {
		sqlStmt = "pragma key = '" + strings.Replace(pin, "'", "''", -1) + "';"
	}
	sqlStmt += `
	create table stream_metas (id text primary key not null, nstream integer, bitrate integer, caption text, nblocks integer, posterid text);
	
	create table stream_blocks (id text not null, streamid text , blockindex integer , blocksize integer , isroot integer, payload text, primary key(streamid, blockindex));
	create index stream_blocks_streamid on stream_blocks (streamid);
    `
	if _, err := db.Exec(sqlStmt); err != nil {
        fmt.Println(err)
		return err
	}
	return nil
}
