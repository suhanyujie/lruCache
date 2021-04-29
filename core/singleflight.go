package core

import "sync"

type call struct {
	wg sync.WaitGroup
	val interface{}
	err error
}

type SFGroup struct {
	mu sync.Mutex
	m map[string]*call
}

func (sfg *SFGroup) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	sfg.mu.Lock()
	if sfg.m == nil {
		sfg.m = make(map[string]*call, 0)
	}
	if c, ok := sfg.m[key]; ok {
		// 如果重复释放会怎样呢？
		sfg.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	sfg.m[key] = c
	sfg.mu.Unlock()
	c.val, c.err = fn()
	c.wg.Done()

	sfg.mu.Lock()
	delete(sfg.m, key)
	sfg.mu.Unlock()

	return c.val, c.err
}
