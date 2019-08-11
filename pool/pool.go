package pool

import "github.com/discordapp/lilliput"

func NewLilliputPool(count int,  opsSize int, bufferSize int) *LilliputPool {


	pool := make(chan *OpsPair, count)

	for i:=0;i<count;i++ {

		pool <- &OpsPair{
			Ops: lilliput.NewImageOps(opsSize),
			Buf: make([]byte, bufferSize, bufferSize),
		}
	}

	return &LilliputPool{
		pool,
	}

}

type OpsPair struct {
	Ops *lilliput.ImageOps
	Buf []byte
}

type LilliputPool struct {

	pool chan *OpsPair

}

func (l *LilliputPool) Get() *OpsPair {

	return <- l.pool

}

func (l *LilliputPool) Put(o *OpsPair) {

	o.Buf = o.Buf[:0]

	l.pool <- o
}
