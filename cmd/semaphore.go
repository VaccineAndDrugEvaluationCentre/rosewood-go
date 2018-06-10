package main

//TODO: move to util package
type semToken struct{} //empty because we only care about the count
//CountingSemaphore implemented as a buffered channel of empty structs
type CountingSemaphore chan semToken

//NewCountingSemaphore returns a counting semaphore with a certain capacity
func NewCountingSemaphore(capacity int) CountingSemaphore {
	if capacity < 1 {
		capacity = 1
	}
	return make(CountingSemaphore, capacity)
}

//Reserve acquires a token by putting a token in the channel's buffer (taking up a spot)
func (cs CountingSemaphore) Reserve(count int) {
	for i := 0; i < count; i++ {
		cs <- semToken{}
	}
}

//Free releases a token by reading from the channel's buffer (freeing up a spot)
func (cs CountingSemaphore) Free(count int) {
	for i := 0; i < count; i++ {
		<-cs
	}
}
