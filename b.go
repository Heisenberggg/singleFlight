
package singleflight // import "golang.org/x/sync/singleflight"

import "sync"

type call struct {
	wg sync.WaitGroup  // 用于阻塞这个调用call的其他请求
	val interface{}	   // 函数执行后的结果
	err error		   // 函数执行后的error
}

type Group struct {
	mu sync.Mutex       // protects m
	m  map[string]*call // 对于每一个需要获取的key有一个对应的call
}

type Result struct {
	Val    interface{}
	Err    error
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error) {
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	//如果获取当前key的函数正在被执行，则获取它的执行结果
	if c, ok := g.m[key]; ok {
		return c.val, c.err, true
	}
	c := new(call)
	g.m[key] = c

	g.doCall(c, key, fn)
	return c.val, c.err
}

func (g *Group) doCall(c *call, key string, fn func() (interface{}, error)) {
	//执行传入的方法
	c.val, c.err = fn()
}
