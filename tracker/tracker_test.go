package tracker

import (
	"reflect"
	"testing"
)

func TestGetPort(t *testing.T) {
	addresses := []string{
		"1.2.3.4:65536",
		"1.2.3.4:36",
		"1.2.3.4:1",
		"1.2.3.4:80",
		"1.2.3.4:95",
		"1.2.3.4:123456",
		"1.2.3.4:111",
		"1.2.3.4:8080",
		"1.2.3.4:",
		"1.2.3.4:0",
	}

	answers := []int32{65536, 36, 1, 80, 95, 123456, 111, 8080, 0, 0}

	for i, add := range addresses {
		ret := getPort(add)
		if answers[i] != ret {
			t.Errorf("\naddress: %v\nport: %v\nfunc return: %v", add, answers[i], ret)
		}
	}
}

func TestToBit(t *testing.T) {
	arr := [][]byte{
		[]byte{0, 0, 0},
		[]byte{0, 0, 1},
		[]byte{1, 1, 0},
		[]byte{1, 0, 1},
		[]byte{255, 255, 255},
		[]byte{0, 50, 100},
		[]byte{1, 3, 2},
		[]byte{1, 2, 3},
	}

	answers := [][]byte{
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
		[]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
		[]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		[]byte{0, 0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]byte{0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
		[]byte{1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
	}

	for i, testArr := range arr {
		res := toBit(testArr)
		if !reflect.DeepEqual(answers[i], res) {
			t.Errorf("\nbyter arr: %v\nshould be: %v\nbut got  : %v\n", testArr, answers[i], res)
		}
	}
}

func TestToInt(t *testing.T) {
	arr := [][]byte{
		[]byte{0, 0},
		[]byte{0, 1},
		[]byte{255, 255},
		[]byte{0, 123},
		[]byte{123, 0},
		[]byte{0, 0, 0, 0},
		[]byte{0, 0, 0, 1},
		[]byte{255, 255, 255, 255},
		[]byte{0, 0, 0, 123},
		[]byte{123, 0, 0, 0},
		[]byte{0, 0, 0, 0, 0, 0, 0, 0},
		[]byte{0, 0, 0, 0, 0, 0, 0, 1},
		[]byte{255, 255, 255, 255, 255, 255, 255, 255},
		[]byte{0, 0, 0, 0, 0, 0, 0, 123},
		[]byte{123, 0, 0, 0, 0, 0, 0, 0},
	}

	answer := []int64{
		0, 1, 65535, 123, (123) * (1 << 8),
		0, 1, -1, 123, 123 * (1 << 24),
		0, 1, -1, 123, 123 * (1 << 56),
	}

	for i, testArr := range arr {
		res := toInt(testArr)
		switch res.(type) {
		case int32:
			if int32(answer[i]) != res.(int32) {
				t.Errorf("\nbyter arr: %v\nshould be: %v\nbut got  : %v\n", testArr, answer[i], res)
			}
		case int64:
			if int64(answer[i]) != res.(int64) {
				t.Errorf("\nbyter arr: %v\nshould be: %v\nbut got  : %v\n", testArr, answer[i], res)
			}
		}
	}
}

func TestToBuf(t *testing.T) {
	arr := []int64{
		0, 1, 65535, 123, (123) * (1 << 8),
		0, 1, -1, 123, 123 * (1 << 24),
		0, 1, -1, 123, 123 * (1 << 56),
	}

	answer := [][]byte{
		[]byte{0, 0, 0, 0},
		[]byte{0, 0, 0, 1},
		[]byte{0, 0, 255, 255},
		[]byte{0, 0, 0, 123},
		[]byte{0, 0, 123, 0},
		[]byte{0, 0, 0, 0},
		[]byte{0, 0, 0, 1},
		[]byte{255, 255, 255, 255},
		[]byte{0, 0, 0, 123},
		[]byte{123, 0, 0, 0},
		[]byte{0, 0, 0, 0, 0, 0, 0, 0},
		[]byte{0, 0, 0, 0, 0, 0, 0, 1},
		[]byte{255, 255, 255, 255, 255, 255, 255, 255},
		[]byte{0, 0, 0, 0, 0, 0, 0, 123},
		[]byte{123, 0, 0, 0, 0, 0, 0, 0},
	}

	for i, val := range arr {
		var res []byte

		if i < 10 {
			res = ToBuf(int32(val))
		} else {
			res = ToBuf(int64(val))
		}

		if !reflect.DeepEqual(res, answer[i]) {
			t.Errorf("\n    val: %v\nshould be: %v\nbut got  : %v\n", val, answer[i], res)
		}
	}
}
