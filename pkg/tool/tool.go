package tool

// DoThat 错误的函数式编程，优化 if err != nil 代码，外部传入 err 进行判断执行
func DoThat(err error, f func() error) error {
	if err != nil {
		return err
	}
	return f()
}
