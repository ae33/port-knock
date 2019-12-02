package knock

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ae33/port-knock/config"
)

func UdpKnocks(c config.Config) error {
	var knocksPerformed int
	knockSignals := make(chan struct{})

	// Use channels to synchronize the ordering of the ports to knock. Multiple
	// goroutines will wait their "turn" to knock their port.
	knockTurns := make([]chan struct{}, len(c.Ports))
	for i := range knockTurns {
		knockTurns[i] = make(chan struct{})
	}

	// Signal that it's the first knock's turn.
	go func() {
		knockTurns[0] <- struct{}{}
	}()

	// Count the number of knocks performed, and control who's turn it is to
	// knock, next.
	go func(k chan struct{}, numDesired int) {
		for knocksPerformed < numDesired {
			<-k
			knocksPerformed = knocksPerformed + 1
			if knocksPerformed < numDesired {
				knockTurns[knocksPerformed] <- struct{}{}
			}
		}
		close(k)
	}(knockSignals, len(c.Ports))

	// Spin up the goroutines that will wait their turn to knock their
	// individual ports.
	for i, p := range c.Ports {
		go func(k chan struct{}, port uint16, ourTurn int) {
			<-knockTurns[ourTurn]

			err := udpKnock(c.Host, port)
			if err != nil {
				log.Printf("error knocking %v at port %d: '%v'", c.Host, p, err)
			}

			k <- struct{}{}
		}(knockSignals, p, i)
	}

	elapsedTime := 0 * time.Nanosecond
	for knocksPerformed < len(c.Ports) {
		if c.QuitAfter == nil {
			return errors.New("nil `quit_after` found in config")
		}

		if elapsedTime > *c.QuitAfter {
			return fmt.Errorf("more than max of %v has elapsed while waiting for udp knocks to happen", *c.QuitAfter)
		}

		time.Sleep(*c.WaitSleep)
		elapsedTime = elapsedTime + *c.WaitSleep
	}

	return nil
}

func udpKnock(host string, port uint16) error {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return fmt.Errorf("error dialing udp connection: '%v'", err)
	}

	bytesWritten, err := conn.Write([]byte("k"))
	if err != nil {
		return fmt.Errorf("error writing udp: '%v'", err)
	}

	pluralControl := "bytes"
	if bytesWritten == 1 {
		pluralControl = "byte"
	}
	log.Printf("%d %s written to %s:%d", bytesWritten, pluralControl, host, port)

	return nil
}
