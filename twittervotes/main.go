package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	nsq "github.com/bitly/go-nsq"
	"gopkg.in/mgo.v2"
)

func main() {
	// create a stop bool with an associated mutex
	// so we can access it from many goroutines at the same time
	var stoplock sync.Mutex
	stop := false
	// signalChan is sent any SIGINT or SIGTERM signals when
	// someone tries to halt the program
	signalChan := make(chan os.Signal, 1)
	// stopChan is passed on to startTwitterStream as a way to
	// tell it to terminate its process when the program is stopped
	stopChan := make(chan struct{}, 1)
	go func() {
		// block until there is a signal on the signalChan
		<-signalChan
		// set stop to true when a signal is received
		stoplock.Lock()
		stop = true
		stoplock.Unlock()
		// stop twitterStream and close the connection
		log.Println("Stopping...")
		stopChan <- struct{}{}
		closeConn()
	}()
	// relays incoming SIGINT and SIGTERM signals to signalChan
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// open and defer the closing of the database connection
	if err := dialdb(); err != nil {
		log.Fatalln("failed to dial MongoDB:", err)
	}
	defer closedb()

	// start things
	// chan for votes
	votes := make(chan string)
	// stop signal channel that lets us know the publisher is stopped
	publisherStoppedChan := publishVotes(votes)
	// stop signal channel that lets us know we stopped reading
	// from the twitter stream
	twitterStoppedChan := startTwitterStream(stopChan, votes)
	go func() {
		// infinite loop that closes the connection every minute
		// (a new connection will be established by the loop inside
		// startTwitterStream)
		for {
			time.Sleep(1 * time.Minute)
			closeConn()
			// check if we should break out of this loop
			stoplock.Lock()
			if stop {
				stoplock.Unlock()
				break
			}
			stoplock.Unlock()
		}
	}()
	// block until the we have stopped reading the twitter stream
	<-twitterStoppedChan
	// close the votes channel, which will signal the publisher
	// goroutine to exit
	close(votes)
	// wait for the publisher to stop before exiting
	<-publisherStoppedChan
}

var db *mgo.Session

func dialdb() error {
	var err error
	log.Println("dialing mongodb: localhost")
	db, err = mgo.Dial("localhost")
	return err
}

func closedb() {
	db.Close()
	log.Println("closed database connection")
}

type poll struct {
	Options []string
}

func loadOptions() ([]string, error) {
	var options []string
	iter := db.DB("ballots").C("polls").Find(nil).Iter()
	var p poll
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}
	iter.Close()
	return options, iter.Err()
}

func publishVotes(votes <-chan string) <-chan struct{} {
	// setup stop channel that will signal the exit of the goroutine
	stopchan := make(chan struct{}, 1)
	// create a new NSQ producer
	pub, _ := nsq.NewProducer("localhost:4150", nsq.NewConfig())
	go func() {
		// range over the votes channel
		for vote := range votes {
			pub.Publish("votes", []byte(vote)) // publish vote
		}
		// when the votes channel is closed, stop publishing
		log.Println("Publisher: Stopping")
		pub.Stop()
		log.Println("Publisher: Stopped")
		// signal the stop channel that this goroutine is exiting.
		// this could alternatively be done in a deferred function
		stopchan <- struct{}{}
	}()
	return stopchan
}
