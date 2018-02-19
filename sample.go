package hello

import (
	"fmt"
	"github.com/guilhebl/go-worker-pool"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Module struct {
	JobQueue   chan job.Job
	Dispatcher *job.WorkerPool
}

var instance *Module
var once sync.Once

func BuildInstance() *Module {
	once.Do(func() {
		instance = newModule()
	})
	return instance
}

func GetInstance() *Module {
	return instance
}

// Builds a new module which is a container for the running app instance
func newModule() *Module {

	// init worker pool
	numCPUs := runtime.NumCPU()
	maxWorkers := numCPUs
	workerPool := job.NewWorkerPool(maxWorkers)
	jobQueue := make(chan job.Job)

	module := Module{
		Dispatcher: &workerPool,
		JobQueue:   jobQueue,
	}

	// A buffered channel that we can send work requests on.
	module.Dispatcher.Run(jobQueue)

	return &module
}

// represents a task that generates random numbers (producers) and sums the results (consumers)
type RandomIntSumTask struct{}

func (e *RandomIntSumTask) Run(payload job.Payload) job.JobResult {
	x, _ := strconv.ParseInt(payload.Params["x"], 10, 0)
	y, _ := strconv.ParseInt(payload.Params["y"], 10, 0)
	z, _ := strconv.ParseInt(payload.Params["z"], 10, 0)
	return job.NewJobResult(x+y+z, nil)
}

func NewRandomIntSumTask() RandomIntSumTask {
	return RandomIntSumTask{}
}

func handler(w http.ResponseWriter, r *http.Request) {
	// let's create some sample jobs
	work := NewRandomSampleJob()
	m1 := work.Payload.Params
	fmt.Fprintf(w, "worker 1 - numbers: %s, %s, %s \n\n", m1["x"], m1["y"], m1["z"])

	work2 := NewRandomSampleJob()
	m2 := work2.Payload.Params
	fmt.Fprintf(w, "worker 2 - numbers: %s, %s, %s \n\n", m2["x"], m2["y"], m2["z"])

	work3 := NewRandomSampleJob()
	m3 := work3.Payload.Params
	fmt.Fprintf(w, "worker 3 - numbers: %s, %s, %s \n\n", m3["x"], m3["y"], m3["z"])

	work4 := NewRandomSampleJob()
	m4 := work4.Payload.Params
	fmt.Fprintf(w, "worker 4 - numbers: %s, %s, %s \n\n", m4["x"], m4["y"], m4["z"])

	// Push each job onto the queue.
	GetInstance().JobQueue <- work
	GetInstance().JobQueue <- work2
	GetInstance().JobQueue <- work3
	GetInstance().JobQueue <- work4

	// Consume the merged output from all jobs and output matching sum result
	total := int64(1)
	for n := range job.Merge(work.ReturnChannel, work2.ReturnChannel, work3.ReturnChannel, work4.ReturnChannel) {
		result := n.Value
		total *= result.(int64)
	}

	fmt.Fprintf(w, "Total: %d", total)
}

const (
	MaxRand = 100
)

// Returns an int >= min, < max
func randomInt(min, max int) int {
    return min + rand.Intn(max-min)
}

func NewRandomSampleJob() job.Job {
    rand.Seed(time.Now().UnixNano())
	ret := job.NewJobResultChannel()
	m := make(map[string]string)
	x := randomInt(1, MaxRand)
	y := randomInt(1, MaxRand)
	z := randomInt(1, MaxRand)
	m["x"] = strconv.Itoa(x)
	m["y"] = strconv.Itoa(y)
	m["z"] = strconv.Itoa(z)
	task := NewRandomIntSumTask()
	return job.NewJob(&task, m, ret)
}

// entry function called by Google App Engine
func init() {

	// create module
	BuildInstance()

	http.HandleFunc("/", handler)
}
