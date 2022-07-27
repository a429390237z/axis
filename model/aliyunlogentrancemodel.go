package model

import (
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	aliyunLogEntranceFieldNames          = builder.RawFieldNames(&AliyunLogEntrance{})
	aliyunLogEntranceRows                = strings.Join(aliyunLogEntranceFieldNames, ",")
	aliyunLogEntranceRowsExpectAutoSet   = strings.Join(stringx.Remove(aliyunLogEntranceFieldNames, "`id`"), ",")
	aliyunLogEntranceRowsWithPlaceHolder = strings.Join(stringx.Remove(aliyunLogEntranceFieldNames, "`id`"), "=?,") + "=?"

	cacheAliyunLogEntranceIdPrefix = "cache:aliyunLogEntrance:id:"
)

type (
	AliyunLogEntranceModel interface {
		Insert(session sqlx.Session, data *AliyunLogEntrance) (sql.Result, error)
		FindOne(id int64) (*AliyunLogEntrance, error)
		FindAll() ([]*AliyunLogEntrance, error)
		Update(session sqlx.Session, data *AliyunLogEntrance) error
		Delete(session sqlx.Session, id int64) error
		Trans(fn func(session sqlx.Session) error) error
	}

	defaultAliyunLogEntranceModel struct {
		sqlc.CachedConn
		table string
	}

	AliyunLogEntrance struct {
		Id               int64  `db:"id"`
		Region           string `db:"region"`            // 地域(英文）
		RegionCn         string `db:"region_cn"`         // 地域(中文)
		InternetEntrance string `db:"internet_entrance"` // 公网入口
		IntranetEntrance string `db:"intranet_entrance"` // 私网入口
	}
)

func NewAliyunLogEntranceModel(conn sqlx.SqlConn, c cache.CacheConf) AliyunLogEntranceModel {
	return &defaultAliyunLogEntranceModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`aliyun_log_entrance`",
	}
}

func (m *defaultAliyunLogEntranceModel) Insert(session sqlx.Session, data *AliyunLogEntrance) (result sql.Result, err error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, aliyunLogEntranceRowsExpectAutoSet)
	if session != nil {
		result, err = session.Exec(query, data.Region, data.RegionCn, data.InternetEntrance, data.IntranetEntrance)
	} else {
		result, err = m.ExecNoCache(query, data.Region, data.RegionCn, data.InternetEntrance, data.IntranetEntrance)
	}
	if err != nil {
		return
	}
	lastInsertId, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()
	if lastInsertId != 0 && rowsAffected != 0 {
		err = m.DelCache(m.formatPrimary(lastInsertId))
		if err != nil {
			logx.Errorf("删除缓存失败，错误信息：%v", err)
		}
	}
	return
}

func (m *defaultAliyunLogEntranceModel) FindOne(id int64) (*AliyunLogEntrance, error) {
	aliyunLogEntranceIdKey := fmt.Sprintf("%s%v", cacheAliyunLogEntranceIdPrefix, id)
	var resp AliyunLogEntrance
	err := m.QueryRow(&resp, aliyunLogEntranceIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", aliyunLogEntranceRows, m.table)
		return conn.QueryRow(v, query, id)
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

func (m *defaultAliyunLogEntranceModel) FindAll() ([]*AliyunLogEntrance, error) {
	var resp []*AliyunLogEntrance
	query := fmt.Sprintf("select %s from %s", aliyunLogEntranceRows, m.table)
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

func (m *defaultAliyunLogEntranceModel) Update(session sqlx.Session, data *AliyunLogEntrance) error {
	var (
		result sql.Result
		err    error
	)

	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, aliyunLogEntranceRowsWithPlaceHolder)
	if session != nil {
		result, err = session.Exec(query, data.Region, data.RegionCn, data.InternetEntrance, data.IntranetEntrance, data.Id)
	} else {
		result, err = m.ExecNoCache(query, data.Region, data.RegionCn, data.InternetEntrance, data.IntranetEntrance, data.Id)
	}
	if err != nil {
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()
	if lastInsertId != 0 && rowsAffected != 0 {
		err = m.DelCache(m.formatPrimary(lastInsertId))
		if err != nil {
			logx.Errorf("删除缓存失败，错误信息：%v", err)
		}
	}
	return err
}

func (m *defaultAliyunLogEntranceModel) Delete(session sqlx.Session, id int64) error {

	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		if session != nil {
			return session.Exec(query, id)
		}
		return conn.Exec(query, id)
	}, m.formatPrimary(id))
	return err
}

func (m *defaultAliyunLogEntranceModel) Trans(fn func(session sqlx.Session) error) error {
	err := m.Transact(func(session sqlx.Session) error {
		return fn(session)
	})
	return err
}

func (m *defaultAliyunLogEntranceModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheAliyunLogEntranceIdPrefix, primary)
}

func (m *defaultAliyunLogEntranceModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", aliyunLogEntranceRows, m.table)
	return conn.QueryRow(v, query, primary)
}
