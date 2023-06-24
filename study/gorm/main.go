package main

// https://gorm.io/zh_CN/docs/
// https://learnku.com/docs/gorm/v2
import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type UserInfo struct {
	gorm.Model // GORM定义一个gorm.Model结构体，含有若干字段(ID字段为主键)
	// AUTO_INCREMENT已被设置在ID字段了，后续再设置是无效的
	Name   string `gorm:"type:varchar(20);default:ice"`
	Age    uint   `gorm:"default:18;not null"`
	Number *uint  `gorm:""`
	// Number的类型为指针，这样当没有指定Number时才会去用默认值，若无默认值用NULL(没有指定not null)
	// 如果不是指针类型，那么未指定时采用默认值(有设置的话)或类型零值，永远不会采用到NULL
}

// TableName gorm默认有一套命名机制。可以通过设置TableName方法来指定表名
func (UserInfo) TableName() string {
	return "users"
}

func main() {
	dsn := "root:380024507@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// db 是GORM创建的默认大小为10的连接池，并自动管理回收
	if err != nil {
		log.Fatal(err)
	}

	//db.Raw("DROP TABLE users;") // 执行一条原生SQL (危险操作)

	// 创建表自动迁移（把结构体和数据表进行对应）
	// 若表以存在会根据模型对象的定义检查并更新表结构
	if err = db.AutoMigrate(&UserInfo{}); err != nil {
		log.Fatal(err)
	}

	db.Unscoped().Where("1 = 1").Delete(&UserInfo{}) // 清空记录(Unscoped方法表示采用硬删除)

	if true {
		u1 := UserInfo{Name: "", Age: 20}                 // 不指定Name或者传零值，效果都是采用默认值("ice")。Number为NULL
		db.Create(&u1)                                    // 插入值。数据较大时可以通传入数据的指针
		u2 := UserInfo{Name: "blover", Number: new(uint)} // Number为零值。通过指针类型可以避免因为是零值而去找默认值
		db.Create(&u2)
	}

	if true {
		var u UserInfo
		db.Where("id < ?", "3").Where("NOT ? IS NULL", "number").First(&u) // Frist取主键排序后第一条，Last取最后一条
		// SELECT * FROM `users` WHERE id < 3 AND NOT 'number' IS NULL AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LI
		//MIT 1;
		// 安全的，会被转义，可避免SQL注入
		// Where与Where之间是And关系，上面这些方法都返回可继续操作的指针(链式操作)
		fmt.Printf("u: %v\n", u.Name)
		// Debug()可以打印最后生成的SQL语句
		db.Debug().Where("id IN (?)", []uint{1, 5, 9}).Or("name LIKE ?", "%e").Take(&u) // Take是随机取一条
		// SELECT * FROM `users` WHERE (id IN (1,5,9) OR name LIKE '%e') AND `users`.`deleted_at` IS NULL AND `users`.`id` = 1 LIMIT 1
		fmt.Printf("u: %v\n", u.Name)
		// 其他演示
		var us []UserInfo
		db.Where(&UserInfo{Name: "", Age: 18}).Find(&us) // age = 18 。注意，当使用结构作为条件查询时，GORM只会查询非零值字段(这里Name是无效的)
		db.Where([]int64{20, 21, 22}).Find(&us)          // id IN (20, 21, 22);
		tx := db.Table("users").Select("name")           // 指定所查表、指定所查字段。返回一个指针，包装了已有的查询
		tx.First(&u)                                     // 当遇到像Take First这样的*立即执行方法*时，才会真正地发起查询
		tx.Not("id < ?", "1000").Take(&u)                // 继续叠加过滤条件。这里由于查询为空(id都小于1000)，所以u的值没有被覆盖
		fmt.Printf("u: %v\n", u.Name)
		// 更多方法见文档。gorm
	}

	if true {
		var u UserInfo
		db.First(&u)
		db.Debug().Model(&u).Where("name = ?", "ice").Update("name", "Ice") // u的主键值(有的话)会作为过滤条件
		// UPDATE `users` SET `name`='Ice',`updated_at`='2023-06-23 21:09:53.759' WHERE id < '10' AND `users`.`deleted_at` IS NULL AND `id` = 1
		db.Debug().Where("1 = 1").Model(&UserInfo{}).Update("age", gorm.Expr("age + ?", 2)) // 所有age加2 (并且更新updated_at)
		// db.Where("id < 10").Model(&UserInfo{}).UpdateColumn("age", gorm.Expr("age + ?", 2)) // 同上，但只更新age字段
		// 过滤条件是必要的，gorm阻止了全局更新和全局删除
	}

	if false {
		var u UserInfo
		u.ID = 1
		// 由于有gorm.DeletedAt字段，会作软删除
		db.Debug().Delete(&u) // u的主键值(有的话)会作为过滤条件
		// UPDATE `users` SET `deleted_at`='2023-06-23 21:04:54.542' WHERE id < '10' AND `users`.`id` = 1 AND `users`.`deleted_at` IS NULL
		// db.Delete(&User{}, "1")
		// db.Delete(&User{}, []int{1})
	}

}
