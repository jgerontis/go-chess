package uci

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os/exec"
	"strings"
	"time"
)

// Connection represents a UCI connection (TCP or process)
type Connection interface {
	io.ReadWriteCloser
}

// ProcessConnection wraps exec.Cmd to implement Connection interface
type ProcessConnection struct {
	stdin  io.WriteCloser
	stdout io.ReadCloser
	cmd    *exec.Cmd
}

func (p *ProcessConnection) Read(b []byte) (int, error) {
	return p.stdout.Read(b)
}

func (p *ProcessConnection) Write(b []byte) (int, error) {
	return p.stdin.Write(b)
}

func (p *ProcessConnection) Close() error {
	// Close stdin and stdout first
	if p.stdin != nil {
		p.stdin.Close()
	}
	if p.stdout != nil {
		p.stdout.Close()
	}
	
	// Wait for process to finish or kill it
	if p.cmd != nil && p.cmd.Process != nil {
		// Give process some time to exit gracefully
		done := make(chan error, 1)
		go func() {
			done <- p.cmd.Wait()
		}()
		
		select {
		case <-done:
			// Process exited
		case <-time.After(2 * time.Second):
			// Force kill if it doesn't exit
			p.cmd.Process.Kill()
		}
	}
	return nil
}

// Client represents a UCI client that can communicate with engines
type Client struct {
	conn   Connection
	reader *bufio.Reader
	writer *bufio.Writer
	cmd    *exec.Cmd // Only used for exec mode
}

// NewClientTCP creates a UCI client connected to a TCP server
func NewClientTCP(address string) (*Client, error) {
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to UCI server at %s: %v", address, err)
	}

	return &Client{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}, nil
}

// NewClientExec creates a UCI client that launches and communicates with an engine process
func NewClientExec(enginePath string, args ...string) (*Client, error) {
	cmd := exec.Command(enginePath, args...)
	
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %v", err)
	}
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		return nil, fmt.Errorf("failed to create stdout pipe: %v", err)
	}
	
	// Start the process
	if err := cmd.Start(); err != nil {
		stdin.Close()
		stdout.Close()
		return nil, fmt.Errorf("failed to start engine process: %v", err)
	}
	
	processConn := &ProcessConnection{
		stdin:  stdin,
		stdout: stdout,
		cmd:    cmd,
	}

	return &Client{
		conn:   processConn,
		reader: bufio.NewReader(processConn),
		writer: bufio.NewWriter(processConn),
		cmd:    cmd,
	}, nil
}

// Close closes the UCI connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// SendCommand sends a command to the UCI engine
func (c *Client) SendCommand(command string) error {
	_, err := c.writer.WriteString(command + "\n")
	if err != nil {
		return err
	}
	return c.writer.Flush()
}

// ReadResponse reads a single response line from the UCI engine
func (c *Client) ReadResponse() (string, error) {
	response, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(response), nil
}

// StartResponseListener starts a goroutine that continuously reads responses
// and sends them to the provided channel. This is useful for interactive applications.
func (c *Client) StartResponseListener(responseChan chan<- string, errorChan chan<- error) {
	go func() {
		defer close(responseChan)
		defer close(errorChan)
		
		for {
			response, err := c.ReadResponse()
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					return // Connection closed gracefully
				}
				errorChan <- err
				return
			}
			responseChan <- response
		}
	}()
}