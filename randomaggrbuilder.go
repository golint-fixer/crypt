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
	"crypto/rand"
	"io"
)

// A RandomAggrBuilder provides methods to build a new RandomAggr.
type RandomAggrBuilder interface {
	// Add a custom random source and specify a weight.
	Add(io.Reader, int) RandomAggrBuilder

	// AddSSTDEG adds a SSTDEG pseudo-random generator and specifies a weight.
	AddSSTDEG(int) RandomAggrBuilder

	// AddSys adds a system pseudo-random generator and specifies a weight.
	AddSys(int) RandomAggrBuilder

	// Build creates and returns a new RandomAggr.
	Build() *RandomAggr

	// FastSet get a RandomSource backed 100% by system pseudo-random generator.
	FastSet() *RandomAggr

	// InsecureSet get a RandomSource backed 100% by SSTDEG pseudo-random
	// generator.
	InsecureSet() *RandomAggr

	// SecureSet get a RandomSource backed 84% by system pseudo-random generator
	// and 16% by SSTDEG pseudo-random generator.
	SecureSet() *RandomAggr
}

type rndaggb struct {
	sources []source
}

// NewRandomAggr creates a new instance of RandomAggrBuilder.
func NewRandomAggr() RandomAggrBuilder {
	return &rndaggb{
		make([]source, 0),
	}
}

func (b *rndaggb) Add(r io.Reader, w int) RandomAggrBuilder {
	b.sources = append(b.sources, source{r, w})
	return b
}

func (b *rndaggb) AddSSTDEG(w int) RandomAggrBuilder {
	b.sources = append(b.sources, source{NewSSTDEG(), w})
	return b
}

func (b *rndaggb) AddSys(w int) RandomAggrBuilder {
	b.sources = append(b.sources, source{rand.Reader, w})
	return b
}

func (b *rndaggb) Build() *RandomAggr {
	sum := 0
	for _, v := range b.sources {
		sum += v.Weight
	}

	return &RandomAggr{b.sources, sum}
}

func (b *rndaggb) FastSet() *RandomAggr {
	return &RandomAggr{
		[]source{
			{
				Reader: rand.Reader,
				Weight: 1,
			},
		},
		1,
	}
}

func (b *rndaggb) InsecureSet() *RandomAggr {
	return &RandomAggr{
		[]source{
			{
				Reader: NewSSTDEG(),
				Weight: 1,
			},
		},
		1,
	}
}

func (b *rndaggb) SecureSet() *RandomAggr {
	sstdeg := NewSSTDEG()
	return &RandomAggr{
		[]source{
			{
				Reader: rand.Reader,
				Weight: 8,
			},
			{
				Reader: sstdeg,
				Weight: 1,
			},
			{
				Reader: rand.Reader,
				Weight: 4,
			},
			{
				Reader: sstdeg,
				Weight: 1,
			},
			{
				Reader: rand.Reader,
				Weight: 4,
			},
			{
				Reader: sstdeg,
				Weight: 1,
			},
		},
		8 + 1 + 4 + 1 + 4 + 1,
	}
}
