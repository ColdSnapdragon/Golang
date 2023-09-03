package models

import "bubble/dao"

// models含有程序中的核心数据类型

type Todo struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Status *bool  `json:"status" gorm:"default:false"` // 前端传false时相当于零值，不用指针会影响更新操作
}

// 以下是一些与dao交互的逻辑操作

func AddOneTodo(item *Todo) error {
	return dao.DB.Create(&item).Error
}

func GetAllTodo() (list []Todo, err error) {
	err = dao.DB.Find(&list).Error
	return
}

func GetOneTodo(id string) (p *Todo, err error) {
	p = new(Todo)
	err = dao.DB.First(p, id).Error
	return
}

func UpdateOneTodo(id string, item *Todo) error {
	return dao.DB.Where("id = ?", id).Updates(item).Error
}

func DelOneTodo(id string) error {
	return dao.DB.Delete(&Todo{}, id).Error
}
