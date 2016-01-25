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

import "io"

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

// Close iterate over io.Closer sources to close them.
func (s *RandomAggr) Close() error {
	var err error
	for _, v := range s.sources {
		if closer, ok := v.Reader.(io.Closer); ok {
			itemErr := closer.Close()
			if err == nil {
				err = itemErr
			}
		}
	}

	return err
}

// Read fills specified byte array with random data from all sources.
func (s *RandomAggr) Read(b []byte) (n int, err error) {
	remainder := len(b)
	pos := 0
	sumWeight := s.sumWeight

	for _, v := range s.sources {
		count := int(float32(remainder) * (float32(v.Weight) / float32(sumWeight)))
		n, err = io.ReadFull(v.Reader, b[pos:pos+count])
		if err != nil && err != io.ErrUnexpectedEOF {
			return
		}

		pos += n
		remainder -= n
		sumWeight -= v.Weight
	}

	return len(b) - remainder, nil
}

var _ io.ReadCloser = (*RandomAggr)(nil)
