/*
 *
 * tools.go
 * tools
 *
 * Created by lin on 2018/12/10 5:19 PM
 * Copyright © 2017-2018 PYL. All rights reserved.
 *
 */

package tools

import (
	"crypto/md5"
	"encoding/hex"
	"math"
)

func Md5Encode(str string) (md string, err error) {
	h := md5.New()
	_, err = h.Write([]byte(str))
	md = hex.EncodeToString(h.Sum(nil))
	return
}

type Coin int64

const Ratio float64 = 100000

func (c Coin) Yuan() float64 {
	return float64(c) / Ratio
}

func ToCoin(yuan float64) Coin {
	return Coin(Round(yuan*Ratio, 0))
}

//Round 四舍五入
func Round(x float64, pre int) float64 {
	a := x
	switch pre {
	case 0:
		_, rem := math.Modf(a)
		if rem >= 0.5 {
			a = math.Ceil(a)
		} else {
			a = math.Floor(a)
		}
	default:
		preNumber := math.Pow10(pre)
		a = x * preNumber
		_, rem := math.Modf(a)
		if rem >= 0.5 {
			a = math.Ceil(a)
		} else {
			a = math.Floor(a)
		}
		a = a / preNumber

	}
	return a
}


