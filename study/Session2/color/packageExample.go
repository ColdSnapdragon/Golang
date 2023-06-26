package color

type Color string

const (
	Green = "green" // 以大写字母开头，让其他包能访问成员
	Blue  = "blue"
	red   = "red" // 不能访问
)