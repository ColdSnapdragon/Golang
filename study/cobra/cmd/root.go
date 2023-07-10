package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

// main.go作为程序的入口，而所有命令都应该放在 cmd/ 目录下

var (
	who     string
	print_b bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hello",
	Short: "the brief description of hello",    // 简短描述(没有详细描述时会替代)
	Long:  "the detailed description of hello", // 详细描述(--help）
	Run: func(cmd *cobra.Command, args []string) {
		if print_b {
			fmt.Println("hello!!!")
		}
		fmt.Printf("To %v\n", who)
	},
}

// 子命令
var versionCmd = &cobra.Command{
	Use: "version",
	// Args: cobra.NoArgs, // 参数验证。下面是自定义的类似做法
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("过多的参数")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:\t1.0.0")
	},
	// 了解四个Hooks函数及其Error版本
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Go Version:\tgo1.20.4")
	},
	PostRunE: func(cmd *cobra.Command, args []string) error { // 存在<Hooks>E版本时，会覆盖<Hooks>版本
		fmt.Println("Go Version:\tgo1.20.4")
		return nil
	},
}

func Execute() {
	log.Fatal(rootCmd.Execute()) // 命令的执行入口
}

// init函数会自动在初始时执行
func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cog.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	rootCmd.AddCommand(versionCmd) // 添加子命令

	/*
		Flags()和PersistentFlags()返回*pflag.FlagSet
		pflag完全兼容flag，可以用flag的方法
	*/

	// 本地标志。注册在什么命令上就只对什么命令可用
	flagSet := rootCmd.Flags()
	// print_b = flagSet.BoolP("print", "p", false, "print another message") // 'P'表示接受短标志
	// flagSet.Lookup("print").NoOptDefVal = "true"                           // 指定标志但没有传值时设定的值(不过对于bool类型，设与不设一样)
	// 上面两行没办法-p指定true，尚不明原因(怀疑是cobra参与后的问题)
	flagSet.BoolVarP(&print_b, "print", "p", false, "print another message")

	// 持续标志。对注册命令的子命令也是可用的
	rootCmd.PersistentFlags().StringVarP(&who, "who", "w", "", "the receiver") // 'Var'表示以传入参数传值
	rootCmd.MarkFlagRequired("who")                                            // 指定who标志必选

	// 解析命令行参数(对于没有'Var'的方法必须执行)
	// pflag.Parse() // 一旦用了，却什么参数也用不了
}
