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
	defaultDatasize.UnmarshalText([]byte(fmt.Sprintf("%d", bytes)))
	s = defaultDatasize.HumanReadable()
	return
}

func (b ByteSize) toBytes(unit ByteSize) float64 {
	v := b / unit
	r := b % unit
	return float64(uint64((float64(v)+float64(r)/float64(unit))*100)) / 100
}

func (b ByteSize) Bytes() uint64 {
	return uint64(b)
}

func (b ByteSize) KBytes() float64 {
	return b.toBytes(KB)
}

func (b ByteSize) MBytes() float64 {
	return b.toBytes(MB)
}

func (b ByteSize) GBytes() float64 {
	return b.toBytes(GB)
}

func (b ByteSize) TBytes() float64 {
	return b.toBytes(TB)
}

func (b ByteSize) PBytes() float64 {
	return b.toBytes(PB)
}

func (b ByteSize) EBytes() float64 {
	return b.toBytes(EB)
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
	t0 := t
	e := &strconv.NumError{Func: fnUnmarshalText, Num: string(t0), Err: ErrBits}
	var val float64 = -1
	var unit string
	if len(t) == 0 {
		return &strconv.NumError{Func: fnUnmarshalText, Err: ErrBits}
	}
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
	case "Kb", "Mb", "Gb", "Tb", "Pb", "Eb", "kB", "mB", "gB", "tB", "pB", "eB":
		return e
	}
	unit = strings.ToLower(unit)
	switch unit {
	case "", "b", "byte":
		if val == -1 {
			val, err = strconv.ParseFloat(string(t), 10)
			if err != nil {
				return
			}
		}
	case "k", "kb", "kilo", "kilobyte", "kilobytes":
		val *= float64(KB)
	case "m", "mb", "mega", "megabyte", "megabytes":
		val *= float64(MB)
	case "g", "gb", "giga", "gigabyte", "gigabytes":
		val *= float64(GB)
	case "t", "tb", "tera", "terabyte", "terabytes":
		val *= float64(TB)
	case "p", "pb", "peta", "petabyte", "petabytes":
		val *= float64(PB)
	case "e", "eb":
		val *= float64(EB)
	default:
		*b = 0
		return e
	}
	*b = ByteSize(val)
	return
}
