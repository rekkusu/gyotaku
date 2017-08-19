package crawler

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"

	"golang.org/x/net/context"
)

const CrawlURL = "http://127.0.0.1:9999/admin_53cr37api/"

var (
	ChromePath string
	CrawlQueue chan interface{} = make(chan interface{}, 256)
	pool       chan *worker
)

type worker struct {
	port int
	quit chan interface{}
	data chan interface{}
}

func NewWorker(port int) *worker {
	return &worker{
		port: port,
		quit: make(chan interface{}),
		data: make(chan interface{}),
	}
}

func StartCrawler(thread int) {
	pool = make(chan *worker, thread)
	workers := make([]*worker, thread)

	go func() {
		for {
			select {
			case <-CrawlQueue:
				(<-pool).data <- struct{}{}
			}
		}
	}()

	for i := 0; i < thread; i++ {
		workers[i] = NewWorker(9222 + i)
		workers[i].Start()
		log.Printf("New worker port: %d\n", 9222+i)
	}
}

func (w *worker) Start() {
	ctx, _ := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--headless", "--disable-gpu", "--remote-debugging-port="+strconv.Itoa(w.port))
	cmd.Start()

	go func() {
		for {
			pool <- w

			select {
			case <-w.data:
				w.Run()
			case <-w.quit:
				return
			}
		}
	}()
}

func (w *worker) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dt := devtool.New(fmt.Sprintf("http://127.0.0.1:%d", w.port))
	pt, err := dt.Get(ctx, devtool.Page)
	if err != nil {
		log.Println(1, err)
		pt, err = dt.Create(ctx)
		if err != nil {
			log.Println(2, err)
			return
		}
	}

	conn, err := rpcc.DialContext(ctx, pt.WebSocketDebuggerURL)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	c := cdp.NewClient(conn)

	domContent, err := c.Page.DOMContentEventFired(ctx)
	if err != nil {
		log.Println(4, err)
	}
	defer domContent.Close()

	if err = c.Page.Enable(ctx); err != nil {
		log.Println(err)
	}

	c.Page.Navigate(ctx, page.NewNavigateArgs(CrawlURL))

	domContent.Recv()
}

func (w *worker) Stop() {
	w.quit <- struct{}{}
}
