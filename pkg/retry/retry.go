package retry

import (
	"fmt"
	"time"
)

/*
Retry 对执行失败的函数，在进行几次重试
*/

// Try 重新尝试函数，如果函数执行失败，则延长时间重试
type Try struct {
	Name     string        // 重试任务名称
	F        func() error  // 需要重新尝试的函数
	Duration time.Duration // 重新尝试的时间间隔
	MaxTimes int           // 最大重试次数
}

func NewTry(name string, f func() error, duration time.Duration, maxTimes int) *Try {
	return &Try{
		Name:     name,
		F:        f,
		Duration: duration,
		MaxTimes: maxTimes,
	}
}

// Report 尝试重试的报告
type Report struct {
	Name        string        // 重试任务名称
	Result      bool          // 函数执行的结果
	Times       int           // 重试的次数
	SumDuration time.Duration // 总执行时间
	Errs        []error       // 函数执行的错误记录
}

func (r *Report) Error() string {
	return fmt.Sprintf("[retry]名称：%s，结果：%v，尝试次数：%v，总时间：%v，错误：%v", r.Name, r.Result, r.Times, r.SumDuration, r.Errs)
}

// Run 尝试重试，返回 chan 可以用于接收尝试报告
func (try *Try) Run() <-chan Report {
	result := make(chan Report, 1)
	go func() {
		defer close(result)
		start := time.Now()
		var errs []error
		for i := 0; i < try.MaxTimes; i++ {
			time.Sleep(try.Duration)
			err := try.F()
			if err == nil {
				result <- Report{
					Name:        try.Name,
					Result:      true,
					Times:       i + 1,
					SumDuration: time.Since(start),
					Errs:        errs,
				}
				return
			}
			errs = append(errs, err)
		}
		result <- Report{
			Name:        try.Name,
			Result:      false,
			Times:       try.MaxTimes,
			SumDuration: time.Since(start),
			Errs:        errs,
		}
	}()
	return result
}
