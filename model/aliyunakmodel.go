package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	aliyunAkFieldNames          = builder.RawFieldNames(&AliyunAk{})
	aliyunAkRows                = strings.Join(aliyunAkFieldNames, ",")
	aliyunAkRowsExpectAutoSet   = strings.Join(stringx.Remove(aliyunAkFieldNames, "`ID`", "`createTime`"), ",")
	aliyunAkRowsWithPlaceHolder = strings.Join(stringx.Remove(aliyunAkFieldNames, "`ID`", "`createTime`"), "=?,") + "=?"

	cacheAliyunAkIDPrefix = "cache:aliyunAk:id:"
)

type (
	AliyunAkModel interface {
		Insert(session sqlx.Session, data *AliyunAk) (sql.Result, error)
		FindOne(iD int64) (*AliyunAk, error)
		FindAll() ([]*AliyunAk, error)
		Update(session sqlx.Session, data *AliyunAk) error
		Delete(session sqlx.Session, iD int64) error
		Trans(fn func(session sqlx.Session) error) error
	}

	defaultAliyunAkModel struct {
		sqlc.CachedConn
		table string
	}

	AliyunAk struct {
		ID               int64          `db:"ID"`
		Account          sql.NullString `db:"Account"`          // 账号ID
		PrimaryAccount   sql.NullString `db:"PrimaryAccount"`   // 主账号
		SecondaryAccount sql.NullString `db:"SecondaryAccount"` // 子账号
		AccessKeyID      sql.NullString `db:"AccessKeyID"`      // AccessKeyID
		AccessKey        sql.NullString `db:"AccessKey"`        // AccessKey
		Permission       sql.NullString `db:"Permission"`       // 权限
		Info             sql.NullString `db:"Info"`             // 说明
		CreateTime       time.Time      `db:"CreateTime"`
	}
)

func NewAliyunAkModel(conn sqlx.SqlConn, c cache.CacheConf) AliyunAkModel {
	return &defaultAliyunAkModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`aliyun_ak`",
	}
}

func (m *defaultAliyunAkModel) Insert(session sqlx.Session, data *AliyunAk) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, aliyunAkRowsExpectAutoSet)
	if session != nil {
		return session.Exec(query, data.Account, data.PrimaryAccount, data.SecondaryAccount, data.AccessKeyID, data.AccessKey, data.Permission, data.Info)
	}
	return m.ExecNoCache(query, data.Account, data.PrimaryAccount, data.SecondaryAccount, data.AccessKeyID, data.AccessKey, data.Permission, data.Info)
}

func (m *defaultAliyunAkModel) FindOne(iD int64) (*AliyunAk, error) {
	aliyunAkIDKey := fmt.Sprintf("%s%v", cacheAliyunAkIDPrefix, iD)
	var resp AliyunAk
	err := m.QueryRow(&resp, aliyunAkIDKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `ID` = ? limit 1", aliyunAkRows, m.table)
		return conn.QueryRow(v, query, iD)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultAliyunAkModel) FindAll() ([]*AliyunAk, error)  {
	var resp []*AliyunAk
	query := fmt.Sprintf("select %s from %s", aliyunAkRows, m.table)
	err := m.QueryRowsNoCache(&resp, query)

	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultAliyunAkModel) Update(session sqlx.Session, data *AliyunAk) error {
	aliyunAkIDKey := fmt.Sprintf("%s%v", cacheAliyunAkIDPrefix, data.ID)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `ID` = ?", m.table, aliyunAkRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.Account, data.PrimaryAccount, data.SecondaryAccount, data.AccessKeyID, data.AccessKey, data.Permission, data.Info, data.ID)
		}
		return conn.Exec(query, data.Account, data.PrimaryAccount, data.SecondaryAccount, data.AccessKeyID, data.AccessKey, data.Permission, data.Info, data.ID)
	}, aliyunAkIDKey)
	return err
}

func (m *defaultAliyunAkModel) Delete(session sqlx.Session, iD int64) error {

	aliyunAkIDKey := fmt.Sprintf("%s%v", cacheAliyunAkIDPrefix, iD)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `ID` = ?", m.table)
		if session != nil {
			return session.Exec(query, iD)
		}
		return conn.Exec(query, iD)
	}, aliyunAkIDKey)
	return err
}

func (m *defaultAliyunAkModel) Trans(fn func(session sqlx.Session) error) error  {
	err := m.Transact(func(session sqlx.Session) error {
		return fn(session)
	})
	return err
}

func (m *defaultAliyunAkModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheAliyunAkIDPrefix, primary)
}

func (m *defaultAliyunAkModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `ID` = ? limit 1", aliyunAkRows, m.table)
	return conn.QueryRow(v, query, primary)
}
