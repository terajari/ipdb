package main

import (
	"container/heap"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func (app *application) recoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header("Connection", "close")
				app.serverErrorResponse(c, fmt.Errorf("%s", err))
				c.Abort()
				return
			}
		}()
		c.Next()
	}
}

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type clientHeap []*client

func (h clientHeap) Len() int           { return len(h) }
func (h clientHeap) Less(i, j int) bool { return h[i].lastSeen.Before(h[j].lastSeen) }
func (h clientHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *clientHeap) Push(x interface{}) {
	*h = append(*h, x.(*client))
}

func (h *clientHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

var (
	mu          sync.Mutex
	clients     = sync.Map{}
	heapClients clientHeap
)

func init() {
	heapClients = make(clientHeap, 0)
	heap.Init(&heapClients)

	go cleanUpClients()
}

func cleanUpClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		now := time.Now()

		for heapClients.Len() > 0 {
			c := heapClients[0]
			if now.Sub(c.lastSeen) > 3*time.Minute {
				heap.Pop(&heapClients)
				clients.Delete(c)
			} else {
				break
			}
		}

		mu.Unlock()
	}
}

func (app *application) rateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if app.config.limiter.enabled {
			host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(c, err)
				c.Abort()
				return
			}

			val, ok := clients.Load(host)
			var cl *client

			if !ok {
				lim := rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst)
				cl = &client{
					limiter:  lim,
					lastSeen: time.Now(),
				}
				mu.Lock()
				heap.Push(&heapClients, cl)
				clients.Store(host, cl)
				mu.Unlock()
			} else {
				cl = val.(*client)
				cl.lastSeen = time.Now()
			}

			if !cl.limiter.Allow() {
				app.rateLimitExceededResponse(c)
				c.Abort()
				return
			}
			c.Next()
		}
	}
}
