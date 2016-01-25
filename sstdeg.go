/*
 * Copyright (C) 2015 Fabrício Godoy <skarllot@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA  02111-1307, USA.
 */

package crypt

import (
	"io"
	"sync"
	"time"
)

const (
	// Defines default sleep time to ensure unpredictability.
	defaultSleepTime = time.Nanosecond * 200

	// defaultPoolSize defines the size of entropy pool for SSTDEG.
	defaultPoolSize = 4096
)

// A SSTDEG (System Sleep Time Delta Entropy Gathering) provides a pseudo-random
// generator based on unpredictable syscall time deltas of sleep calls.
type SSTDEG struct {
	pool  [defaultPoolSize]byte
	size  int
	mutex *sync.Mutex
	stop  chan bool
}

// NewSSTDEG creates a new instance of SSTDEG.
func NewSSTDEG() *SSTDEG {
	result := &SSTDEG{
		size:  0,
		mutex: &sync.Mutex{},
		stop:  make(chan bool, 0),
	}

	go result.generator()
	time.Sleep(defaultSleepTime)

	return result
}

// Close stops background routine that fills entropy pool.
func (s *SSTDEG) Close() error {
	if s.mutex == nil {
		return nil
	}

	s.stop <- true

	s.mutex = nil
	s.size = 0
	s.stop = nil
	return nil
}

// EntropyAvailable returns the entropy pool size of current instance.
func (s *SSTDEG) EntropyAvailable() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.size
}

// pop removes the n-elements from pool if available.
func (s *SSTDEG) pop(b []byte) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	l := len(b)

	if s.size >= l {
		copy(b, s.pool[s.size-l:s.size])
		s.size -= l
		return true
	}

	return false
}

// Read fills specified byte array with random data.
// Always return parameter array length and no errors.
func (s *SSTDEG) Read(b []byte) (n int, err error) {
	var chunk []byte
	if len(b) > defaultPoolSize {
		chunk = b[:defaultPoolSize]
	} else {
		chunk = b
	}

	ok := s.pop(chunk)

	for !ok {
		time.Sleep(defaultSleepTime)
		ok = s.pop(chunk)
	}

	return len(chunk), err
}

// generator fills entropy pool for this instance.
func (s *SSTDEG) generator() {
	var rndBits [2]byte
	var index byte
	var overflowCounter int

	for {
		rndDuration := time.Duration(getUInt16FromBytes(rndBits))
		before := time.Now()

		select {
		case <-time.After(defaultSleepTime + rndDuration):
			diff := time.Now().Sub(before)
			n := byte(diff.Nanoseconds() / 100)

			rndBits[index] = n
			index ^= 1

			s.mutex.Lock()
			if s.size < defaultPoolSize {
				s.pool[s.size] = n
				s.size++
			} else {
				if overflowCounter == defaultPoolSize {
					overflowCounter = 0
				}

				s.pool[overflowCounter] ^= n
				overflowCounter++
			}
			s.mutex.Unlock()
		case <-s.stop:
			return
		}
	}
}

// getUInt16FromBytes convert a 2-byte array to 16-bit unsigned integer.
func getUInt16FromBytes(input [2]byte) uint16 {
	return uint16(input[0]) + uint16(input[1])*256
}

var _ io.ReadCloser = (*SSTDEG)(nil)
