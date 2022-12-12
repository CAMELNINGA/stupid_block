package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	INSERT_BLOCK = `INSERT INTO block (data,hash,previous_hash,timestamp,pow, sign, sign_ts) values ($1, $2, $3, $4,$5, $6, $7)`
	GET_BLOCK    = `SELECT data,hash,previous_hash,timestamp,pow,sign, sign_ts FROM block ORDER BY timestamp `
)

// Postgres -.
type Postgres struct {
	Config *Db
	Db     *sqlx.DB
}

type Db struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string

	Collector       string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifeTime int64
}

func (c *Db) ConnectionString() string {
	uri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		c.Host, c.Port,
		c.User, c.Name,
		c.Password,
	)

	return uri
}

// New -.
func New(cfg *Db) (*Postgres, error) {
	pgdb := &Postgres{
		Config: cfg,
	}

	db, err := sqlx.Open("pgx", cfg.ConnectionString())
	if err != nil {
		return nil, err
	}

	pgdb.Db = db

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifeTime) * time.Minute)

	return pgdb, nil
}

func (p *Postgres) Close() {
	if p.Db != nil {
		p.Db.Close()
	}
}

func (s Postgres) SetBlock(ctx context.Context, data, hash, previous_hash string, timestamp time.Time, pow int, sign string, signTs time.Time) error {
	if _, err := s.Db.ExecContext(ctx, INSERT_BLOCK, data, hash, previous_hash, timestamp, pow, sign, signTs); err != nil {

		return err
	}

	return nil
}

type DBBlock struct {
	Data         string    `db:"data"`
	Hash         string    `db:"hash"`
	PreviousHash string    `db:"previous_hash"`
	Timestamp    time.Time `db:"timestamp"`
	Pow          int       `db:"pow"`
	Sign         string    `db:"sign"`
	SignTs       time.Time `db:"sign_ts"`
}

func (d DBBlock) toBlock() (Block, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(d.Data), &data)
	if err != nil {
		return Block{}, err
	}
	return Block{
		data:         data,
		hash:         d.Hash,
		previousHash: d.PreviousHash,
		timestamp:    d.Timestamp,
		pow:          d.Pow,
	}, nil
}

type DbBlocks []DBBlock

func (dd DbBlocks) toBlocks() ([]Block, error) {
	var bloks []Block
	for _, v := range dd {
		s, err := v.toBlock()
		if err != nil {
			return []Block{}, err
		}
		bloks = append(bloks, s)
	}
	return bloks, nil
}

func (s Postgres) GetBlock(ctx context.Context) ([]Block, error) {
	var bloks DbBlocks

	err := s.Db.SelectContext(ctx, &bloks, GET_BLOCK)
	if err != nil {
		return []Block{}, err
	}
	return bloks.toBlocks()
}
