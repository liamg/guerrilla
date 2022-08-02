package guerrilla

import (
	"sync"
	"time"
)

type Poller interface {
	Poll() <-chan Email
	Close()
	Error() error
}

type poller struct {
	sync.Mutex
	client       Client
	channel      chan Email
	err          error
	wg           sync.WaitGroup
	busy         bool
	closeChan    chan struct{}
	pollInterval time.Duration
}

var _ Poller = (*poller)(nil)

type PollOption func(*poller)

func PollOptionWithInterval(interval time.Duration) PollOption {
	return func(p *poller) {
		p.pollInterval = interval
	}
}

func NewPoller(c Client, options ...PollOption) Poller {
	p := &poller{
		client:       c,
		channel:      make(chan Email),
		closeChan:    make(chan struct{}),
		pollInterval: time.Second * 30,
	}
	for _, option := range options {
		option(p)
	}
	return p
}

func (p *poller) Poll() <-chan Email {
	p.Lock()
	defer p.Unlock()
	if p.err != nil {
		return nil
	}
	if p.busy {
		return p.channel
	}
	p.busy = true
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			summaries, err := p.client.GetNewEmails()
			if err != nil {
				p.setError(err)
				return
			}
			for _, summary := range summaries {
				email, err := p.client.GetEmail(summary.ID)
				if err != nil {
					p.setError(err)
					return
				}
				p.channel <- *email
			}
			timer := time.NewTimer(p.pollInterval)
			select {
			case <-p.closeChan:
				timer.Stop()
				return
			case <-timer.C:
			}
		}
	}()
	return p.channel
}

func (p *poller) Close() {
	p.Lock()
	defer p.Unlock()
	if !p.busy {
		return
	}
	close(p.closeChan)
	p.wg.Wait()
	close(p.channel)
}

func (p *poller) setError(err error) {
	defer p.Close()
	p.Lock()
	defer p.Unlock()
	p.err = err
}

func (p *poller) Error() error {
	p.Lock()
	defer p.Unlock()
	return p.err
}
