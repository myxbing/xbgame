package timer

import (
	"sync"
	"sync/atomic"
	"time"
)

//定时器结构
type _Timer struct {
	Id int32
	Timeout int64
	CH chan int32
}

const (
	TIMER_INTERVAL_LEVEL = uint(16) // 时间间隔最大分级
)

var (
	_events [TIMER_INTERVAL_LEVEL]map[uint32]*_Timer // 事件列表
	_queue      map[uint32]*_Timer // 事件添加队列
	_lock sync.Mutex
	_identity uint32 // 内部事件编号
)

func init() {
	for k := range _events {
		_events[k] = make(map[uint32]*_Timer)
	}
	_queue = make(map[uint32]*_Timer)
	go _timer()
}

//------------------------------------------------
// 定时器 goroutine
// 根据程序启动后经过的秒数计数
func _timer() {
	timer_count := uint32(0)
	last := time.Now().Unix()

	for {
		time.Sleep(100 * time.Millisecond)

		// 处理排队
		// 最小的时间间隔，处理为1s
		_lock.Lock()
		for k, v := range _queue {
			diff := v.Timeout - time.Now().Unix()
			if diff <= 0 {
				diff = 1
			}

			for i := TIMER_INTERVAL_LEVEL - 1; i >= 0; i-- {
				if diff >= 1<<i {
					_events[i][k] = v
					break
				}
			}
		}
		_queue = make(map[uint32]*_Timer)
		_lock.Unlock()

		// 检查事件触发
		// 累计距离上一次触发的秒数,并逐秒触发
		// 如果校正了系统时间，时间前移，nsec为负数的时候，last的值不应该变动，否则会出现秒数的重复计数
		now := time.Now().Unix()
		nsec := now - last

		if nsec <= 0 {
			continue
		} else {
			last = now
		}

		for c := int64(0); c < nsec; c++ {
			timer_count++

			for i := TIMER_INTERVAL_LEVEL - 1; i > 0; i-- {
				mask := (uint32(1) << i) - 1
				if timer_count&mask == 0 {
					_trigger(i)
				}
			}

			_trigger(0)
		}
	}
}

//触发定时器
func _trigger(level uint) {
	now := time.Now().Unix()
	list := _events[level]

	for k, v := range list {
		if v.Timeout-now < 1<<level {
			// 移动到前一个更短间距的LIST
			if level == 0 {
				func() {
					defer func() {
						recover() // ignore closed channel
					}()

					v.CH <- v.Id
				}()
			} else {
				_events[level-1][k] = v
			}

			delete(list, k)
		}
	}
}

// 添加一个定时，timeout为到期的Unix时间
// id 是调用者定义的编号, 事件发生时，会把id发送到ch
func SetTimer(id int32, timeout int64, ch chan int32) {
	event := &_Timer{Id: id, CH: ch, Timeout: timeout}
	timer_id := atomic.AddUint32(&_identity, 1)
	_lock.Lock()
	_queue[timer_id] = event
	_lock.Unlock()
}
