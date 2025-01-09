// compile with:  gcc -g main.c ringbuf.c -lpthread

#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>
#include <semaphore.h>
#include <time.h>

#include "ringbuf.h"

struct BoundedBuffer {
    struct RingBuffer*  buf;
    pthread_mutex_t     mtx;
    sem_t               numElements;
    sem_t               capacity;
    
    
};

struct BoundedBuffer* buf_new(int size){
    struct BoundedBuffer* buf = malloc(sizeof(struct BoundedBuffer));
    buf->buf = rb_new(size);
    
    pthread_mutex_init(&buf->mtx, NULL);
    // TODO: initialize semaphores
    //sem_init(&buf->capacity,      0, /*starting value?*/);
	//sem_init(&buf->numElements,   0, /*starting value?*/);
    sem_init(&buf->capacity, 0, size);
    sem_init(&buf->numElements, 0, 0);
    
    return buf;    
}

void buf_destroy(struct BoundedBuffer* buf){
    rb_destroy(buf->buf);
    pthread_mutex_destroy(&buf->mtx);
    sem_destroy(&buf->numElements);
    sem_destroy(&buf->capacity);
    free(buf);
}




void buf_push(struct BoundedBuffer* buf, int val){    
    
    sem_wait(&buf->capacity);       // wait until there is space in the buffer

    pthread_mutex_lock(&buf->mtx); // lock the buffer
    rb_push(buf->buf, val);        // push the value into the buffer
    pthread_mutex_unlock(&buf->mtx); // unlock the buffer
    
    sem_post(&buf->numElements);   // signal that there are new elements in the buffer
}

int buf_pop(struct BoundedBuffer* buf){
    // TODO: same, but different?

    sem_wait(&buf->numElements);    // wait until there are elements in the buffer
    
    pthread_mutex_lock(&buf->mtx);  // lock the buffer
    int val = rb_pop(buf->buf);     // pop the value from the buffer
    pthread_mutex_unlock(&buf->mtx); // unlock the buffer

    sem_post(&buf->capacity);       // signal that there is space in the buffer
    
    return val;
}





void* producer(void* args){
    struct BoundedBuffer* buf = (struct BoundedBuffer*)(args);
    
    for(int i = 0; i < 10; i++){
        nanosleep(&(struct timespec){0, 100*1000*1000}, NULL);
        printf("[producer]: pushing %d\n", i);
        buf_push(buf, i);
    }
    return NULL;
}

void* consumer(void* args){
    struct BoundedBuffer* buf = (struct BoundedBuffer*)(args);
    
    // give the producer a 1-second head start
    nanosleep(&(struct timespec){1, 0}, NULL);
    while(1){
        int val = buf_pop(buf);
        printf("[consumer]: %d\n", val);
        nanosleep(&(struct timespec){0, 50*1000*1000}, NULL);
    }
}

int main(){ 
    
    struct BoundedBuffer* buf = buf_new(5);
    
    pthread_t producer_thr;
    pthread_t consumer_thr;
    pthread_create(&producer_thr, NULL, producer, buf);
    pthread_create(&consumer_thr, NULL, consumer, buf);
    
    pthread_join(producer_thr, NULL);
    pthread_cancel(consumer_thr);
    
    buf_destroy(buf);
    
    return 0;
}
