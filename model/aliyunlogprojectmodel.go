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
	aliyunLogProjectFieldNames           = builder.RawFieldNames(&AliyunLogProject{})
	aliyunLogProjectRows                 = strings.Join(aliyunLogProjectFieldNames, ",")
	aliyunLogProjectRowsExpectAutoSet    = strings.Join(stringx.Remove(aliyunLogProjectFieldNames, "`id`"), ",")
	aliyunLogProjectRowsWithPlaceHolder  = strings.Join(stringx.Remove(aliyunLogProjectFieldNames, "`id`", "`create_time`"), "=?,") + "=?"
	aliyunLogProjectRowsWithPlaceHolder2 = strings.Join(stringx.Remove(aliyunLogProjectFieldNames, "`id`", "`name`", "`region`", "`create_time`", "`endpoint`"), "=?,") + "=?"

	cacheAliyunLogProjectIdPrefix         = "cache:aliyunLogProject:id:"
	cacheAliyunLogProjectNameRegionPrefix = "cache:aliyunLogProject:owner:region:name:"
)

type (
	AliyunLogProjectModel interface {
		Insert(session sqlx.Session, data *AliyunLogProject) (sql.Result, error)
		FindOne(id int64) (*AliyunLogProject, error)
		FindOneByOwnerNameRegion(owner string, region string, name string) (*AliyunLogProject, error)
		Update(session sqlx.Session, data *AliyunLogProject) error
		InsertOrUpdate(session sqlx.Session, data *AliyunLogProject) error
		Delete(session sqlx.Session, id int64) error
		Trans(fn func(session sqlx.Session) error) error
	}

	defaultAliyunLogProjectModel struct {
		sqlc.CachedConn
		table string
	}

	AliyunLogProject struct {
		Id             int64     `db:"id"`
		Name           string    `db:"name"`             // project名称：作为Host的一部分，project名称在阿里云地域内全局唯一,创建后不可修改
		Description    string    `db:"description"`      // Project描述
		Region         string    `db:"region"`           // project所有地域
		Status         int64     `db:"status"`           // project状态：1：Normal(正常）0：Disable(禁用)
		Owner          string    `db:"owner"`            // 日志项目拥有者
		Maintainer     string    `db:"maintainer"`       // 日志项目维护者
		CreateTime     time.Time `db:"create_time"`      // project创建时间
		LastModifyTime time.Time `db:"last_modify_time"` // 最后一次更新project时间
		Endpoint       string    `db:"endpoint"`         // 地域
	}
)

func NewAliyunLogProjectModel(conn sqlx.SqlConn, c cache.CacheConf) AliyunLogProjectModel {
	return &defaultAliyunLogProjectModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`aliyun_log_project`",
	}
}

func (m *defaultAliyunLogProjectModel) Insert(session sqlx.Session, data *AliyunLogProject) (result sql.Result, err error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, aliyunLogProjectRowsExpectAutoSet)

	if session != nil {
		result, err = session.Exec(query, data.Name, data.Description, data.Region, data.Status, data.Owner, data.Maintainer, data.CreateTime, data.LastModifyTime, data.Endpoint)
	} else {
		result, err = m.ExecNoCache(query, data.Name, data.Description, data.Region, data.Status, data.Owner, data.Maintainer, data.CreateTime, data.LastModifyTime, data.Endpoint)
	}
	if err != nil {
		return nil, err
	}
	lastInsertId, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()
	if lastInsertId != 0 && rowsAffected != 0 {
		err = m.DelCache(m.formatPrimary(lastInsertId), m.formatProjectNameRegionKey(data.Owner, data.Region, data.Name))
		if err != nil {
			logx.Errorf("删除缓存失败，错误信息：%v", err)
		}
	}
	return result, err
}

func (m *defaultAliyunLogProjectModel) FindOne(id int64) (*AliyunLogProject, error) {
	var resp AliyunLogProject
	err := m.QueryRow(&resp, m.formatPrimary(id), func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", aliyunLogProjectRows, m.table)
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

func (m *defaultAliyunLogProjectModel) FindOneByOwnerNameRegion(owner string, region string, name string) (*AliyunLogProject, error) {
	var resp AliyunLogProject
	err := m.QueryRowIndex(&resp, m.formatProjectNameRegionKey(owner, region, name), m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `name` = ? and `region` = ? limit 1", aliyunLogProjectRows, m.table)
		if err := conn.QueryRow(&resp, query, name, region); err != nil {
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

func (m *defaultAliyunLogProjectModel) Update(session sqlx.Session, data *AliyunLogProject) error {
	var (
		result sql.Result
		err    error
	)

	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, aliyunLogProjectRowsWithPlaceHolder)
	if session != nil {
		result, err = session.Exec(query, data.Name, data.Description, data.Region, data.Status, data.Owner, data.Maintainer, data.LastModifyTime, data.Endpoint, data.Id)
	} else {
		result, err = m.ExecNoCache(query, data.Name, data.Description, data.Region, data.Status, data.Owner, data.Maintainer, data.LastModifyTime, data.Endpoint, data.Id)
	}
	if err != nil {
		return err
	}

	lastInsertId, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()
	if lastInsertId != 0 && rowsAffected != 0 {
		err = m.DelCache(m.formatPrimary(lastInsertId), m.formatProjectNameRegionKey(data.Owner, data.Region, data.Name))
		if err != nil {
			logx.Errorf("删除缓存失败，错误信息：%v", err)
		}
	}
	return err
}

func (m *defaultAliyunLogProjectModel) InsertOrUpdate(session sqlx.Session, data *AliyunLogProject) error {
	var (
		result sql.Result
		err    error
	)

	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?) on DUPLICATE KEY UPDATE %s", m.table, aliyunLogProjectRowsExpectAutoSet, aliyunLogProjectRowsWithPlaceHolder2)

	if session != nil {
		result, err = session.Exec(query, data.Name, data.Description, data.Region, data.Status, data.Owner, data.Maintainer, data.CreateTime, data.LastModifyTime, data.Endpoint, data.Description, data.Status, data.Owner, data.Maintainer, data.LastModifyTime)
	} else {
		result, err = m.ExecNoCache(query, data.Name, data.Description, data.Region, data.Status, data.Owner, data.Maintainer, data.CreateTime, data.LastModifyTime, data.Endpoint, data.Description, data.Status, data.Owner, data.Maintainer, data.LastModifyTime)
	}
	if err != nil {
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()
	if lastInsertId != 0 && rowsAffected != 0 {
		err = m.DelCache(m.formatPrimary(lastInsertId), m.formatProjectNameRegionKey(data.Owner, data.Region, data.Name))
		if err != nil {
			logx.Errorf("删除缓存失败，错误信息：%v", err)
		}
	}
	return err
}

func (m *defaultAliyunLogProjectModel) Delete(session sqlx.Session, id int64) error {
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
	}, m.formatPrimary(id), m.formatProjectNameRegionKey(data.Owner, data.Region, data.Name))
	return err
}

func (m *defaultAliyunLogProjectModel) Trans(fn func(session sqlx.Session) error) error {
	err := m.Transact(func(session sqlx.Session) error {
		return fn(session)
	})
	return err
}

func (m *defaultAliyunLogProjectModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheAliyunLogProjectIdPrefix, primary)
}

func (m *defaultAliyunLogProjectModel) formatProjectNameRegionKey(owner, region, name interface{}) string {
	return fmt.Sprintf("%s%v:%v:%v", cacheAliyunLogProjectNameRegionPrefix, owner, region, name)
}

func (m *defaultAliyunLogProjectModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", aliyunLogProjectRows, m.table)
	return conn.QueryRow(v, query, primary)
}
