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
	"fmt"
	"io"
)

// A source defines a source of random data and its weight from total.
type source struct {
	// The reader of random data.
	Reader io.Reader
	// The weight of current random source.
	Weight int
}

// A RandomAggr represents an aggregation of random data sources.
type RandomAggr struct {
	sources   []source
	sumWeight int
}

// Read fills specified byte array with random data from all sources.
func (s *RandomAggr) Read(b []byte) (n int, err error) {
	l := len(b)
	pos := 0

	for _, v := range s.sources {
		count := int(float32(l) * (float32(v.Weight) / float32(s.sumWeight)))
		n, err = v.Reader.Read(b[pos : pos+count])
		if err != nil || n != count {
			fmt.Printf("Error(%d): %v\n%#v\n\n", n, err, v)
			return
		}
		pos += count
	}

	return l, nil
}

var _ io.Reader = (*RandomAggr)(nil)
