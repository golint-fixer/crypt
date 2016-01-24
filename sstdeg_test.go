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
)

const (
	DefaultUnpredRounds            = 1000
	DefaultPredictabilityThreshold = .05
)

func testUnpred(r io.Reader) int {
	dict := make(map[int16]bool)
	buff := make([]byte, 2)
	count := 0

	for i := 0; i < DefaultUnpredRounds; i++ {
		r.Read(buff)
		val := int16(buff[0]) + int16(buff[1])*256

		if _, ok := dict[val]; ok {
			count++
		} else {
			dict[val] = true
		}
	}

	return count
}

func TestSystemUnpredictability(t *testing.T) {
	count := testUnpred(rand.Reader)

	if count > DefaultUnpredRounds*DefaultPredictabilityThreshold {
		t.Errorf(
			"System random generator could not generate unpredictable data: %d of %d",
			count, DefaultUnpredRounds)
	}
	t.Logf(
		"System random generator predictability: %.2f%%",
		(float32(count)/float32(DefaultUnpredRounds))*100)
}

func TestSSTDEGUnpredictability(t *testing.T) {
	rnd := NewSSTDEG()
	defer rnd.Dispose()

	count := testUnpred(rnd)

	if count > DefaultUnpredRounds*DefaultPredictabilityThreshold {
		t.Errorf(
			"SSTDEG random generator could not generate unpredictable data: %d of %d",
			count, DefaultUnpredRounds)
	}
	t.Logf(
		"SSTDEG random generator predictability: %.2f%%",
		(float32(count)/float32(DefaultUnpredRounds))*100)
}

func TestSSTDEGFillEntropyBuffer(t *testing.T) {
	rnd := NewSSTDEG()
	defer rnd.Dispose()

	for rnd.EntropyAvailable() < SSTDEGPoolSize {
		time.Sleep(defaultSleepTime)
	}
}

func BenchmarkSSTDEG(b *testing.B) {
	rnd := NewSSTDEG()
	buff := make([]byte, 1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rnd.Read(buff)
	}

	b.StopTimer()
	rnd.Dispose()
}

func BenchmarkSSTDEGBatch(b *testing.B) {
	rnd := NewSSTDEG()
	buff := make([]byte, DefaultTokenSize)
	for rnd.EntropyAvailable() < SSTDEGPoolSize {
		// Waits for entropy buffer filling
		time.Sleep(time.Millisecond)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rnd.Read(buff)
	}

	b.StopTimer()
	rnd.Dispose()
}
