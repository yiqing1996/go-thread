package thread

import (
	"context"
	"errors"
	"sync"
)

type threadInfo struct {
	running bool
	cancel  context.CancelFunc
}

type Pool struct {
	sync.RWMutex
	threads map[string]*threadInfo
}

func Test(a, b int) int {
	return a / b
}

func NewThreadPool() *Pool {
	return &Pool{threads: make(map[string]*threadInfo)}
}

func (p *Pool) Start(name string) (context.Context, error) {
	p.Lock()
	defer p.Unlock()
	if _, exists := p.threads[name]; exists {
		return nil, errors.New("线程名称已存在")
	}

	ctx, cancel := context.WithCancel(context.Background())
	p.threads[name] = &threadInfo{cancel: cancel, running: true}
	return ctx, nil
}

// Stop 通知指定线程停止
func (p *Pool) Stop(name string) error {
	p.Lock()
	defer p.Unlock()
	if t, exists := p.threads[name]; exists {
		t.cancel()
		return nil
	}
	return errors.New("指定名称线程不存在")
}

// StopAll 通知所有线程停止
func (p *Pool) StopAll() {
	p.Lock()
	defer p.Unlock()
	for _, t := range p.threads {
		t.cancel()
	}
}

// SetStop 由线程自己在退出时调用
func (p *Pool) SetStop(name string) error {
	p.Lock()
	defer p.Unlock()
	if _, exists := p.threads[name]; exists {
		delete(p.threads, name)
		return nil
	}
	return errors.New("线程不存在")
}

// IsRunning 线程是否运行中
func (p *Pool) IsRunning(name string) bool {
	p.RLock()
	defer p.RUnlock()
	if _, exists := p.threads[name]; exists {
		return p.threads[name].running
	}
	return false
}

func (p *Pool) GetAllRunning() []string {
	p.RLock()
	defer p.RUnlock()
	list := make([]string, 0)
	for name, t := range p.threads {
		if t.running {
			list = append(list, name)
		}
	}
	return list
}
