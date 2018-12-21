package _strings

import (
	"fmt"
	"testing"
)

func TestRoundStr(t *testing.T) {
	num := "12"
	result := RoundStrWithNoError(num, 3)
	expect := "12.000"
	if result != expect {
		t.Fail()
		fmt.Printf("expect:%s result:%s \r\n", expect, result)
	}

	num = "12.01"
	result = RoundStrWithNoError(num, 3)
	expect = "12.010"
	if result != expect {
		t.Fail()
		fmt.Printf("expect:%s result:%s \r\n", expect, result)
	}

	num = "12.9995"
	result = RoundStrWithNoError(num, 3)
	expect = "13.000"
	if result != expect {
		t.Fail()
		fmt.Printf("expect:%s result:%s \r\n", expect, result)
	}

	num = "12.0155"
	result = RoundStrWithNoError(num, 3)
	expect = "12.016"
	if result != expect {
		t.Fail()
		fmt.Printf("expect:%s result:%s \r\n", expect, result)
	}

	num = "12.0115"
	result = RoundStrWithNoError(num, 3)
	expect = "12.012"
	if result != expect {
		t.Fail()
		fmt.Printf("expect:%s result:%s \r\n", expect, result)
	}

	num = "12.0125"
	result = RoundStrWithNoError(num, 3)
	expect = "12.013"
	if result != expect {
		t.Fail()
		fmt.Printf("expect:%s result:%s \r\n", expect, result)
	}

	num = "12.5678"
	result = RoundStrWithNoError(num, 0)
	expect = "13"
	if result != expect {
		t.Fail()
		fmt.Printf("expect:%s result:%s \r\n", expect, result)
	}

	num = "12.12345"
	result = RoundStrWithNoError(num, 3)
	expect = "12.123"
	if result != expect {
		t.Fail()
		fmt.Printf("expect:%s result:%s \r\n", expect, result)
	}

	num = "12.12345"
	result = RoundStrWithNoError(num, 4)
	expect = "12.1235"
	if result != expect {
		t.Fail()
		fmt.Printf("expect:%s result:%s \r\n", expect, result)
	}

	num = "-12"
	result = RoundStrWithNoError(num, 3)
	expect = "-12.000"
	if result != expect {
		t.Fail()
		fmt.Printf("expect:%s result:%s \r\n", expect, result)
	}

	num = "-12.01"
	result = RoundStrWithNoError(num, 3)
	expect = "-12.010"
	if result != expect {
		t.Fail()
		fmt.Printf("expect:%s result:%s \r\n", expect, result)
	}

	num = "-12.9995"
	result = RoundStrWithNoError(num, 3)
	expect = "-13.000"
	if result != expect {
		t.Fail()
		fmt.Printf("expect:%s result:%s \r\n", expect, result)
	}

	num = "-12.0155"
	result = RoundStrWithNoError(num, 3)
	expect = "-12.016"
	if result != expect {
		t.Fail()
		fmt.Printf("expect:%s result:%s \r\n", expect, result)
	}

	num = "-12.0115"
	result = RoundStrWithNoError(num, 3)
	expect = "-12.012"
	if result != expect {
		t.Fail()
		fmt.Printf("expect:%s result:%s \r\n", expect, result)
	}

	num = "-12.0125"
	result = RoundStrWithNoError(num, 3)
	expect = "-12.013"
	if result != expect {
		t.Fail()
		fmt.Printf("expect:%s result:%s \r\n", expect, result)
	}

	num = "-12.5678"
	result = RoundStrWithNoError(num, 0)
	expect = "-13"
	if result != expect {
		t.Fail()
		fmt.Printf("expect:%s result:%s \r\n", expect, result)
	}

	num = "-12.12345"
	result = RoundStrWithNoError(num, 3)
	expect = "-12.123"
	if result != expect {
		t.Fail()
		fmt.Printf("expect:%s result:%s \r\n", expect, result)
	}

	num = "-12.12345"
	result = RoundStrWithNoError(num, 4)
	expect = "-12.1235"
	if result != expect {
		t.Fail()
		fmt.Printf("expect:%s result:%s \r\n", expect, result)
	}
}
