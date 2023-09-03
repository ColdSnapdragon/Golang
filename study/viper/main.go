package main

// https://jianghushinian.cn/2023/04/25/how-to-use-viper-for-configuration-management-in-go/

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"time"
)

func main() {
	v := viper.New() // 代表一个配置文件
	v1 := viper.New()

	// 为key设置默认值value (设置默认配置，优先级最低)
	v.SetDefault("username", "ice") // viper对于键不区分大小写
	v.SetDefault("server", map[string]string{"ip": "127.0.0.1", "port": "8080"})

	// 方式一
	v.SetConfigFile("./viper/config/test.yaml") // 指定配置文件(路径 + 配置文件名)
	v.SetConfigType("yaml")                     // 如果配置文件名中没有扩展名，则需要显式指定配置文件的格式
	// 读取配置文件到内存中解析为Viper配置对象
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { // 判断是哪类错误。viper提供了一系列错误类型用于诊断
			fmt.Println(errors.New("配置文件不存在"))
		} else {
			fmt.Println(errors.New("其他错误"), err)
		}
		os.Exit(1)
	}

	// 方式二
	v1.AddConfigPath(".")              // 添加配置路径。对一个配置文件会按指定的路径顺序查询，一旦找到便不再向下
	v1.AddConfigPath("./viper/config") // 相对路径是基于go.mod所在目录
	v1.SetConfigName("test1")          // 指定配置文件名(没有拓展名)
	_ = v1.ReadInConfig()
	// 监听配置文件变化
	v1.WatchConfig()
	// 注册每次配置文件发生变更后都会调用的回调函数
	v1.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("%s changed\n", in.Name)
	})
	time.Sleep(3 * time.Second) // 暂停3秒，可修改配置文件观察效果

	var yaml = []byte(`
message: hello
sender: ice
`)
	v2 := viper.New()
	v2.SetConfigType("yaml")                 // 设置配置文件类型
	_ = v2.ReadConfig(bytes.NewBuffer(yaml)) //从io.Reader读取配置(这里把string包装为buffer，后者提供read方法)

	// 打印最后一次ReadInConfig所使用的配置文件路径
	fmt.Println("读取到配置文件:", v.ConfigFileUsed()) // 返回用于填充配置注册表的配置文件路径
	fmt.Printf("%+v\n", v.Get("username"))      // 值不存在会返回零值
	fmt.Printf("%+v\n", v.Get("contributor"))
	fmt.Printf("%+v\n", v1.Get("uid"))
	fmt.Printf("%+v\n", v2.Get("message"))

	v3 := viper.New()
	_ = v3.BindEnv("GOPATH") // 读取(绑定)一个环境变量(区分大小写)
	fmt.Printf("%+v\n", v3.Get("GOPATH"))

	type Config struct {
		username string // 不会参与序列化和反序列化
		// Viper 支持嵌套结构体
		Server struct {
			IP   string // viper不区分键大小写(如这里的'P')
			Port string
		}
	}
	var cfg Config
	_ = v.Unmarshal(&cfg) // 反序列化
	fmt.Printf("%+v\n", cfg)
	var u string
	_ = v.UnmarshalKey("username", &u) // 反序列化指定键
	fmt.Printf("%+v\n", u)
	// 序列化另见文章

}
