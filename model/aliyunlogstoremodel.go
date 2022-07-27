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
	aliyunLogstoreFieldNames           = builder.RawFieldNames(&AliyunLogstore{})
	aliyunLogstoreRows                 = strings.Join(aliyunLogstoreFieldNames, ",")
	aliyunLogstoreRowsExpectAutoSet    = strings.Join(stringx.Remove(aliyunLogstoreFieldNames, "`id`"), ",")
	aliyunLogstoreRowsWithPlaceHolder  = strings.Join(stringx.Remove(aliyunLogstoreFieldNames, "`id`", "`create_time`"), "=?,") + "=?"
	aliyunLogstoreRowsWithPlaceHolder2 = strings.Join(stringx.Remove(aliyunLogstoreFieldNames, "`id`", "`name`", "`project_name`", "`endpoint`", "`create_time`", "`owner`"), "=?,") + "=?"

	cacheAliyunLogstoreIdPrefix          = "cache:aliyunLogstore:id:"
	cacheAliyunLogstoreNameProjectPrefix = "cache:aliyunLogstore:owner:endpoint:project:name:"
)

type (
	AliyunLogstoreModel interface {
		Insert(session sqlx.Session, data *AliyunLogstore) (sql.Result, error)
		FindOne(id int64) (*AliyunLogstore, error)
		FindAll() ([]*AliyunLogstore, error)
		FindOneByOwnerEndpointProjectName(owner, endpoint, projectName string, name string) (*AliyunLogstore, error)
		Update(session sqlx.Session, data *AliyunLogstore) error
		InsertOrUpdate(session sqlx.Session, data *AliyunLogstore) error
		Delete(session sqlx.Session, id int64) error
		Trans(fn func(session sqlx.Session) error) error
	}

	defaultAliyunLogstoreModel struct {
		sqlc.CachedConn
		table string
	}

	AliyunLogstore struct {
		Id             int64     `db:"id"`
		Name           string    `db:"name"`             // 日志库名称
		ProjectName    string    `db:"project_name"`     // 日志项目名称
		Endpoint       string    `db:"endpoint"`         // 日志库所在地域
		Ttl            int64     `db:"ttl"`              // 数据的保存时间
		ShardCount     int64     `db:"shardCount"`       // shard分区数
		EnableTracking int64     `db:"enable_tracking"`  // 是否开启Webtracking功能: 1:开启 0：关闭
		AutoSplit      int64     `db:"auto_split"`       // 是否自动分裂shard: 1:自动分裂 0：不自动分裂
		MaxSplitShard  int64     `db:"max_split_shard"`  // 自动分裂时最大的shard个数，最小值为1，最大值为64
		AppendMeta     int64     `db:"appendMeta"`       // 是否记录外网IP地址的功能：1：记录 0：不记录
		TelemetryType  int64     `db:"telemetry_type"`   // 要查询的日志类型：1：Metrics(时序数据）0：None：非时序存储
		CreateTime     time.Time `db:"create_time"`      // 日志库创建时间
		LastModifyTime time.Time `db:"last_modify_time"` // 日志库更新时间
		Owner          string    `db:"owner"`            // 日志库拥有者
		Maintainer     string    `db:"maintainer"`       // 日志库维护者
	}
)

func NewAliyunLogstoreModel(conn sqlx.SqlConn, c cache.CacheConf) AliyunLogstoreModel {
	return &defaultAliyunLogstoreModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`aliyun_logstore`",
	}
}

func (m *defaultAliyunLogstoreModel) Insert(session sqlx.Session, data *AliyunLogstore) (result sql.Result, err error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, aliyunLogstoreRowsExpectAutoSet)

	if session != nil {
		result, err = session.Exec(query, data.Name, data.ProjectName, data.Endpoint, data.Ttl, data.ShardCount, data.EnableTracking, data.AutoSplit, data.MaxSplitShard, data.AppendMeta, data.TelemetryType, data.CreateTime, data.LastModifyTime, data.Owner, data.Maintainer)
	} else {
		result, err = m.ExecNoCache(query, data.Name, data.ProjectName, data.Endpoint, data.Ttl, data.ShardCount, data.EnableTracking, data.AutoSplit, data.MaxSplitShard, data.AppendMeta, data.TelemetryType, data.CreateTime, data.LastModifyTime, data.Owner, data.Maintainer)
	}
	if err != nil {
		return
	}
	lastInsertId, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()
	if lastInsertId != 0 && rowsAffected != 0 {
		err = m.DelCache(m.formatPrimary(lastInsertId), m.formatNameProjectKey(data.Owner, data.Endpoint, data.ProjectName, data.Name))
		if err != nil {
			logx.Errorf("删除缓存失败，错误信息：%v", err)
		}
	}
	return
}

func (m *defaultAliyunLogstoreModel) FindOne(id int64) (*AliyunLogstore, error) {
	var resp AliyunLogstore
	err := m.QueryRow(&resp, m.formatPrimary(id), func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", aliyunLogstoreRows, m.table)
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

func (m *defaultAliyunLogstoreModel) FindAll() ([]*AliyunLogstore, error) {
	var resp []*AliyunLogstore
	query := fmt.Sprintf("select %s from %s", aliyunLogstoreRows, m.table)
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

func (m *defaultAliyunLogstoreModel) FindOneByOwnerEndpointProjectName(owner, endpoint, projectName string, name string) (*AliyunLogstore, error) {
	aliyunLogStoreNameProjectKey := fmt.Sprintf("%s%v:%v:%v:%v", cacheAliyunLogstoreNameProjectPrefix, owner, endpoint, projectName, name)
	var resp AliyunLogstore
	err := m.QueryRowIndex(&resp, aliyunLogStoreNameProjectKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `owner` = ? and `endpoint` = ? and `project_name` = ? and `name` = ? limit 1", aliyunLogstoreRows, m.table)
		if err := conn.QueryRow(&resp, query, owner, endpoint, projectName, name); err != nil {
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

func (m *defaultAliyunLogstoreModel) Update(session sqlx.Session, data *AliyunLogstore) error {
	var (
		result sql.Result
		err    error
	)
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, aliyunLogstoreRowsWithPlaceHolder)

	if session != nil {
		result, err = session.Exec(query, data.Name, data.ProjectName, data.Endpoint, data.Ttl, data.ShardCount, data.EnableTracking, data.AutoSplit, data.MaxSplitShard, data.AppendMeta, data.TelemetryType, data.LastModifyTime, data.Owner, data.Maintainer, data.Id)
	} else {
		result, err = m.ExecNoCache(query, data.Name, data.ProjectName, data.Endpoint, data.Ttl, data.ShardCount, data.EnableTracking, data.AutoSplit, data.MaxSplitShard, data.AppendMeta, data.TelemetryType, data.LastModifyTime, data.Owner, data.Maintainer, data.Id)
	}

	if err != nil {
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()
	if lastInsertId != 0 && rowsAffected != 0 {
		err = m.DelCache(m.formatPrimary(lastInsertId), m.formatNameProjectKey(data.Owner, data.Endpoint, data.ProjectName, data.Name))
		if err != nil {
			logx.Errorf("删除缓存失败，错误信息：%v", err)
		}
	}
	return err
}

func (m *defaultAliyunLogstoreModel) InsertOrUpdate(session sqlx.Session, data *AliyunLogstore) error {
	var (
		result sql.Result
		err    error
	)
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) on DUPLICATE KEY UPDATE %s", m.table, aliyunLogstoreRowsExpectAutoSet, aliyunLogstoreRowsWithPlaceHolder2)

	if session != nil {
		result, err = session.Exec(query, data.Name, data.ProjectName, data.Endpoint, data.Ttl, data.ShardCount, data.EnableTracking, data.AutoSplit, data.MaxSplitShard, data.AppendMeta, data.TelemetryType, data.CreateTime, data.LastModifyTime, data.Owner, data.Maintainer,
			data.Ttl, data.ShardCount, data.EnableTracking, data.AutoSplit, data.MaxSplitShard, data.AppendMeta, data.TelemetryType, data.LastModifyTime, data.Maintainer)
	} else {
		result, err = m.ExecNoCache(query, data.Name, data.ProjectName, data.Endpoint, data.Ttl, data.ShardCount, data.EnableTracking, data.AutoSplit, data.MaxSplitShard, data.AppendMeta, data.TelemetryType, data.CreateTime, data.LastModifyTime, data.Owner, data.Maintainer,
			data.Ttl, data.ShardCount, data.EnableTracking, data.AutoSplit, data.MaxSplitShard, data.AppendMeta, data.TelemetryType, data.LastModifyTime, data.Maintainer)
	}

	if err != nil {
		return err
	}
	lastInsertId, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()
	if lastInsertId != 0 && rowsAffected != 0 {
		err = m.DelCache(m.formatPrimary(lastInsertId), m.formatNameProjectKey(data.Owner, data.Endpoint, data.ProjectName, data.Name))
		if err != nil {
			logx.Errorf("删除缓存失败，错误信息：%v", err)
		}
	}
	return err
}

func (m *defaultAliyunLogstoreModel) Delete(session sqlx.Session, id int64) error {

	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		if session != nil {
			return session.Exec(query, id)
		}
		return conn.Exec(query, id)
	}, m.formatPrimary(id))
	return err
}

func (m *defaultAliyunLogstoreModel) Trans(fn func(session sqlx.Session) error) error {
	err := m.Transact(func(session sqlx.Session) error {
		return fn(session)
	})
	return err
}

func (m *defaultAliyunLogstoreModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheAliyunLogstoreIdPrefix, primary)
}

func (m *defaultAliyunLogstoreModel) formatNameProjectKey(owner, endpoint, project, name interface{}) string {
	return fmt.Sprintf("%s%v:%v:%v:%v", cacheAliyunLogstoreNameProjectPrefix, owner, endpoint, project, name)
}

func (m *defaultAliyunLogstoreModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", aliyunLogstoreRows, m.table)
	return conn.QueryRow(v, query, primary)
}
