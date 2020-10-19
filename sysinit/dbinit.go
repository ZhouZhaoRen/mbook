package sysinit

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // 必须得引入这个mysql驱动
)

// 由于需要初始化不同的数据库，所以传进来的参数不同，参数个数也不同
func dbinit(alias ...string){
	// 先看是否为开发模式，如果是，则显示命令信息
	fmt.Println("开始初始化数据库")
	isDev:=beego.AppConfig.String("runmode")=="dev"
	if isDev{
		orm.Debug=isDev  // 显示开发命令信息
	}

	if len(alias)>0{
		for _,value:=range alias{
			registerDatabase(value)
			// 主库自动建表
			if "w"==value{
				orm.RunSyncdb("default",false,isDev) //
			}
		}
	}else{
		registerDatabase("w")
		orm.RunSyncdb("default",false,isDev)
	}

}

// 初始化数据库，对数据库的连接参数进行拼接
func registerDatabase(alias string){
	fmt.Println("开始注册数据库")
	if len(alias)==0{
		return
	}

	// 连接名称
	dbAlias:=alias
	if alias=="w" || alias=="default"{
		dbAlias="default"
		alias="w"
	}
	fmt.Printf("数据库别名为%s：",dbAlias)
	// 数据库名称
	dbName:=beego.AppConfig.String("db_"+alias+"_database")
	// 数据库连接用户名
	dbUser := beego.AppConfig.String("db_" + alias + "_username")
	// 数据库连接密码
	dbPwd := beego.AppConfig.String("db_" + alias + "_password")
	// 数据库IP
	dbHost := beego.AppConfig.String("db_" + alias + "_host")
	// 数据库端口号
	dbPort := beego.AppConfig.String("db_" + alias + "_port")

	// 进行数据库的注册
	orm.RegisterDataBase(dbAlias,"mysql",dbUser+":"+dbPwd+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8", 30)

}
