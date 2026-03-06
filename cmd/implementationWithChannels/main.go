package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type blockingQueue struct {
	channel chan int32
}


// initialisation, enqueue and dequeue operations are needed

func initialise(capacity int32) *blockingQueue{
	q := &blockingQueue{
		channel : make(chan int32,capacity),
	}

	return q
}

func (q *blockingQueue) enqueue (item int32) {
	q.channel<-item
}

func (q *blockingQueue) dequeue () int32 {
	return <-q.channel	
}

func (q *blockingQueue) size () int32 {
	return int32(len(q.channel))
}

func main(){
	q := initialise(32)
	var WgE,WgD sync.WaitGroup
	for i:=0;i<1000000;i++ {
		WgE.Add(1)
		go func(){
			defer WgE.Done()
			q.enqueue(rand.Int31())
		}()
	}

	for i:=0;i<1000000;i++ {
		WgD.Add(1)
		go func(){
			defer WgD.Done()
			q.dequeue()
		}()
	}
	WgE.Wait()
	WgD.Wait()

	fmt.Println("size of queue is: ", q.size())
}