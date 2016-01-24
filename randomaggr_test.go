/*
 * Copyright (C) 2016 Fabrício Godoy <skarllot@gmail.com>
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
	"testing"
)

type InfiniteSource int

func (s InfiniteSource) Read(b []byte) (int, error) {
	for i := range b {
		b[i] = byte(s)
	}

	return len(b), nil
}

type LimitedSource struct {
	val  int
	size int
}

func (s *LimitedSource) Read(b []byte) (int, error) {
	counter := 0

	for i := range b {
		if s.size <= 0 {
			return counter, nil
		}

		b[i] = byte(s.val)
		s.size--
		counter++
	}

	return counter, nil
}

func countByValue(b []byte) map[byte]int {
	result := make(map[byte]int, 0)
	for _, v := range b {
		if _, ok := result[v]; ok {
			result[v]++
		} else {
			result[v] = 1
		}
	}

	return result
}

func TestRandomAggrDistribution(t *testing.T) {
	rnd := NewRandomAggr().
		Add(&LimitedSource{1, 10}, 2).
		Add(InfiniteSource(2), 3).
		Add(InfiniteSource(3), 5).
		Build()

	buf := make([]byte, 100)
	n, err := rnd.Read(buf)
	if err != nil {
		t.Errorf("Error reading from aggregation: %v", err)
	}
	if n != len(buf) {
		t.Errorf("Should fill entire buffer: read %d bytes", n)
	}

	testValues(buf[:10], 1, t)
	testValues(buf[10:43], 2, t)
	testValues(buf[43:], 3, t)
}

func TestRandomAggrDistribution2(t *testing.T) {
	rnd := NewRandomAggr().
		Add(&LimitedSource{1, 10}, 2).
		Add(&LimitedSource{2, 20}, 3).
		Add(InfiniteSource(3), 5).
		Build()

	buf := make([]byte, 100)
	n, err := rnd.Read(buf)
	if err != nil {
		t.Errorf("Error reading from aggregation: %v", err)
	}
	if n != len(buf) {
		t.Errorf("Should fill entire buffer: read %d bytes", n)
	}

	testValues(buf[:10], 1, t)
	testValues(buf[10:30], 2, t)
	testValues(buf[30:], 3, t)
}

func TestRandomAggrInsufficient(t *testing.T) {
	rnd := NewRandomAggr().
		Add(&LimitedSource{1, 10}, 2).
		Add(&LimitedSource{2, 20}, 3).
		Add(&LimitedSource{3, 30}, 5).
		Build()

	buf := make([]byte, 100)
	n, err := rnd.Read(buf)
	if err != nil {
		t.Errorf("Error reading from aggregation: %v", err)
	}
	if n != 60 {
		t.Errorf("Should read 60 bytes: read %d bytes", n)
	}

	testValues(buf[:10], 1, t)
	testValues(buf[10:30], 2, t)
	testValues(buf[30:60], 3, t)
	testValues(buf[60:], 0, t)
}

func TestRandomAggrWeight(t *testing.T) {
	rnd := NewRandomAggr().
		Add(InfiniteSource(1), 10).
		Add(InfiniteSource(2), 8).
		Add(InfiniteSource(3), 3).
		Build()

	buf := make([]byte, 10+8+3)
	n, err := rnd.Read(buf)
	if err != nil {
		t.Errorf("Error reading from aggregation: %v", err)
	}
	if n != len(buf) {
		t.Errorf("Should fill entire buffer: read %d bytes", n)
	}

	testValues(buf[:10], 1, t)
	testValues(buf[10:18], 2, t)
	testValues(buf[18:], 3, t)
}

func TestRandomAggrRepeatingDecimal(t *testing.T) {
	rnd := NewRandomAggr().
		Add(InfiniteSource(1), 5).
		Add(InfiniteSource(2), 4).
		Build()

	buf := make([]byte, 9)
	n, err := rnd.Read(buf)
	if err != nil {
		t.Errorf("Error reading from aggregation: %v", err)
	}
	if n != len(buf) {
		t.Errorf("Should fill entire buffer: read %d bytes", n)
	}

	testValues(buf[:5], 1, t)
	testValues(buf[5:], 2, t)
}

func TestRandomAggrRepeatingDecimal2(t *testing.T) {
	rnd := NewRandomAggr().
		Add(InfiniteSource(1), 1).
		Add(InfiniteSource(2), 9).
		Build()

	buf := make([]byte, 10)
	n, err := rnd.Read(buf)
	if err != nil {
		t.Errorf("Error reading from aggregation: %v", err)
	}
	if n != len(buf) {
		t.Errorf("Should fill entire buffer: read %d bytes", n)
	}

	testValues(buf[:1], 1, t)
	testValues(buf[1:], 2, t)
}

func testValues(b []byte, expected byte, t *testing.T) {
	for _, v := range b {
		if v != expected {
			t.Errorf("Unexpected value: got %d instead of %d", v, expected)
		}
	}
}
