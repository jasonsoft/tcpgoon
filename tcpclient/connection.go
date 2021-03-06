package tcpclient

import (
	"fmt"
	"time"
)

type Connection struct {
	ID      int
	status  ConnectionStatus
	metrics connectionMetrics
}

type ConnectionStatus int

type connectionMetrics struct {
	tcpEstablishedDuration time.Duration
	tcpErroredDuration     time.Duration
	// packets lost, retransmissions and other metrics could come
}

// Types of connection status
const (
	ConnectionNotInitiated ConnectionStatus = iota + 0
	ConnectionDialing
	ConnectionEstablished
	ConnectionClosed
	ConnectionError
)

// ConnectionFunc type to use connection functions as an argument
type ConnectionFunc func(Connection) bool

// NewConnection initializes a connection given all values that are actually stored internally. This is just being used
// as the first (dirty) approach for tests
func NewConnection(id int, status ConnectionStatus, procTime time.Duration) Connection {
	return Connection{
		ID:     id,
		status: status,
		metrics: connectionMetrics{
			tcpEstablishedDuration: procTime,
			tcpErroredDuration:     procTime,
		},
	}

}

func (c Connection) GetConnectionStatus() ConnectionStatus {
	return c.status
}

func (c Connection) String() string {
	var status string
	switch c.status {
	case ConnectionNotInitiated:
		status = "not initiated"
	case ConnectionDialing:
		status = "dialing"
	case ConnectionEstablished:
		status = "established"
	case ConnectionClosed:
		status = "closed"
	case ConnectionError:
		status = "errored"
	}

	switch c.status {
	case ConnectionEstablished:
		return fmt.Sprintf("Connection %d has become %s after %s", c.ID, status, c.metrics.tcpEstablishedDuration)
	default:
		return fmt.Sprintf("Connection %d is %s", c.ID, status)
	}

}

// GetTCPProcessingDuration returns the time spent processing the connection
func (c Connection) GetTCPProcessingDuration() time.Duration {
	if WentOk(c) {
		return c.metrics.tcpEstablishedDuration
	}
	return c.metrics.tcpErroredDuration
}

func (c Connection) isStatusIn(statuses []ConnectionStatus) bool {
	for _, s := range statuses {
		if c.GetConnectionStatus() == s {
			return true
		}
	}
	return false
}

// WentOk return true when the Connection is Established or Closed state
func WentOk(c Connection) bool {
	return c.isStatusIn([]ConnectionStatus{ConnectionEstablished, ConnectionClosed})
}

// IsOk return true when the Connection is Established
func IsOk(c Connection) bool {
	return c.isStatusIn([]ConnectionStatus{ConnectionEstablished})
}

// WithError return true when the Connection is in Error state
func WithError(c Connection) bool {
	return c.isStatusIn([]ConnectionStatus{ConnectionError})
}

//PendingToProcess return true when the Connection is Established or Closed state
func PendingToProcess(c Connection) bool {
	return c.isStatusIn([]ConnectionStatus{ConnectionNotInitiated, ConnectionDialing})
}
