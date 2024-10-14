package model

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

func NewQueueJob(queue string, payload any, opts ...QueueJobOption) (QueueJob, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return QueueJob{}, err
	}
	payloadStr := string(payloadBytes)
	h := sha256.New()
	_, err = io.WriteString(h, queue+payloadStr)
	if err != nil {
		return QueueJob{}, err
	}
	j := QueueJob{
		Fingerprint:      fmt.Sprintf("%x", h.Sum(nil)),
		Queue:            queue,
		Status:           QueueJobStatusPending,
		Payload:          payloadStr,
		RunAfter:         time.Now(),
		ArchivalDuration: Duration(7 * 24 * time.Hour),
	}
	for _, opt := range opts {
		opt(&j)
	}
	return j, nil
}

type QueueJobOption func(*QueueJob)

func QueueJobMaxRetries(maxRetries uint) QueueJobOption {
	return func(j *QueueJob) {
		j.MaxRetries = maxRetries
	}
}

func QueueJobPriority(priority int) QueueJobOption {
	return func(j *QueueJob) {
		j.Priority = priority
	}
}

func QueueJobDelayBy(duration time.Duration) QueueJobOption {
	return func(j *QueueJob) {
		j.RunAfter = j.RunAfter.Add(duration)
	}
}
