package dispatcher

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

type Client struct {
	ID      string
	MsgChan chan string
}

type Dispatcher struct {
	ctx     context.Context
	joinCh  chan Client
	leaveCh chan Client
	msgCh   chan string
}

var (
	clientsGauge  = prometheus.NewGauge(prometheus.GaugeOpts{Name: "chat_clients_connected", Help: "Active clients"})
	messagesTotal = prometheus.NewCounter(prometheus.CounterOpts{Name: "chat_messages_total", Help: "Total broadcast messages"})
)

func init() {
	if err := prometheus.Register(clientsGauge); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			clientsGauge = are.ExistingCollector.(prometheus.Gauge)
		} else {
			log.Fatalf("unable to register metric: %v", err)
		}
	}
}

func New(ctx context.Context) *Dispatcher {
	return &Dispatcher{
		ctx:     ctx,
		joinCh:  make(chan Client),
		leaveCh: make(chan Client),
		msgCh:   make(chan string, 10),
	}
}

func (d *Dispatcher) Run() {
	clients := make(map[string]chan string)
	for {
		select {
		case c := <-d.joinCh:
			clientsGauge.Inc()
			clients[c.ID] = c.MsgChan
		case c := <-d.leaveCh:
			clientsGauge.Dec()
			delete(clients, c.ID)
			close(c.MsgChan)
		case msg := <-d.msgCh:
			messagesTotal.Inc()
			for _, out := range clients {
				out <- msg
			}
		case <-d.ctx.Done():
			for _, out := range clients {
				close(out)
			}
			return
		}
	}
}

func (d *Dispatcher) Join(c Client)        { d.joinCh <- c }
func (d *Dispatcher) Leave(c Client)       { d.leaveCh <- c }
func (d *Dispatcher) Broadcast(msg string) { d.msgCh <- msg }
