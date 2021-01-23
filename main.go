package main 
 
import ( 
    "fmt"
    "sync"
	"time"
)

type workerChannel chan uint64
var tasksDone = make(chan string)
 
func printTime(msg string) { 
    fmt.Println(msg, time.Now().Format("15:04:05")) 
} 
 
// Task that will be done over time 
func worker1(wg *sync.WaitGroup, worker string, w uint64) {
    msg := fmt.Sprintf("worker %s working on task %d", worker, w)
    printTime(msg) 
    tasksDone <- fmt.Sprintf("worker %s finish task %d", worker, w)
    wg.Done() 
} 

func worker2(wg *sync.WaitGroup, worker string, w uint64) { 
    msg := fmt.Sprintf("worker %s working on %d", worker, w)
    printTime(msg) 
    tasksDone <- fmt.Sprintf("worker %s finish task %d", worker, w)
    wg.Done() 
} 

func producer(wg *sync.WaitGroup) { 
    for taskDone := range tasksDone {
        fmt.Println(taskDone)
    }
    wg.Done() 
} 

func consumer(workCh chan uint64) { 
    i := 0
    for {
        time.Sleep(time.Second * 5)
        workCh <- uint64(i)
        i++
    }
} 
 
// Task done in parallel 
func listenForever() { 
    for { 
        // printTime("Listening...") 
    } 
} 
 
func main() { 

    num_routines := 3
    workCh := make(workerChannel)
    
    var waitGroup sync.WaitGroup 
    waitGroup.Add(num_routines) 
 
    go listenForever() 
 
    // Give some time for listenForever to start 
    time.Sleep(time.Nanosecond * 10) 

    // Send work to be done
    go consumer(workCh) 
 
    // Each worker work on a task in parrallel and print out completed task
    go func() {
        for w := range workCh {
            go worker1(&waitGroup, "worker1", w) 
            go worker2(&waitGroup, "worker2", w) 
            go producer(&waitGroup) 
            waitGroup.Add(num_routines)
        }
	}()

	waitGroup.Wait()
}