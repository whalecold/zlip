// When encoding or decoding the input file, the file may be too large for the memory to load all data, so we separate the whole
// data into many small chunks. Every chunk is considered as a task and the scheduler dispatches processors to deal with the task,
// the processors performer task parallel which can improve effectiveness. After processing the data,
// scheduler will collects and sorts out the result. The out data may be out of order as the parallel performing,
// we need mark sequence to every task for ordering result more conveniently.
package scheduler
