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
	"fmt"
)

type readError struct {
	expected int
	got      int
}

func newReadError(expected, got int) readError {
	return readError{
		expected,
		got,
	}
}

func (e readError) Error() string {
	return fmt.Sprintf("Error reading: got %d bytes instead of %d",
		e.got, e.expected)
}

func (e readError) String() string {
	return e.Error()
}
