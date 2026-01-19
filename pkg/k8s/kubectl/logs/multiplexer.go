package logs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

// StreamSource represents a single log stream source
type StreamSource struct {
	ComponentID   string
	ComponentName string
	Namespace     string
	PodName       string
	Container     string
	Streamer      *LogStreamer
}

// Multiplexer coordinates concurrent log streaming from multiple sources
type Multiplexer struct {
	Sources      []*StreamSource
	Prefix       bool
	NoColor      bool
	OutputFormat string // "stylish", "json", "yaml"
	errChan      chan error
	doneChan     chan struct{}
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
}

// PrefixWriter wraps an io.Writer with thread-safe line prefixing
type PrefixWriter struct {
	Writer     io.Writer
	Prefix     string
	Color      *color.Color
	mu         sync.Mutex
	lineBuffer bytes.Buffer
}

// NewMultiplexer creates a new log multiplexer
func NewMultiplexer(sources []*StreamSource, prefix bool, noColor bool, outputFormat string) *Multiplexer {
	ctx, cancel := context.WithCancel(context.Background())

	return &Multiplexer{
		Sources:      sources,
		Prefix:       prefix,
		NoColor:      noColor,
		OutputFormat: outputFormat,
		errChan:      make(chan error, len(sources)),
		doneChan:     make(chan struct{}),
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Start begins streaming from all sources concurrently
func (m *Multiplexer) Start() error {
	// Start a goroutine for each source
	for _, source := range m.Sources {
		m.wg.Add(1)
		go m.streamSource(source)
	}

	// Wait for all streams to complete
	go func() {
		m.wg.Wait()
		close(m.doneChan)
	}()

	return nil
}

// Wait blocks until all streams complete or are stopped
func (m *Multiplexer) Wait() []error {
	<-m.doneChan

	// Collect any errors
	var errs []error
	close(m.errChan)
	for err := range m.errChan {
		errs = append(errs, err)
	}

	return errs
}

// Stop cancels all active streams
func (m *Multiplexer) Stop() {
	m.cancel()
}

// streamSource handles streaming from a single source
func (m *Multiplexer) streamSource(source *StreamSource) {
	defer m.wg.Done()

	// Start the log stream
	err := source.Streamer.Start(m.ctx)
	if err != nil {
		m.errChan <- fmt.Errorf("[%s] Failed to start log stream: %w",
			m.formatSourcePrefix(source), err)
		return
	}
	defer source.Streamer.Close()

	// Create appropriate writer based on output format
	var writer io.Writer = os.Stdout

	switch m.OutputFormat {
	case "json":
		writer = NewJSONWriter(os.Stdout, source)
	case "yaml":
		writer = NewYAMLWriter(os.Stdout, source)
	default: // "stylish" or any other format
		if m.Prefix {
			prefix := m.formatSourcePrefix(source)
			prefixColor := m.getColorForSource(source)
			writer = NewPrefixWriter(os.Stdout, prefix, prefixColor, m.NoColor)
		}
	}

	// Copy logs to output
	_, err = io.Copy(writer, source.Streamer)
	if err != nil && err != io.EOF {
		// Check if error is due to context cancellation
		if m.ctx.Err() == nil {
			m.errChan <- fmt.Errorf("[%s] Stream error: %w",
				m.formatSourcePrefix(source), err)
		}
	}
}

// formatSourcePrefix creates a prefix string for a source
func (m *Multiplexer) formatSourcePrefix(source *StreamSource) string {
	name := source.ComponentName
	if name == "" {
		name = source.ComponentID
	}

	if source.Container != "" {
		return fmt.Sprintf("%s/%s/%s", name, source.PodName, source.Container)
	}

	return fmt.Sprintf("%s/%s", name, source.PodName)
}

// getColorForSource returns a consistent color for a source based on hash
func (m *Multiplexer) getColorForSource(source *StreamSource) *color.Color {
	colors := []*color.Color{
		color.New(color.FgCyan),
		color.New(color.FgGreen),
		color.New(color.FgYellow),
		color.New(color.FgBlue),
		color.New(color.FgMagenta),
	}

	// Hash the source identifier to get consistent color
	h := fnv.New32a()
	h.Write([]byte(m.formatSourcePrefix(source)))
	index := h.Sum32() % uint32(len(colors))

	return colors[index]
}

// NewPrefixWriter creates a new prefix writer
func NewPrefixWriter(w io.Writer, prefix string, c *color.Color, noColor bool) *PrefixWriter {
	return &PrefixWriter{
		Writer: w,
		Prefix: prefix,
		Color:  c,
	}
}

// Write writes data with prefix for each line
func (pw *PrefixWriter) Write(p []byte) (int, error) {
	pw.mu.Lock()
	defer pw.mu.Unlock()

	totalWritten := 0

	// Process input byte by byte to handle line breaks
	for _, b := range p {
		pw.lineBuffer.WriteByte(b)

		// When we hit a newline, flush the line with prefix
		if b == '\n' {
			line := pw.lineBuffer.String()
			pw.lineBuffer.Reset()

			// Write prefix and line
			var written int
			var err error
			if pw.Color != nil {
				written, err = pw.Color.Fprintf(pw.Writer, "[%s] %s", pw.Prefix, line)
			} else {
				written, err = fmt.Fprintf(pw.Writer, "[%s] %s", pw.Prefix, line)
			}

			if err != nil {
				return totalWritten, err
			}

			totalWritten += written
		}
	}

	return len(p), nil
}

// Flush flushes any remaining data in the buffer
func (pw *PrefixWriter) Flush() error {
	pw.mu.Lock()
	defer pw.mu.Unlock()

	if pw.lineBuffer.Len() > 0 {
		line := pw.lineBuffer.String()
		pw.lineBuffer.Reset()

		var err error
		if pw.Color != nil {
			_, err = pw.Color.Fprintf(pw.Writer, "[%s] %s\n", pw.Prefix, line)
		} else {
			_, err = fmt.Fprintf(pw.Writer, "[%s] %s\n", pw.Prefix, line)
		}

		return err
	}

	return nil
}

// LogEntry represents a structured log entry for JSON/YAML output
type LogEntry struct {
	Timestamp     string `json:"timestamp" yaml:"timestamp"`
	Component     string `json:"component" yaml:"component"`
	ComponentID   string `json:"componentId" yaml:"componentId"`
	Pod           string `json:"pod" yaml:"pod"`
	Container     string `json:"container" yaml:"container"`
	Namespace     string `json:"namespace" yaml:"namespace"`
	Message       string `json:"message" yaml:"message"`
}

// JSONWriter formats log lines as JSON objects
type JSONWriter struct {
	Writer    io.Writer
	Source    *StreamSource
	mu        sync.Mutex
	lineBuffer bytes.Buffer
}

// NewJSONWriter creates a new JSON writer
func NewJSONWriter(w io.Writer, source *StreamSource) *JSONWriter {
	return &JSONWriter{
		Writer: w,
		Source: source,
	}
}

// Write writes data as JSON formatted log entries
func (jw *JSONWriter) Write(p []byte) (int, error) {
	jw.mu.Lock()
	defer jw.mu.Unlock()

	totalWritten := 0

	// Process input byte by byte to handle line breaks
	for _, b := range p {
		jw.lineBuffer.WriteByte(b)

		// When we hit a newline, format as JSON and flush
		if b == '\n' {
			line := strings.TrimSuffix(jw.lineBuffer.String(), "\n")
			jw.lineBuffer.Reset()

			if line == "" {
				continue
			}

			entry := LogEntry{
				Timestamp:   time.Now().UTC().Format(time.RFC3339),
				Component:   jw.Source.ComponentName,
				ComponentID: jw.Source.ComponentID,
				Pod:         jw.Source.PodName,
				Container:   jw.Source.Container,
				Namespace:   jw.Source.Namespace,
				Message:     line,
			}

			jsonData, err := json.Marshal(entry)
			if err != nil {
				return totalWritten, err
			}

			written, err := fmt.Fprintf(jw.Writer, "%s\n", jsonData)
			if err != nil {
				return totalWritten, err
			}

			totalWritten += written
		}
	}

	return len(p), nil
}

// YAMLWriter formats log lines as YAML documents
type YAMLWriter struct {
	Writer     io.Writer
	Source     *StreamSource
	mu         sync.Mutex
	lineBuffer bytes.Buffer
}

// NewYAMLWriter creates a new YAML writer
func NewYAMLWriter(w io.Writer, source *StreamSource) *YAMLWriter {
	return &YAMLWriter{
		Writer: w,
		Source: source,
	}
}

// Write writes data as YAML formatted log entries
func (yw *YAMLWriter) Write(p []byte) (int, error) {
	yw.mu.Lock()
	defer yw.mu.Unlock()

	totalWritten := 0

	// Process input byte by byte to handle line breaks
	for _, b := range p {
		yw.lineBuffer.WriteByte(b)

		// When we hit a newline, format as YAML and flush
		if b == '\n' {
			line := strings.TrimSuffix(yw.lineBuffer.String(), "\n")
			yw.lineBuffer.Reset()

			if line == "" {
				continue
			}

			entry := LogEntry{
				Timestamp:   time.Now().UTC().Format(time.RFC3339),
				Component:   yw.Source.ComponentName,
				ComponentID: yw.Source.ComponentID,
				Pod:         yw.Source.PodName,
				Container:   yw.Source.Container,
				Namespace:   yw.Source.Namespace,
				Message:     line,
			}

			yamlData, err := yaml.Marshal(entry)
			if err != nil {
				return totalWritten, err
			}

			written, err := fmt.Fprintf(yw.Writer, "---\n%s", yamlData)
			if err != nil {
				return totalWritten, err
			}

			totalWritten += written
		}
	}

	return len(p), nil
}
