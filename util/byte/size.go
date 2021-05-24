// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gbyte

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ByteSize uint64

const (
	B  ByteSize = 1
	KB          = B << 10
	MB          = KB << 10
	GB          = MB << 10
	TB          = GB << 10
	PB          = TB << 10
	EB          = PB << 10

	fnUnmarshalText string = "UnmarshalText"
	maxUint64       uint64 = (1 << 64) - 1
)

var ErrBits = errors.New("unit with capital unit prefix and lower case unit (b) - bits, not bytes ")

var defaultDatasize ByteSize

func StrToSize(s string) (bytes uint64, err error) {
	err = defaultDatasize.UnmarshalText([]byte(s))
	if err != nil {
		return
	}
	bytes = defaultDatasize.Bytes()
	return
}
func SizeToStr(bytes uint64) (s string, err error) {
	err = defaultDatasize.UnmarshalText([]byte(fmt.Sprintf("%d", bytes)))
	if err != nil {
		return
	}
	s = defaultDatasize.HumanReadable()
	return
}
func (b ByteSize) Bytes() uint64 {
	return uint64(b)
}

func (b ByteSize) KBytes() float64 {
	if b < KB {
		return 0
	}
	v := b / KB
	r := b % KB
	return float64(v) + float64(r)/float64(KB)
}

func (b ByteSize) MBytes() float64 {
	if b < MB {
		return 0
	}
	v := b / MB
	r := b % MB
	return float64(v) + float64(r)/float64(MB)
}

func (b ByteSize) GBytes() float64 {
	if b < GB {
		return 0
	}
	v := b / GB
	r := b % GB
	return float64(v) + float64(r)/float64(GB)
}

func (b ByteSize) TBytes() float64 {
	if b < TB {
		return 0
	}
	v := b / TB
	r := b % TB
	return float64(v) + float64(r)/float64(TB)
}

func (b ByteSize) PBytes() float64 {
	if b < PB {
		return 0
	}
	v := b / PB
	r := b % PB
	return float64(v) + float64(r)/float64(PB)
}

func (b ByteSize) EBytes() float64 {
	if b < EB {
		return 0
	}
	v := b / EB
	r := b % EB
	return float64(v) + float64(r)/float64(EB)
}

func (b ByteSize) HumanReadable() (s string) {
	uint := "B"
	defer func() {
		if strings.Contains(s, ".") {
			s = strings.TrimRight(s, "0")
			s = strings.TrimSuffix(s, ".")
		}
		s += " " + uint
	}()
	switch {
	case b == 0:
		return fmt.Sprint("0")
	case b/EB >= 1:
		uint = "EB"
		return fmt.Sprintf("%.2f", float64(b)/float64(EB))
	case b/PB >= 1:
		uint = "PB"
		return fmt.Sprintf("%.2f", float64(b)/float64(PB))
	case b/TB >= 1:
		uint = "TB"
		return fmt.Sprintf("%.2f", float64(b)/float64(TB))
	case b/GB >= 1:
		uint = "GB"
		return fmt.Sprintf("%.2f", float64(b)/float64(GB))
	case b/MB >= 1:
		uint = "MB"
		return fmt.Sprintf("%.2f", float64(b)/float64(MB))
	case b/KB >= 1:
		uint = "KB"
		return fmt.Sprintf("%.2f", float64(b)/float64(KB))
	default:
		return fmt.Sprintf("%d", b)
	}
}

func (b ByteSize) HR() string {
	return b.HumanReadable()
}

func (b ByteSize) String() (s string) {
	return strings.Replace(b.HumanReadable(), " ", "", 1)
}

func (b ByteSize) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b *ByteSize) UnmarshalText(t []byte) (err error) {
	var val float64 = -1
	var unit string
	if len(t) == 0 {
		return &strconv.NumError{Func: fnUnmarshalText, Err: ErrBits}
	}
	// copy for error message
	t0 := t
	for i, v := range t {
		if !(v >= '0' && v <= '9' || v == '.') {
			unit = strings.TrimSpace(string(t[i:]))
			val, err = strconv.ParseFloat(string(t[:i]), 10)
			if err != nil {
				return
			}
			break
		}
	}

	switch unit {
	case "Kb", "Mb", "Gb", "Tb", "Pb", "Eb":
		goto BitsError
	}
	unit = strings.ToLower(unit)
	switch unit {
	case "", "b", "byte":
		defer func() {
			*b = ByteSize(val)
		}()
		if val >= 0 {
			return
		}
		val, err = strconv.ParseFloat(string(t), 10)
		if err != nil {
			return
		}
	case "k", "kb", "kilo", "kilobyte", "kilobytes":
		if val > float64(maxUint64/uint64(KB)) {
			goto Overflow
		}
		val *= float64(KB)

	case "m", "mb", "mega", "megabyte", "megabytes":
		if val > float64(maxUint64/uint64(MB)) {
			goto Overflow
		}
		val *= float64(MB)

	case "g", "gb", "giga", "gigabyte", "gigabytes":
		if val > float64(maxUint64/uint64(GB)) {
			goto Overflow
		}
		val *= float64(GB)

	case "t", "tb", "tera", "terabyte", "terabytes":
		if val > float64(maxUint64/uint64(TB)) {
			goto Overflow
		}
		val *= float64(TB)

	case "p", "pb", "peta", "petabyte", "petabytes":
		if val > float64(maxUint64/uint64(PB)) {
			goto Overflow
		}
		val *= float64(PB)
	case "e", "eb":
		if val > float64(maxUint64/uint64(EB)) {
			goto Overflow
		}
		val *= float64(EB)

	default:
		goto SyntaxError
	}
	*b = ByteSize(val)
	return nil

Overflow:
	*b = ByteSize(maxUint64)
	return &strconv.NumError{Func: fnUnmarshalText, Num: string(t0), Err: strconv.ErrRange}

SyntaxError:
	*b = 0
	return &strconv.NumError{Func: fnUnmarshalText, Num: string(t0), Err: strconv.ErrSyntax}

BitsError:
	*b = 0
	return &strconv.NumError{Func: fnUnmarshalText, Num: string(t0), Err: ErrBits}
}
