package client

import (
	"context"

	"k8s.io/klog/v2"
)

// Progress contains progress information about this operation.
type Progress struct {
	Finished      bool
	Image         string
	Tag           string
	BytesReceived uint64
	Error         error
}

// ContainerInfo contains information about a container on the target system.
type ContainerInfo struct {
	ID        string
	Name      string
	ImageName string
	State     string

	Error error
}

// LogMessage contains the log message retrieved from the target system as well as any error that
// may have occurred
type LogMessage struct {
	Msg string

	Error error
}

type startOptions struct {
	envs  []string
	ports []string
}

type nonBlockTypes interface {
	*Progress | *ContainerInfo | *LogMessage
}

// StartOption is an option passed to a start container call.
type StartOption func(*startOptions)

// WithEnv sets the environment to be passed to the start operation.
func WithEnv(envs []string) StartOption {
	return func(opt *startOptions) {
		opt.envs = envs
	}
}

// WithPorts sets the ports to be passed to the start operation.
func WithPorts(ports []string) StartOption {
	return func(opt *startOptions) {
		opt.ports = ports
	}
}

// nonBlockingChannelSend attempts to send a message in a non blocking manner. If the context is
// cancelled it simply returns with an indication that the context was cancelled
func nonBlockingChannelSend[T nonBlockTypes](ctx context.Context, ch chan T, data T) bool {
	select {
	case <-ctx.Done():
		return true
	case ch <- data:
		return false
	default:
		klog.Warningf("unable to send message; dropping")
	}
	return false
}