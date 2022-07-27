package model

import (
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	aliyunLogRecordFieldNames           = builder.RawFieldNames(&AliyunLogRecord{})
	aliyunLogRecordRows                 = strings.Join(aliyunLogRecordFieldNames, ",")
	aliyunLogRecordRowsExpectAutoSet    = strings.Join(stringx.Remove(aliyunLogRecordFieldNames, "`id`", "`create_time`", "`modify_time`"), ",")
	aliyunLogRecordRowsWithPlaceHolder  = strings.Join(stringx.Remove(aliyunLogRecordFieldNames, "`id`", "`create_time`"), "=?,") + "=?"
	aliyunLogRecordRowsWithPlaceHolder2 = strings.Join(stringx.Remove(aliyunLogRecordFieldNames, "`id`", "`create_time`", "`date`", "`project_name`", "`logstore_name`"), "=?,") + "=?"

	cacheAliyunLogRecordIdPrefix                          = "cache:aliyunLogRecord:id:"
	cacheAliyunLogRecordDateProjectNameLogstoreNamePrefix = "cache:aliyunLogRecord:date:projectName:logstoreName:"
)

type (
	AliyunLogRecordModel interface {
		Insert(session sqlx.Session, data *AliyunLogRecord) (sql.Result, error)
		FindOne(id int64) (*AliyunLogRecord, error)
		FindOneByDateProjectNameLogstoreName(date time.Time, projectName string, logstoreName string) (*AliyunLogRecord, error)
		Update(session sqlx.Session, data *AliyunLogRecord) error
		InsertOrUpdate(session sqlx.Session, data *AliyunLogRecord) error
		Delete(session sqlx.Session, id int64) error
		Trans(fn func(session sqlx.Session) error) error
	}

	defaultAliyunLogRecordModel struct {
		sqlc.CachedConn
		table string
	}

	AliyunLogRecord struct {
		Id           int64     `db:"id"`
		Date         time.Time `db:"date"`          // 统计时间
		ProjectName  string    `db:"project_name"`  // project名称
		LogstoreName string    `db:"logstore_name"` // 日志库名称
		Tp           int64     `db:"type"`          // 敏感信息类型：1. 日志存在电话号码
		Count        int64     `db:"count"`         // 敏感信息数量
		Info         string    `db:"info"`          // 单条敏感信息例子：如日志存在电话号码
		CreateTime   time.Time `db:"create_time"`   // 创建时间
		ModifyTime   time.Time `db:"modify_time"`   // 更新时间
	}
)

func NewAliyunLogRecordModel(conn sqlx.SqlConn, c cache.CacheConf) AliyunLogRecordModel {
	return &defaultAliyunLogRecordModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`aliyun_log_record`",
	}
}

func (m *defaultAliyunLogRecordModel) Insert(session sqlx.Session, data *AliyunLogRecord) (result sql.Result, err error) {

	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, aliyunLogRecordRowsExpectAutoSet)
	if session != nil {
		result, err = session.Exec(query, data.Date, data.ProjectName, data.LogstoreName, data.Tp, data.Count, data.Info)
	} else {
		result, err = m.ExecNoCache(query, data.Date, data.ProjectName, data.LogstoreName, data.Tp, data.Count, data.Info)
	}
	if err != nil {
		return nil, err
	}
	lastInsertId, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()
	if lastInsertId != 0 && rowsAffected != 0 {
		err = m.DelCache(m.formatPrimary(lastInsertId), m.formatDateProjectLogstoreKey(data.Date, data.ProjectName, data.LogstoreName))
		if err != nil {
			logx.Errorf("删除缓存失败，错误信息：%v", err)
		}
	}
	return result, err
}

func (m *defaultAliyunLogRecordModel) FindOne(id int64) (*AliyunLogRecord, error) {
	var resp AliyunLogRecord
	err := m.QueryRow(&resp, m.formatPrimary(id), func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", aliyunLogRecordRows, m.table)
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

func (m *defaultAliyunLogRecordModel) FindOneByDateProjectNameLogstoreName(date time.Time, projectName string, logstoreName string) (*AliyunLogRecord, error) {
	var resp AliyunLogRecord
	err := m.QueryRowIndex(&resp, m.formatDateProjectLogstoreKey(date, projectName, logstoreName), m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `date` = ? and `project_name` = ? and `logstore_name` = ? limit 1", aliyunLogRecordRows, m.table)
		if err := conn.QueryRow(&resp, query, date, projectName, logstoreName); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultAliyunLogRecordModel) Update(session sqlx.Session, data *AliyunLogRecord) error {
	var (
		result sql.Result
		err    error
	)

	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, aliyunLogRecordRowsWithPlaceHolder)

	if session != nil {
		result, err = session.Exec(query, data.Date, data.ProjectName, data.LogstoreName, data.Tp, data.Count, data.Info, data.ModifyTime, data.Id)
	} else {
		result, err = m.ExecNoCache(query, data.Date, data.ProjectName, data.LogstoreName, data.Tp, data.Count, data.Info, data.ModifyTime, data.Id)
	}

	if err != nil {
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()
	if lastInsertId != 0 && rowsAffected != 0 {
		err = m.DelCache(m.formatPrimary(lastInsertId), m.formatDateProjectLogstoreKey(data.Date, data.ProjectName, data.LogstoreName))
		if err != nil {
			logx.Errorf("删除缓存失败，错误信息：%v", err)
		}
	}
	return err
}

func (m *defaultAliyunLogRecordModel) InsertOrUpdate(session sqlx.Session, data *AliyunLogRecord) error {
	var (
		result sql.Result
		err    error
	)

	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?) on DUPLICATE KEY UPDATE %s", m.table, aliyunLogRecordRowsExpectAutoSet, aliyunLogRecordRowsWithPlaceHolder2)

	if session != nil {
		result, err = session.Exec(query, data.Date, data.ProjectName, data.LogstoreName, data.Tp, data.Count, data.Info, data.Tp, data.Count, data.Info, data.ModifyTime)
	} else {
		result, err = m.ExecNoCache(query, data.Date, data.ProjectName, data.LogstoreName, data.Tp, data.Count, data.Info, data.Tp, data.Count, data.Info, data.ModifyTime)
	}

	if err != nil {
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()
	if lastInsertId != 0 && rowsAffected != 0 {
		err = m.DelCache(m.formatPrimary(lastInsertId), m.formatDateProjectLogstoreKey(data.Date, data.ProjectName, data.LogstoreName))
		if err != nil {
			logx.Errorf("删除缓存失败，错误信息：%v", err)
		}
	}
	return err
}

func (m *defaultAliyunLogRecordModel) Delete(session sqlx.Session, id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		if session != nil {
			return session.Exec(query, id)
		}
		return conn.Exec(query, id)
	}, m.formatPrimary(id), m.formatDateProjectLogstoreKey(data.Date, data.ProjectName, data.LogstoreName))
	return err
}

func (m *defaultAliyunLogRecordModel) Trans(fn func(session sqlx.Session) error) error {
	err := m.Transact(func(session sqlx.Session) error {
		return fn(session)
	})
	return err
}

func (m *defaultAliyunLogRecordModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheAliyunLogRecordIdPrefix, primary)
}

func (m *defaultAliyunLogRecordModel) formatDateProjectLogstoreKey(date, project, logstore interface{}) string {
	return fmt.Sprintf("%s%v:%v:%v", cacheAliyunLogRecordDateProjectNameLogstoreNamePrefix, date, project, logstore)
}

func (m *defaultAliyunLogRecordModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", aliyunLogRecordRows, m.table)
	return conn.QueryRow(v, query, primary)
}
