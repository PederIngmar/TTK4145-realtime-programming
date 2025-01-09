## Results

### Task 3
Each run of the programmes "foo.c" and "foo.go" prints a different result between -1000000 and 1000000, even after compilation.
Both threads execute independently, but they operate on the same memory location. The 2 functions consists of some atomic operations like:
- reading the value from memory
- modify the value
- write the updated value back to memory.

If the two threads interleave these steps, the final value of i can be incorrect because one thread may overwrite changes made by the other.


GOMAXPROCS (used in foo.go) decides how many CPU cores the program uses. If I set it to 1 the program only uses one core, so no atomic operations will overlap, and the result will therefore always be 0. But the program will run slower.