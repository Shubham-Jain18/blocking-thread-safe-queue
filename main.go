package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
)


type blockingQueue struct{
	queue []int32
	mu sync.Mutex
	notFull *sync.Cond
	notEmpty *sync.Cond
	capacity int32
}


// need functions for initialisation, enqueue and dequeue operation

func initialise(capacity int32) *blockingQueue {
	q:= &blockingQueue{
		queue : make([]int32,0,capacity),
		capacity: capacity,
	}
	q.notFull = sync.NewCond(&q.mu)
	q.notEmpty = sync.NewCond(&q.mu)

	return q
}

func (q *blockingQueue) enqueue (item int32) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// using for loop instead of if to handle spurious wakeups and lock thefts from other threads
	for len(q.queue)>=int(q.capacity){
		q.notFull.Wait()
	}

	q.queue = append(q.queue, item)
	q.notEmpty.Signal()

	log.Println(item, " inserted in the queue")

}

func (q *blockingQueue) dequeue () int32 {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.queue) == 0 {
		q.notEmpty.Wait()
	}

	item := q.queue[0]
	q.queue = q.queue[1:]
	q.notFull.Signal()

	log.Println(item, " dequeued")

	return item
}

func (q *blockingQueue) size () int32 {
	q.mu.Lock()
	defer q.mu.Unlock()

	return int32(len(q.queue))
}

func main(){
	q:= initialise(32)

	var WgE,WgD sync.WaitGroup

	for i:=0; i<1000000; i++ {
		WgE.Add(1)
		go func(){
			q.enqueue(rand.Int31())
			WgE.Done()
		}()
	}

	for i:=0; i<1000000; i++ {
		WgD.Add(1)
		go func(){
			q.dequeue()
			WgD.Done()
		}()
	}

	WgE.Wait()
	WgD.Wait()

	fmt.Println("size of queue is: ", q.size())

	
}
