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
	err  error
	read int
	msg  string
}

func newReadError(msg string, n int, err error) readError {
	return readError{
		err,
		n,
		msg,
	}
}

func (e readError) Error() string {
	return fmt.Sprintf("%s (read '%d' bytes)\nInner error: %v",
		e.msg, e.read, e.err)
}

func (e readError) String() string {
	return e.Error()
}
