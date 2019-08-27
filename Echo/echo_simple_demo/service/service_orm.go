package service

import (
	"echodemo/conf"
	"github.com/donnie4w/go-logger/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strconv"
)

var (
	db *sqlx.DB
)

func Start()  {
	logger.Info("start connection databases...")
	logger.Debug(conf.Config.Database.Db_user+":"+conf.Config.Database.Db_pass+"@tcp("+conf.Config.Database.Db_host+":"+strconv.Itoa(conf.Config.Database.Db_port)+")/"+conf.Config.Database.Db_name+"?parseTime=true")
	sqlSession, err := sqlx.Connect("mysql", conf.Config.Database.Db_user+":"+conf.Config.Database.Db_pass+"@tcp("+conf.Config.Database.Db_host+":"+strconv.Itoa(conf.Config.Database.Db_port)+")/"+conf.Config.Database.Db_name+"?parseTime=true")
	if err != nil{
		logger.Error(err)
	}
	db = sqlSession
	logger.Info("connection database success...")

	// 测试数据库查询
	var count int
	if err = db.Get(&count, "SELECT count(*) FROM t_user"); err != nil {
		logger.Error(err)
		return
	}
	logger.Debug("count t_user :", count)

}