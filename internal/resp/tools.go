/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package resp

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
)

func ShaEncode(str string) (md string, err error) {
	h := sha256.New()
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

// Round 四舍五入
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
