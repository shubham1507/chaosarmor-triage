package main

import (
	"crypto/sha256"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type timelineEvent struct {
	timestamp time.Time
	message   string
}

func readTimeline(input io.Reader) ([]timelineEvent, error) {
	reader := csv.NewReader(input)
	reader.FieldsPerRecord = 2
	var events []timelineEvent

	for record := 1; ; record++ {
		fields, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("record %d: %w", record, err)
		}

		timestamp, err := time.Parse(time.RFC3339, strings.TrimSpace(fields[0]))
		if err != nil {
			return nil, fmt.Errorf("record %d timestamp: %w", record, err)
		}

		message := strings.TrimSpace(fields[1])
		if message == "" {
			return nil, fmt.Errorf("record %d: event message is empty", record)
		}
		events = append(events, timelineEvent{timestamp: timestamp, message: message})
	}

	sort.SliceStable(events, func(i, j int) bool {
		return events[i].timestamp.Before(events[j].timestamp)
	})
	return events, nil
}

func printTimelineFile(output io.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open timeline: %w", err)
	}
	defer file.Close()

	digest := sha256.New()
	events, err := readTimeline(io.TeeReader(file, digest))
	if err != nil {
		return fmt.Errorf("read timeline %q: %w", path, err)
	}

	for _, event := range events {
		if _, err := fmt.Fprintf(output, "%s  %q\n",
			event.timestamp.UTC().Format(time.RFC3339Nano), event.message); err != nil {
			return fmt.Errorf("write timeline: %w", err)
		}
	}
	if _, err := fmt.Fprintf(output, "SHA-256: %x\n", digest.Sum(nil)); err != nil {
		return fmt.Errorf("write timeline digest: %w", err)
	}
	return nil
}

func reportTimelineError(err error) {
	logger := log.New(os.Stderr, "chaosarmor-triage: ", 0)
	logger.Print(err)
}
