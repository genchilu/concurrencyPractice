package ringbufpadded

import (
	"runtime"
	"sync/atomic"
)

func roundUp(v uint64) uint64 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v |= v >> 32
	v++
	//if v in edge, return default 1024
	if v == 0 {
		return 1024
	}
	return v
}

type entity struct {
	position uint64
	data     interface{}
}

type RingBufferPadded struct {
	head    uint64
	_p1     [8]uint64
	tail    uint64
	_p2     [8]uint64
	mask    uint64
	_p3     [8]uint64
	ringbuf []*entity
	_p4     [8]uint64
}

func NewRingBufferPadded(size uint64) *RingBufferPadded {
	rb := &RingBufferPadded{}
	rb.init(size)
	return rb
}

func (rb *RingBufferPadded) init(size uint64) {
	size = roundUp(size)
	rb.ringbuf = make([]*entity, size, size)
	for i := uint64(0); i < size; i++ {
		rb.ringbuf[i] = &entity{position: i}
	}
	rb.mask = size - 1
}

func (rb *RingBufferPadded) Put(item interface{}) error {
	var ent *entity
	pos := atomic.LoadUint64(&rb.head)
	i := 0
L:
	for {
		ent = rb.ringbuf[pos&rb.mask]
		seq := atomic.LoadUint64(&ent.position)
		switch diff := seq - pos; {
		case diff == 0:
			if atomic.CompareAndSwapUint64(&rb.head, pos, pos+1) {
				break L
			}
		case diff > 0:
			pos = atomic.LoadUint64(&rb.head)
		default:
			panic("error while putting item into RingBufferPadded")

		}
		if i == 10000 {
			runtime.Gosched() // free up the cpu before the next iteration
			i = 0
		} else {
			i++
		}
	}
	ent.data = item
	atomic.StoreUint64(&ent.position, pos+1)
	return nil
}

func (rb *RingBufferPadded) Get() (interface{}, error) {
	var ent *entity
	pos := atomic.LoadUint64(&rb.tail)
	i := 0
L:
	for {
		ent = rb.ringbuf[pos&rb.mask]
		seq := atomic.LoadUint64(&ent.position)
		switch diff := seq - (pos + 1); {
		case diff == 0:
			if atomic.CompareAndSwapUint64(&rb.tail, pos, pos+1) {
				break L
			}
		case diff > 0:
			pos = atomic.LoadUint64(&rb.tail)
		default:
			panic("error while getting item into RingBufferPadded")
		}
		if i == 10000 {
			runtime.Gosched()
			i = 0
		} else {
			i++
		}
	}
	data := ent.data
	atomic.StoreUint64(&ent.position, pos+rb.mask+1)
	return data, nil
}
