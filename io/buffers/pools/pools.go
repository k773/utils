package pools

import (
	"github.com/k773/utils/io/buffers/readerWriterAt"
	"sync"
)

/*
	readerWriterAt.SimpleReaderWriterAt
*/

// Pool item

type PoolItemSimpleReaderWriterAt struct {
	*readerWriterAt.SimpleReaderWriterAt
}

func (p *PoolItemSimpleReaderWriterAt) Release() {
	p.Reset()
	PoolSimpleReaderWriterAt.pool.Put(p)
}

// Pool

type poolSimpleReaderWriterAt struct {
	pool *sync.Pool
}

func (p *poolSimpleReaderWriterAt) Get() *PoolItemSimpleReaderWriterAt {
	return p.pool.Get().(*PoolItemSimpleReaderWriterAt)
}

var PoolSimpleReaderWriterAt = &poolSimpleReaderWriterAt{&sync.Pool{New: func() any { return &PoolItemSimpleReaderWriterAt{readerWriterAt.NewSimpleReaderWriterAt(nil)} }}}
