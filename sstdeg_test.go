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
	"crypto/rand"
	"io"
	"testing"
	"time"

	"github.com/GaryBoone/GoStats/stats"
)

const (
	TestingRounds            = 1000
	PredictabilityThreshold  = .05
	MinimumStandardDeviation = 65.0
)

func testUnpred(r io.Reader) (int, float64) {
	dict := make(map[int16]bool)
	buff := make([]byte, 2)
	count := 0
	var prob stats.Stats

	for i := 0; i < TestingRounds; i++ {
		r.Read(buff)
		val := int16(buff[0]) + int16(buff[1])*256

		if _, ok := dict[val]; ok {
			count++
		} else {
			dict[val] = true
		}
		prob.Update(float64(buff[0]))
		prob.Update(float64(buff[1]))
	}

	return count, prob.PopulationStandardDeviation()
}

func TestSystemUnpredictability(t *testing.T) {
	count, stddev := testUnpred(rand.Reader)

	if count > TestingRounds*PredictabilityThreshold {
		t.Errorf(
			"System random generator: %d dups of %d",
			count, TestingRounds)
	}
	if stddev < MinimumStandardDeviation {
		t.Errorf(
			"System random generator: %.2f STDDEV (%.2f minimum)",
			stddev, MinimumStandardDeviation)
	}
	t.Logf(
		"System random generator: %.2f%% dups/%.2f STDDEV",
		(float32(count)/float32(TestingRounds))*100, stddev)
}

func TestSSTDEGUnpredictability(t *testing.T) {
	rnd := NewSSTDEG()
	defer rnd.Close()

	count, stddev := testUnpred(rnd)

	if count > TestingRounds*PredictabilityThreshold {
		t.Errorf(
			"SSTDEG random generator: %d dups of %d",
			count, TestingRounds)
	}
	if stddev < MinimumStandardDeviation {
		t.Errorf(
			"SSTDEG random generator: %.2f STDDEV (%.2f minimum)",
			stddev, MinimumStandardDeviation)
	}
	t.Logf(
		"SSTDEG random generator: %.2f%% dups/%.2f STDDEV",
		(float32(count)/float32(TestingRounds))*100, stddev)
}

func TestSSTDEGFillEntropyBuffer(t *testing.T) {
	rnd := NewSSTDEG()
	defer rnd.Close()

	for rnd.EntropyAvailable() < SSTDEGPoolSize {
		time.Sleep(defaultSleepTime)
	}
}

func BenchmarkSSTDEG(b *testing.B) {
	rnd := NewSSTDEG()
	buff := make([]byte, b.N)
	b.ResetTimer()

	n, err := rnd.Read(buff)
	idx := n - 1
	for err == io.EOF {
		n, err = rnd.Read(buff[idx:])
		idx += n - 1
	}

	b.StopTimer()
	rnd.Close()
}

func BenchmarkSSTDEGWait(b *testing.B) {
	rnd := NewSSTDEG()
	buff := make([]byte, b.N)
	for rnd.EntropyAvailable() < SSTDEGPoolSize {
		// Waits for entropy buffer filling
		time.Sleep(time.Millisecond)
	}
	b.ResetTimer()

	n, err := rnd.Read(buff)
	idx := n - 1
	for err == io.EOF {
		n, err = rnd.Read(buff[idx:])
		idx += n - 1
	}

	b.StopTimer()
	rnd.Close()
}
