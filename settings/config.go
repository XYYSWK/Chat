package settings

import (
	"Chat/global"
	"Chat/pkg/tool"
	"flag"
	"github.com/XYYSWK/Lutils/pkg/setting"
	"strings"
)

//配置文件绑定到全局结构体上（默认加载）

var (
	configPaths       string //配置文件路径
	privateConfigName string //private 配置文件名
	publicConfigName  string //public 配置文件名
	configType        string //配置文件类型
)

// 通过 flag 包设置命令行参数
func setupFlag() {
	//命令行参数绑定
	flag.StringVar(&configPaths, "config_path", global.RootDir+"/config/app", "指定要使用的配置文件的路径，多个路径用逗号隔开")
	flag.StringVar(&privateConfigName, "private_config_name", "private", "private 配置文件名")
	flag.StringVar(&publicConfigName, "public_config_name", "public", "public 配置文件名")
	flag.StringVar(&configType, "config_type", "yaml", "配置文件类型")
	flag.Parse() //解析命令行参数，并将它们对应的值赋给相应的变量
}

type config struct {
}

// Init  读取配置文件，将配置文件上的内容映射到结构体中
func (config) Init() {
	setupFlag()
	var (
		err            error
		publicSetting  *setting.Setting
		privateSetting *setting.Setting
	)
	// 这是第一个 Init 函数，会在其他组件的 Init 函数执行之前执行，并将配置文件绑定到全局变量上
	err = tool.DoThat(err, func() error {
		//初始化 publicSetting 的基础属性
		publicSetting, err = setting.NewSetting(publicConfigName, configType, strings.Split(configPaths, ",")...) //引入配置文件路径
		return tool.DoThat(err, func() error { return publicSetting.BindAll(&global.PublicSetting) })             //将配置文件中的信息解析到全局变量中
	})
	err = tool.DoThat(err, func() error {
		//初始化 privateSetting 的基础属性
		privateSetting, err = setting.NewSetting(privateConfigName, configType, strings.Split(configPaths, ",")...) //引入配置文件路径
		return tool.DoThat(err, func() error { return privateSetting.BindAll(&global.PrivateSetting) })             //将配置文件中的信息解析到全局变量中
	})
	if err != nil {
		panic("读取配置文件有误：" + err.Error())
	}
}
