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
	"testing"
)

func TestSaltUnpredictability(t *testing.T) {
	dict := make(map[string]bool)
	s := NewSalter(NewRandomAggr().SecureSet(), nil)
	defer s.Dispose()
	count := 0

	for i := 0; i < TestingRounds; i++ {
		val, err := s.Token(0)
		if err != nil {
			t.Fatalf("Error creating a new token: %v", err)
		}

		if _, ok := dict[val]; ok {
			count++
		} else {
			dict[val] = true
		}
	}

	if count > 0 {
		t.Errorf(
			"Salter class could not generate unpredictable data: %d of %d",
			count, TestingRounds)
	}
}

func BenchmarkSalter(b *testing.B) {
	salter := NewSalter(rand.Reader, nil)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		salter.Token(0)
	}

	b.StopTimer()
	salter.Dispose()
}

func BenchmarkSalterSecure(b *testing.B) {
	salter := NewSalter(NewRandomAggr().SecureSet(), nil)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		salter.Token(0)
	}

	b.StopTimer()
	salter.Dispose()
}
