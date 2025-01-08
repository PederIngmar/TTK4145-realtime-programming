// Compile with `gcc foo.c -Wall -std=gnu99 -lpthread`, or use the makefile
// The executable will be named `foo` if you use the makefile, or `a.out` if you use gcc directly

#include <pthread.h>
#include <stdio.h>

int i = 0;

// Note the return type: void*
void* incrementingThreadFunction(){
    // TODO: increment i 1_000_000 times
    for (int j = 0; j < 1000000; j++) {
        i++;
    }
    return NULL;
}

void* decrementingThreadFunction(){
    // TODO: decrement i 1_000_000 times
    for (int j = 0; j < 1000000; j++) {
        i--;
    }
    return NULL;
}


int main(){
    // TODO: 
    // start the two functions as their own threads using `pthread_create`
    // Hint: search the web! Maybe try "pthread_create example"?
    
    pthread_t incrementingThread; // pthread_t is a type that represents a thread
    pthread_t decrementingThread; 
    pthread_create(&incrementingThread, NULL, incrementingThreadFunction, NULL); // pthread_create creates a new thread
    pthread_create(&decrementingThread, NULL, decrementingThreadFunction, NULL);

    // TODO:
    // wait for the two threads to be done before printing the final result
    // Hint: Use `pthread_join`  
    pthread_join(incrementingThread, NULL); // pthread_join waits for the thread to finish. The thread and the return value is joined
    pthread_join(decrementingThread, NULL);


    printf("The magic number is: %d\n", i);
    return 0;
}
#include <pthread.h>
#include <stdio.h>

// Global variable and mutex
int i = 0;
pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER;

// Incrementing function
void *incrementingThreadFunction(void *arg) {
    for (int j = 0; j < 1000000; ++j) {
        pthread_mutex_lock(&mutex); // Lock the mutex
        i++;
        pthread_mutex_unlock(&mutex); // Unlock the mutex
    }
    return NULL;
}

// Decrementing function
void *decrementingThreadFunction(void *arg) {
    for (int j = 0; j < 1000000; ++j) {
        pthread_mutex_lock(&mutex); // Lock the mutex
        i--;
        pthread_mutex_unlock(&mutex); // Unlock the mutex
    }
    return NULL;
}

int main() {
    pthread_t incrementingThread, decrementingThread;

    pthread_create(&incrementingThread, NULL, incrementingThreadFunction, NULL);
    pthread_create(&decrementingThread, NULL, decrementingThreadFunction, NULL);

    pthread_join(incrementingThread, NULL);
    pthread_join(decrementingThread, NULL);

    printf("Final value of i: %d\n", i);

    // Destroy the mutex
    pthread_mutex_destroy(&mutex);

    return 0;
}
