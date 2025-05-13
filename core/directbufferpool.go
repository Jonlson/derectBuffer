package core

import (
	"sort"
	"sync"
	"sync/atomic"
)

const (
	minBitSize              = 6  // 2^6 = 64
	maxBitSize              = 32 // 2^32 = 4G
	steps                   = 4
	minSize                 = 1 << minBitSize
	maxSIze                 = 1 << (minBitSize + steps - 1)
	calibrateCallsThreshold = 42000 // 42000 calls per second 就是
	maxPercentile           = 0.95  // 95%的分位点

)

type Pool struct {
	calls       [steps]uint64
	calibrating uint64

	defaultSize uint64
	maxSize     uint64

	pool sync.Pool
}

var defaultPool Pool

func Get() *DirectBuffer { return defaultPool.Get() }

func (p *Pool) Get() *DirectBuffer {
	// 从池中获取一个对象
	v := p.pool.Get()
	// 如果获取的对象不为空
	if v != nil {
		// 将获取的对象转换为 *DirectBuffer 类型并返回
		return v.(*DirectBuffer)
	}
	// 如果获取的对象为空，则创建一个新的 *DirectBuffer 对象并返回
	return &DirectBuffer{
		// 初始化 data 字段，其长度为 0，容量为 p.defaultSize 的值
		data: make([]byte, 0, atomic.LoadUint64(&p.defaultSize)),
	}
}

func Put(b *DirectBuffer) { defaultPool.Put(b) }

// Put 将一个 *DirectBuffer 对象放回池中.如果池中该大小的对象数量超过了阈值，则进行校准操作
func (p *Pool) Put(b *DirectBuffer) {
	idx := index(len(b.data))

	// 如果当前池中该大小的对象调用次数超过了阈值，则进行校准操作

	if atomic.AddUint64(&p.calls[idx], 1) > calibrateCallsThreshold {
		p.calibrate()
	}

	maxSize := int(atomic.LoadUint64(&p.maxSize))
	if maxSize == 0 || cap(b.data) <= maxSize {
		b.Reset()
		p.pool.Put(b)
	}
}

func (p *Pool) calibrate() {
	if !atomic.CompareAndSwapUint64(&p.calibrating, 0, 1) {
		return
	}
	// 统计所有 buffer 大小的调用频率
	a := make(callSizes, 0, steps)
	var callsSum uint64
	for i := uint64(0); i < steps; i++ {
		calls := atomic.SwapUint64(&p.calls[i], 0)
		callsSum += calls
		a = append(a, callSize{
			calls: calls,
			size:  minSize << i,
		})
	}
	// 找出最常用的大小区间
	sort.Sort(a) // 从大到小排序

	// 计算默认大小和最大允许回收的大小
	defaultSize := a[0].size
	maxSize := defaultSize

	// 确定最大允许回收的 buffer 大小
	maxSum := uint64(float64(callsSum) * maxPercentile)
	callsSum = 0
	for i := 0; i < steps; i++ {
		if callsSum > maxSum {
			break
		}
		callsSum += a[i].calls
		size := a[i].size
		if size > maxSize {
			maxSize = size
		}
	}

	atomic.StoreUint64(&p.defaultSize, defaultSize)
	atomic.StoreUint64(&p.maxSize, maxSize)

	atomic.StoreUint64(&p.calibrating, 0)
}

type callSize struct {
	calls uint64
	size  uint64
}

type callSizes []callSize

func (ci callSizes) Len() int {
	return len(ci)
}

func (ci callSizes) Less(i, j int) bool {
	return ci[i].calls > ci[j].calls
}

func (ci callSizes) Swap(i, j int) {
	ci[i], ci[j] = ci[j], ci[i]
}

func index(n int) int {
	n--
	n >>= minBitSize
	idx := 0
	for n > 0 {
		n >>= 1
		idx++
	}
	if idx >= steps {
		idx = steps - 1
	}
	return idx
}
