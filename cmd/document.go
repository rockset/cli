package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rockset/rockset-go-client/openapi"
	"io"
	"log"
)

type DocumentAdder interface {
	AddDocuments(ctx context.Context, workspace, collection string,
		docs []interface{}) ([]openapi.DocumentStatus, error)
}

type StreamConfig struct {
	Workspace  string
	Collection string
	BatchSize  uint64
}

type Streamer struct {
	adder DocumentAdder
	StreamConfig
}

func NewStreamer(s DocumentAdder, cfg StreamConfig) Streamer {
	return Streamer{
		adder:        s,
		StreamConfig: cfg,
	}
}

func (s *Streamer) Stream(ctx context.Context, in io.Reader) (uint64, error) {
	var counter uint64
	buf := make([]interface{}, 0, s.BatchSize)
	d := json.NewDecoder(in)

	for {
		// this is a generic way to describe a json object
		var j map[string]interface{}

		if err := d.Decode(&j); err != nil {
			if err == io.EOF {
				// flush remaining
				cnt, err := s.flush(ctx, buf)
				if err != nil {
					return counter, err
				}
				counter += cnt

				return counter, nil
			}
			// should this just log the error so we skip incorrect json?
			return counter, err
		}

		buf = append(buf, j)
		if uint64(len(buf)) < s.BatchSize {
			continue
		}

		cnt, err := s.flush(ctx, buf)
		if err != nil {
			return counter, err
		}
		counter += cnt

		buf = make([]interface{}, 0, s.BatchSize)
	}
}

// flush the buffered docs and return how many were added
func (s *Streamer) flush(ctx context.Context, buf []interface{}) (uint64, error) {
	res, err := s.adder.AddDocuments(ctx, s.Workspace, s.Collection, buf)
	if err != nil {
		return 0, fmt.Errorf("failed to flush %d documents: %w", len(buf), err)
	}

	var count uint64
	for i, r := range res {
		if r.GetStatus() == "ADDED" {
			count++
			continue
		}
		log.Printf("result %d: %s", i, r.GetStatus())
	}

	return count, nil
}
