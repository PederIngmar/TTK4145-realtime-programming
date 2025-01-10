Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?
Concurrency is about processing different taskes almost at the same time, on the same core or on different cores by switching between them.
Parallelism is about processing different taskes at the same time on differnt cores

What is the difference between a *race condition* and a *data race*? 
Both can happend when concurrent threads operate on the same memory location

Race condition: faulty behaviour that can happen when the timing/order of execution of many processes happening separately matters, but is not properly synchronized.
Data race: faulty behaviour that can happen when two or more processes happening separately is accessing the same memory location at the same time, and at least one of them is editing the content at that location. Data race is a specific type of race condition.

*Very* roughly - what does a *scheduler* do, and how does it do it?
A scheduler determines when different processes/threads are give runtime on the CPU depending og the needs of the system.


### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?
The main advantage of using threads is it allows processes to run concurrently, which is an advantage in realtime systems. When on process is waiting for something in to happen (some input for exsample) another process can use the CPU. Threads also allow for parallelism which can speed up processes.

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?
Fibers/coroutines are managed by the user space rather then the OS kernel (as for threaads). One would use the over threads because of: 
- faster context switching because it avoids krenal overhead
- lower memory consumption because they use smaller stacks and dont rely on OS managed thread pools
- scalability because many of them can run on one single OS thread

Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
It makes the programmers life harder in the sens that procresses needs to be synchronized which can be very complex.

What do you think is best - *shared variables* or *message passing*?
My current view is that message passing is better because it is simple but powerfull because it scales very well. 

