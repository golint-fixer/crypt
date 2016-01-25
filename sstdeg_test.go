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
	MaximumDups              = .05
	MinimumStandardDeviation = 50.0
)

func testUnpred(r io.Reader) (float64, float64) {
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

	return float64(count) / float64(TestingRounds),
		prob.PopulationStandardDeviation()
}

func TestSystemUnpredictability(t *testing.T) {
	dups, stddev := testUnpred(rand.Reader)

	if dups > MaximumDups {
		t.Errorf(
			"System random generator: %d dups of %d",
			int(TestingRounds*dups), TestingRounds)
	}
	if stddev < MinimumStandardDeviation {
		t.Errorf(
			"System random generator: %.2f STDDEV (%.2f minimum)",
			stddev, MinimumStandardDeviation)
	}
	t.Logf(
		"System random generator: %.2f%% dups/%.2f STDDEV",
		dups*100, stddev)
}

func TestSSTDEGUnpredictability(t *testing.T) {
	rnd := NewSSTDEG()
	defer rnd.Close()

	dups, stddev := testUnpred(rnd)

	if dups > MaximumDups {
		t.Errorf(
			"SSTDEG random generator: %d dups of %d",
			int(TestingRounds*dups), TestingRounds)
	}
	if stddev < MinimumStandardDeviation {
		t.Errorf(
			"SSTDEG random generator: %.2f STDDEV (%.2f minimum)",
			stddev, MinimumStandardDeviation)
	}
	t.Logf(
		"SSTDEG random generator: %.2f%% dups/%.2f STDDEV",
		dups*100, stddev)
}

func TestSSTDEGFillEntropyBuffer(t *testing.T) {
	rnd := NewSSTDEG()
	defer rnd.Close()

	for rnd.EntropyAvailable() < defaultPoolSize {
		time.Sleep(defaultSleepTime)
	}
}

func BenchmarkSSTDEG(b *testing.B) {
	rnd := NewSSTDEG()
	buff := make([]byte, b.N)
	b.ResetTimer()

	n, err := io.ReadFull(rnd, buff)
	if err != nil {
		b.Fatalf("Error reading SSTDEG: %v", err)
	} else if n < len(buff) {
		b.Fatalf("Error reading SSTDEG: should read %d bytes but read %d",
			len(buff), n)
	}

	b.StopTimer()
	rnd.Close()
}

func BenchmarkSSTDEGWait(b *testing.B) {
	rnd := NewSSTDEG()
	buff := make([]byte, b.N)
	for rnd.EntropyAvailable() < defaultPoolSize {
		// Waits for entropy buffer filling
		time.Sleep(time.Millisecond)
	}
	b.ResetTimer()

	n, err := io.ReadFull(rnd, buff)
	if err != nil {
		b.Fatalf("Error reading SSTDEG: %v", err)
	} else if n < len(buff) {
		b.Fatalf("Error reading SSTDEG: should read %d bytes but read %d",
			len(buff), n)
	}

	b.StopTimer()
	rnd.Close()
}
