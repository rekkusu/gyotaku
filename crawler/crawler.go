package crawler

import (
	"errors"
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

const (
	CrawlURL = "http://127.0.0.1:9999/admin_53cr37api/"
	NewURL   = "http://127.0.0.1:9999/new"
)

var (
	ChromePath string
	Flag       string
	CrawlQueue chan int = make(chan int, 256)
	pool       chan *worker
)

type worker struct {
	port int
	quit chan interface{}
	data chan int
}

func NewWorker(port int) *worker {
	return &worker{
		port: port,
		quit: make(chan interface{}),
		data: make(chan int),
	}
}

func StartCrawler(thread int) {
	pool = make(chan *worker, thread)
	workers := make([]*worker, thread)

	go func() {
		for {
			select {
			case id := <-CrawlQueue:
				(<-pool).data <- id
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
	flagUrl := fmt.Sprintf("%s?url=%s", NewURL, Flag)
	cmd := exec.CommandContext(ctx, ChromePath, "--headless", "--disable-gpu", "--no-referrers", "--remote-debugging-port="+strconv.Itoa(w.port), flagUrl)
	cmd.Start()

	go func() {
		for {
			pool <- w

			select {
			case id := <-w.data:
				w.Run(id)
			case <-w.quit:
				return
			}
		}
	}()
}

func (w *worker) Run(id int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	finish := make(chan interface{})

	go func() {
		dt := devtool.New(fmt.Sprintf("http://127.0.0.1:%d", w.port))
		pt, err := w.getTarget(ctx, dt)
		if err != nil {
			log.Println(err)
			return
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
			log.Println(err)
		}
		defer domContent.Close()

		if err = c.Page.Enable(ctx); err != nil {
			log.Println(err)
		}

		url := fmt.Sprintf("%s%d", CrawlURL, id)
		c.Page.Navigate(ctx, page.NewNavigateArgs(url))

		domContent.Recv()

		finish <- struct{}{}
	}()

	select {
	case <-finish:
	case <-ctx.Done():
	}
}

func (w *worker) getTarget(ctx context.Context, dt *devtool.DevTools) (*devtool.Target, error) {
	targetCh := make(chan *devtool.Target)

	go func() {
		for i := 0; i < 3; i++ {
			pt, err := dt.Get(ctx, devtool.Page)
			if err != nil {
				log.Println(err)
				pt, err = dt.Create(ctx)
				if err != nil {
					log.Println(err)
					return
				}
			}

			if pt.WebSocketDebuggerURL != "" {
				targetCh <- pt
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	select {
	case t := <-targetCh:
		return t, nil
	case <-ctx.Done():
		return nil, errors.New("timeout")
	}
}

func (w *worker) Stop() {
	w.quit <- struct{}{}
}
