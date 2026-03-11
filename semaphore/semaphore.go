package semaphore

type Semaphore chan struct{}

func NewSempaphore(n int) Semaphore {
	return make(Semaphore, n)
}

func (s Semaphore) Acquire() {
	s <- struct{}{}
}

func (s Semaphore) Release() {
	<-s
}
