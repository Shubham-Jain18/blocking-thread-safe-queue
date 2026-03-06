# Blocking Thread-Safe Queue

This repository implements solutions for the classic producer consumer problem.

## Thread-Safe Queue vs. Blocking Queue

A thread-safe queue prevents race conditions and data corruption when accessed concurrently by multiple threads. Operations on a standard thread-safe queue return immediately. If the queue is full on enqueue or empty on dequeue, the operation fails or returns a boolean indicator.

A blocking queue suspends the execution of the calling thread if it cannot complete the operation immediately. The advantage of a blocking queue over a non-blocking thread-safe queue is the elimination of busy waiting. Threads do not waste CPU cycles continuously polling the queue to check if space or data has become available.

## Condition Variables and Thread Sleep

Condition variables provide a mechanism for threads to sleep while waiting for a specific state change. 

When a thread attempts an operation that cannot proceed, it invokes a wait method on the condition variable. The operating system removes the thread from the CPU and places it into a sleep state. The thread is added to a specific waiting queue tied to that condition variable. Simultaneously, the thread releases the lock it holds, allowing other threads to access the shared resource. 

While in the waiting queue, the thread consumes zero CPU resources. Once another thread modifies the shared state, it signals the condition variable. The OS selects a thread from the waiting queue and wakes it up. The awakened thread attempts to reacquire the lock. Once the lock is secured, it checks the condition again and proceeds.

## Implementations

### Naive Implementation (Mutexes and Condition Variables)

This implementation uses explicit synchronization primitives to guard a slice-based queue.

* The `blockingQueue` structure contains an integer slice, a `sync.Mutex`, and two `sync.Cond` pointers named `notFull` and `notEmpty`.
* Enqueue operations lock the mutex and evaluate the queue length.
* If the queue length meets or exceeds capacity, the thread waits on the `notFull` condition variable.
* Dequeue operations wait on the `notEmpty` condition variable while the queue length equals zero.
* Wait statements are enclosed within `for` loops to handle spurious wakeups and lock thefts from competing threads.
* After appending an item to the slice, the `notEmpty` variable is signaled.
* After slicing an item from the front of the queue, the `notFull` variable is signaled.

### Channel Implementation

This implementation delegates synchronization and blocking to Go's native channel constructs.

* The `blockingQueue` structure consists solely of a buffered `chan int32`.
* The enqueue function performs a standard channel send operation.
* The dequeue function performs a standard channel receive operation.
* The buffered channel natively handles suspending the goroutine if it is full during a send or empty during a receive.

## Channels vs. Mutexes in Go

Channels are the preferred construct in Go for transferring data ownership and orchestrating execution flow between goroutines. For the producer consumer problem, channels abstract away the complexity of condition variables, manual locking, and signaling. 

Channels are superior in aspects of readability, safety, and built-in messaging mechanics. They prevent common errors associated with manual lock management. Mutexes remain superior for synchronizing access to granular internal state or executing simple atomic mutations where channel communication overhead would degrade performance.
