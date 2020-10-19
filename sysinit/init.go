package sysinit

// 进行初始化的函数，进行一系列初始化操作，包括初始化系统，初始化主数据库，初始化从数据库
func init(){
	sysinit()
	dbinit()
	//dbinit("r")
	//dbinit("uaw","uar")
}
