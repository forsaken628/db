package main // import "github.com/forsaken628/db"

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
	"math/rand"
	"github.com/openzipkin/zipkin-go-opentracing"
	"fmt"
	"os"
	"github.com/opentracing/opentracing-go"
)

var wg sync.WaitGroup
var db *sql.DB

const n = 50

func init() {
	db0, err := sql.Open("mysql", "root@/test")
	if err != nil {
		panic(err)
	}
	err = db0.Ping()
	if err != nil {
		panic(err)
	}
	db = db0
	rand.Seed(time.Now().UnixNano())

	collector := NewLogCollector()

	// Create our recorder.
	recorder := zipkintracer.NewRecorder(
		collector,
		false,
		"",
		"test",
	)

	// Create our tracer.
	tracer, err := zipkintracer.NewTracer(
		recorder,
		zipkintracer.ClientServerSameSpan(true),
		zipkintracer.TraceID128Bit(true),
		//todo zipkintracer.WithLogger()
	)

	if err != nil {
		fmt.Printf("unable to create Zipkin tracer: %+v\n", err)
		os.Exit(1)
	}

	opentracing.SetGlobalTracer(tracer)
}

func main() {
	case2()
	wg.Wait()
}
