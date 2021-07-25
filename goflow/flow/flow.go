package flow

import (
	"fmt"
	"sync"
)

type Flow struct {
	// 整个任务结构
	// 通过map来存储，具体task的key和taskItem
	funcs map[string]*flowStruct
}

// task func
// res 依赖上层func的计算结果，通过map存储，key为func name, value为计算结果
// 返回为当前func的计算结果
type flowFunc func(res map[string]interface{}) (interface{}, error)

// 单任务的结构
type flowStruct struct {
	Deps []string // 依赖哪些上层func的返回
	Ctr  int      // 下层有几个func依赖此func的返回
	Fn   flowFunc // 该层func
	C    chan interface{} // 该层func的返回结果，个数等于Ctr. 因为每个下游chan均需要一条返回结果
	once sync.Once // 用来唯一关闭C chan
}

func (fs *flowStruct) Done(r interface{}) {
	// 将task的执行结果，写入到结果C中，chan的条数等于下游依赖的个数
	for i := 0; i < fs.Ctr; i++ {
		fs.C <- r
	}
}

func (fs *flowStruct) Close() {
	// 关闭返回结果chan,once保证唯一次执行，防止deadlock
	fs.once.Do(func() {
		close(fs.C)
	})
}

func (fs *flowStruct) Init() {
	// 初始化返回结果chan的长度
	fs.C = make(chan interface{}, fs.Ctr)
}

// New flow struct
// 实例化一个DAG任务
func New() *Flow {
	return &Flow{
		funcs: make(map[string]*flowStruct),
	}
}

func (flw *Flow) Add(name string, d []string, fn flowFunc) *Flow {
	// 添加DAG节点任务
	// name: 任务名称，d: 上游依赖函数结果，fn: 该节点fn
	flw.funcs[name] = &flowStruct{
		Deps: d,
		Fn:   fn,
		Ctr:  1, // prevent deadlock，最少一个下游依赖，防止deadlock.(此时为返回结果)
	}
	return flw
}

func (flw *Flow) Do() (map[string]interface{}, error) {
	// DAG图实际执行
	for name, fn := range flw.funcs {
		// 遍历所有task进行校验及下游依赖个数的修正
		for _, dep := range fn.Deps {
			// prevent self depends
			// 避免自已依赖自己的环
			if dep == name {
				return nil, fmt.Errorf("Error: Function \"%s\" depends of it self!", name)
			}
			// prevent no existing dependencies
			// 避免依赖不存在func
			if _, exists := flw.funcs[dep]; exists == false {
				return nil, fmt.Errorf("Error: Function \"%s\" not exists!", dep)
			}
			// 修正某func的下游依赖个数，实际add节点时，无法知道有多少下游依赖自己
			flw.funcs[dep].Ctr++
		}
	}
	// 实际执行
	return flw.do()
}

func (flw *Flow) do() (map[string]interface{}, error) {
	// 实际执行
	var err error
	// 用来存最后的返回结果，key为name, value为该func的执行结果
	res := make(map[string]interface{}, len(flw.funcs))

	// 实例化 实际需要返回的C chan个数
	for _, f := range flw.funcs{
		f.Init()
	}
	// 每个节点的func都使用一个goroutine来执行
	for name, f := range flw.funcs {
		go func(name string, fs *flowStruct) {
			defer func() { fs.Close() }()
			// 用来存上游依赖的返回结果
			results := make(map[string]interface{}, len(fs.Deps))

			// drain dependency results
			// 获取上游依赖的返回结果，通过chan阻塞的方式，保证上游依赖都执行完成
			for _, dep := range fs.Deps {
				results[dep] = <-flw.funcs[dep].C
			}

			// 该节点的fun实际执行，传入依赖的上游数据
			r, fnErr := fs.Fn(results)
			if fnErr != nil {
				// close all channels
				// 如果该节点执行失败，则关闭所有chan，并直接退出，不用进行等待
				for _, fn := range flw.funcs {
					fn.Close()
				}
				err = fnErr
				return
			}
			// exit if error
			if err != nil {
				return
			}
			// 对该节点的执行结果进行下游依赖数据个数chan的赋值
			fs.Done(r)

		}(name, f)
	}

	// wait for all
	// 通过chan阻塞的方式，获取所有func的返回结果
	for name, fs := range flw.funcs {
		res[name] = <-fs.C
	}
	// 返回结果
	return res, err
}
