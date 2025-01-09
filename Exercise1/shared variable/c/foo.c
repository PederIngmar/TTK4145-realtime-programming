// Compile with `gcc foo.c -Wall -std=gnu99 -lpthread`, or use the makefile
// The executable will be named `foo` if you use the makefile, or `a.out` if you use gcc directly

#include <pthread.h>
#include <stdio.h>

int i = 0;
pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER; // Using a mutex to protect the shared variable, and not semaphore as it is a more general 

// Incrementing function
void *incrementingThreadFunction(void *arg) {
    for (int j = 0; j < 1000000; ++j) {
        pthread_mutex_lock(&mutex);
        i++;
        pthread_mutex_unlock(&mutex);
    }
    return NULL;
}

// Decrementing function
void *decrementingThreadFunction(void *arg) {
    for (int j = 0; j < 1000001; ++j) {
        pthread_mutex_lock(&mutex);
        i--;
        pthread_mutex_unlock(&mutex);
    }
    return NULL;
}

int main() {
    pthread_t incrementingThread, decrementingThread;

    pthread_create(&incrementingThread, NULL, incrementingThreadFunction, NULL);
    pthread_create(&decrementingThread, NULL, decrementingThreadFunction, NULL);

    pthread_join(incrementingThread, NULL);
    pthread_join(decrementingThread, NULL);

    printf("Magic value i: %d\n", i);

    pthread_mutex_destroy(&mutex);

    return 0;
}
