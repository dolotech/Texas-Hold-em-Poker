// 数据库客户端常用操作 （基于 xorm）
package db

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

// Client 数据库客户端
type Client struct {
	engine *xorm.Engine
}

var db *Client

//var database_addr string
//const configFilePath string = "config.ini"

func Init(url string) error{
	var err error
	if db == nil {
		db, err = NewClient("postgres", url)
		//db.ShowSQL(config.GetCfgDatabase().ShowSql)
		db.ShowSQL(true)
	}
	return err
}
func C() *Client {
	return db
}

// NewClient 创建一个客户端链接
func NewClient(driver, connstr string) (*Client, error) {
	engine, err := xorm.NewEngine(driver, connstr)
	engine.SetMaxIdleConns(200) //最大连接数设置为200
	if err == nil {
		ping := func() <-chan string {
			errmsg := make(chan string)
			go func() {
				err := engine.Ping()
				if err != nil {
					errmsg <- err.Error()
				} else {
					errmsg <- ""
				}
			}()
			return errmsg
		}
		go func() {
			select {
			case msg := <-ping():
				if len(msg) > 0 {
					//logs.Error(msg)
				} else {
					//logs.Info("connect %s succeed.", driver)
				}
			case <-time.After(time.Second * 5):
				//logs.Info("connect timeout 5 second.")
			}
		}()
	}
	return &Client{engine: engine}, err
}

func (cl *Client) Engine() *xorm.Engine {
	return cl.engine
}

// Close 关闭连接
func (cl *Client) Close() error {
	if cl.engine != nil {
		err := cl.engine.Close()
		cl.engine = nil
		return err
	}
	return nil
}

// SerialInteractor 序列接口
type SerialInteractor interface {
	TableName() string
}

// ShowSQL 显示 SQL
func (cl *Client) ShowSQL(show bool) {
	cl.engine.ShowSQL(show)
}

// NewSession 创建一个事务
func (cl *Client) NewSession() *xorm.Session {
	return cl.engine.NewSession()
}

// GetNextSerial 获取下个以序列
func (cl *Client) GetNextSerial(serialbean SerialInteractor) (int, error) {
	sqlStr := fmt.Sprintf("select nextval('%s_id_seq'::regclass)", serialbean.TableName())
	data, err := cl.engine.Query(sqlStr)
	if err != nil || data == nil || len(data) == 0 {
		return 0, err
	}
	nextval := string(data[0]["nextval"])
	return strconv.Atoi(nextval)
}

// SimpleValue 简单结果
func (cl *Client) SimpleValue(sql string) (string, error) {
	data, err := cl.engine.Query(sql)
	if err != nil || data == nil || len(data) == 0 {
		return "", err
	}
	for _, v := range data[0] {
		ret := string(v)
		return ret, nil
	}
	return "", errors.New("not found return value")
}

// SimpleIntValue 简单的 int 值
func (cl *Client) SimpleIntValue(sql string) (int, error) {
	data, err := cl.engine.Query(sql)
	if err != nil || data == nil || len(data) == 0 {
		return 0, err
	}
	for _, v := range data[0] {
		ret := string(v)
		return strconv.Atoi(ret)
	}
	return 0, errors.New("not found return value")
}

// Insert 保存一个数据
func (cl *Client) Insert(bean interface{}, omitColumns ...string) (int64, error) {
	return cl.engine.Omit(omitColumns...).Insert(bean)
}

// Get 查找一条数据
func (cl *Client) Get(bean interface{}) (bool, error) {
	has, err := cl.engine.Get(bean)
	return has, err
}

// WhereGet 条件查询
func (cl *Client) WhereGet(bean interface{}, querystring string, args ...interface{}) (bool, error) {
	return cl.engine.Where(querystring, args...).Get(bean)
}

// WhereFind 条件查询
func (cl *Client) WhereFind(bean interface{}, querystring string, args ...interface{}) error {
	return cl.engine.Where(querystring, args...).Find(bean)
}

// Update 修改一条数据
func (cl *Client) Update(id int, bean interface{}, omitClumns ...string) error {
	session := cl.engine.NewSession()
	defer session.Close()
	rows, err := session.Omit(omitClumns...).Id(id).Update(bean)
	if rows > 1 {
		_ = session.Rollback()
		return fmt.Errorf("修改失败，影响行数: %d ", rows)
	}
	_ = session.Commit()
	return err
}

// Delete 删除数据
func (cl *Client) Delete(bean interface{}) error {
	session := cl.engine.NewSession()
	rows, err := session.Delete(bean)
	if rows > 1 {
		_ = session.Rollback()
		return fmt.Errorf("修改失败，影响行数: %d ", rows)
	}
	_ = session.Commit()
	return err
}

// // Query 多数据查询
// func (cl *Client) Query(beans interface{}, querystring string, args ...interface{}) error {
// 	return cl.engine.Where(querystring, args...).Find(beans)
// }

// Where 条件查询
func (cl *Client) Where(beans interface{}, querystring string, args ...interface{}) error {
	return cl.engine.Where(querystring, args...).Find(beans)
}

// In 条件查询 IN list
func (cl *Client) In(beans interface{}, querystring string, args ...interface{}) error {
	return cl.engine.In(querystring, args...).Find(beans)
}

// SQL 查询
func (cl *Client) SQL(beans interface{}, querystring string, args ...interface{}) error {
	return cl.engine.Sql(querystring, args...).Find(beans)
}

// QueryValue 查询
func (cl *Client) QueryValue(sqlstr string) ([]map[string][]byte, error) {
	data, err := cl.engine.Query(sqlstr)
	return data, err
}

// Create 新建一个表
func (cl *Client) CreateTables(beans ...interface{}) error {
	return cl.engine.CreateTables(beans...)
}

// Sync 同步新表
func (cl *Client) Sync(beans ...interface{}) error {
	return cl.engine.Sync(beans...)
}
