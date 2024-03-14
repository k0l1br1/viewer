package main

import (
    "math"

	"github.com/k0l1br1/converter/bins"
)

func split(bs []bins.Bin, ts []float64, vs []float64) {
	if len(bs) != len(ts) || len(bs) != len(vs) {
		panic("invalid bins lenght")
	}
    for i, v := range bs {
        ts[i] = float64(v.Time)
        vs[i] = float64(v.Volume)
    }
}

func MinMax(in []float64) (float64, float64) {
    if len(in) == 0 {
        return 0, 0
    }
    min := in[0]
    max := in[0]
	for _, n := range in {
		if n > max {
            max = n
        }
        if n < min {
            min = n
        }
	}
    return min, max
}

func Sum(in []float64) float64 {
	var total float64
	for _, n := range in {
		total += n
	}
	return total
}

func Mean(in []float64) float64 {
	if len(in) == 0 {
		return 0
	}
	sum := Sum(in)
	return sum / float64(len(in))
}

func Variance(in []float64, mean float64) float64 {
	if len(in) == 0 {
		return 0
	}
    var v float64
	for _, n := range in {
		v += math.Pow(n - mean, 2)
	}

	return v / float64(len(in))
}

func NormalizeZ(in []float64) {
    mean := Mean(in)
    stdDev := math.Sqrt(Variance(in, mean))
	for i, n := range in {
		in[i] = (n - mean) / stdDev
	}
}

func NormalizeMinMax(in []float64) {
    min, max := MinMax(in)
    diff := max - min
    for i, n := range in {
        in[i] = (n - min) / diff
    }
}
