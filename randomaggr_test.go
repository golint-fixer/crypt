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

type FooSource int

func (f FooSource) Read(b []byte) (int, error) {
	for i := range b {
		b[i] = byte(f)
	}

	return len(b), nil
}

func TestRandomAggrWeight(t *testing.T) {
	rnd := NewRandomAggr().
		Add(FooSource(1), 10).
		Add(FooSource(2), 8).
		Add(FooSource(3), 3).
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
		Add(FooSource(1), 5).
		Add(FooSource(2), 4).
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
		Add(FooSource(1), 1).
		Add(FooSource(2), 9).
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
