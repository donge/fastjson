package fastjson

import (
	"fmt"
	"testing"
	"time"
)

func TestGetStringConcurrent(t *testing.T) {
	const concurrency = 4
	data := []byte(largeFixture)

	ch := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			s := GetString(data, "non-existing-key")
			if s != "" {
				ch <- fmt.Errorf("unexpected non-empty string got: %q", s)
			}
			ch <- nil
		}()
	}

	for i := 0; i < concurrency; i++ {
		select {
		case <-time.After(time.Second * 5):
			t.Fatalf("timeout")
		case err := <-ch:
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
		}
	}
}

func TestGetBytesConcurrent(t *testing.T) {
	const concurrency = 4
	data := []byte(largeFixture)

	ch := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			b := GetBytes(data, "non-existing-key")
			if b != nil {
				ch <- fmt.Errorf("unexpected non-empty string got: %q", b)
			}
			ch <- nil
		}()
	}

	for i := 0; i < concurrency; i++ {
		select {
		case <-time.After(time.Second * 5):
			t.Fatalf("timeout")
		case err := <-ch:
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
		}
	}
}

func TestGetString(t *testing.T) {
	data := []byte(`{"foo":"bar", "baz": 1234}`)

	// normal path
	s := GetString(data, "foo")
	if s != "bar" {
		t.Fatalf("unexpected value obtained; got %q; want %q", s, "bar")
	}

	// non-existing path
	s = GetString(data, "foo", "zzz")
	if s != "" {
		t.Fatalf("unexpected non-empty value obtained: %q", s)
	}

	// invalid type
	s = GetString(data, "baz")
	if s != "" {
		t.Fatalf("unexpected non-empty value obtained: %q", s)
	}

	// invalid json
	s = GetString([]byte("invalid json"), "foobar", "baz")
	if s != "" {
		t.Fatalf("unexpected non-empty value obtained: %q", s)
	}
}

func TestGetStringImproved(t *testing.T) {
	data := []byte(`{
    "-soapenv": "http://schemas.xmlsoap.org/soap/envelope/",
    "Body": {
        "HIPMessageServer": {
            "action": "MES0018",
            "number": 0232
        },
		"messages": [
			{
				"content": "第一段信息",
				"role": "user"
			},
			{
				"content": "倒数第二段信息",
				"role": "assistant"
			},
			{
				"content": "最后一段信息",
				"role": "user"
			}
		]
    },
    "Header": "header here",
    "clin": "www.servicewall.com"
}`)

	// normal path
	s := GetStringImproved(data, "Body", "HIPMessageServer", "action")
	if s != "MES0018" {
		t.Fatalf("unexpected value obtained; got %q; want %q", s, "MES0018")
	}

	// object path
	s = GetStringImproved(data, "Body", "HIPMessageServer")
	if s != `{"action":"MES0018","number":0232}` {
		t.Fatalf("unexpected value obtained; got %q; want %q", s, `{"action":"MES0018","number":0232}`)
	}

	// number path
	s = GetStringImproved(data, "Body", "HIPMessageServer", "number")
	if s != "0232" {
		t.Fatalf("unexpected value obtained; got %q; want %q", s, "0232")
	}

	// non-existing path
	s = GetStringImproved(data, "foo", "zzz")
	if s != "" {
		t.Fatalf("unexpected non-empty value obtained: %q", s)
	}

	// invalid json
	s = GetStringImproved([]byte("invalid json"), "foobar", "baz")
	if s != "" {
		t.Fatalf("unexpected non-empty value obtained: %q", s)
	}

	s = GetStringImproved(data, "Body", "messages", "-1", "content")
	if s != "最后一段信息" {
		t.Fatalf("unexpected value obtained; got %q; want %q", s, "最后一段信息")
	}

	s = GetStringImproved(data, "Body", "messages", "-2", "content")
	if s != "倒数第二段信息" {
		t.Fatalf("unexpected value obtained; got %q; want %q", s, "倒数第二段信息")
	}

	s = GetStringImproved(data, "Body", "messages", "1", "content")
	if s != "倒数第二段信息" {
		t.Fatalf("unexpected value obtained; got %q; want %q", s, "倒数第二段信息")
	}

	s = GetStringImproved(data, "Body", "messages", "-100", "content")
	if s != "" {
		t.Fatalf("unexpected value obtained; got %q; want %q", s, "")
	}
}

func TestGetBytes(t *testing.T) {
	data := []byte(`{"foo":"bar", "baz": 1234}`)

	// normal path
	b := GetBytes(data, "foo")
	if string(b) != "bar" {
		t.Fatalf("unexpected value obtained; got %q; want %q", b, "bar")
	}

	// non-existing path
	b = GetBytes(data, "foo", "zzz")
	if b != nil {
		t.Fatalf("unexpected non-empty value obtained: %q", b)
	}

	// invalid type
	b = GetBytes(data, "baz")
	if b != nil {
		t.Fatalf("unexpected non-empty value obtained: %q", b)
	}

	// invalid json
	b = GetBytes([]byte("invalid json"), "foobar", "baz")
	if b != nil {
		t.Fatalf("unexpected non-empty value obtained: %q", b)
	}
}

func TestGetInt(t *testing.T) {
	data := []byte(`{"foo":"bar", "baz": 1234}`)

	// normal path
	n := GetInt(data, "baz")
	if n != 1234 {
		t.Fatalf("unexpected value obtained; got %d; want %d", n, 1234)
	}

	// non-existing path
	n = GetInt(data, "foo", "zzz")
	if n != 0 {
		t.Fatalf("unexpected non-zero value obtained: %d", n)
	}

	// invalid type
	n = GetInt(data, "foo")
	if n != 0 {
		t.Fatalf("unexpected non-zero value obtained: %d", n)
	}

	// invalid json
	n = GetInt([]byte("invalid json"), "foobar", "baz")
	if n != 0 {
		t.Fatalf("unexpected non-empty value obtained: %d", n)
	}
}

func TestGetFloat64(t *testing.T) {
	data := []byte(`{"foo":"bar", "baz": 12.34}`)

	// normal path
	f := GetFloat64(data, "baz")
	if f != 12.34 {
		t.Fatalf("unexpected value obtained; got %f; want %f", f, 12.34)
	}

	// non-existing path
	f = GetFloat64(data, "foo", "zzz")
	if f != 0 {
		t.Fatalf("unexpected non-zero value obtained: %f", f)
	}

	// invalid type
	f = GetFloat64(data, "foo")
	if f != 0 {
		t.Fatalf("unexpected non-zero value obtained: %f", f)
	}

	// invalid json
	f = GetFloat64([]byte("invalid json"), "foobar", "baz")
	if f != 0 {
		t.Fatalf("unexpected non-empty value obtained: %f", f)
	}
}

func TestGetFloat64Exact(t *testing.T) {
	data := []byte(`{"foo":"bar", "baz": 12.34, "zero":0}`)
	f, e := GetFloat64Exact(data, "foo")
	if e == nil || f != 0 {
		t.Fatalf("unexpected non-zero value obtained: %f", f)
	}
	f, e = GetFloat64Exact(data, "baz")
	if e != nil || f != 12.34 {
		t.Fatalf("unexpected non-zero value obtained: %f", f)
	}
	f, e = GetFloat64Exact(data, "zero")
	if e != nil || f != 0 {
		t.Fatalf("unexpected non-zero value obtained: %f", f)
	}
}

func TestGetBool(t *testing.T) {
	data := []byte(`{"foo":"bar", "baz": true}`)

	// normal path
	b := GetBool(data, "baz")
	if !b {
		t.Fatalf("unexpected value obtained; got %v; want %v", b, true)
	}

	// non-existing path
	b = GetBool(data, "foo", "zzz")
	if b {
		t.Fatalf("unexpected true value obtained")
	}

	// invalid type
	b = GetBool(data, "foo")
	if b {
		t.Fatalf("unexpected true value obtained")
	}

	// invalid json
	b = GetBool([]byte("invalid json"), "foobar", "baz")
	if b {
		t.Fatalf("unexpected true value obtained")
	}
}

func TestExists(t *testing.T) {
	data := []byte(`{"foo": [{"bar": 1234, "baz": 0}]}`)

	if !Exists(data, "foo") {
		t.Fatalf("cannot find foo")
	}
	if !Exists(data, "foo", "0") {
		t.Fatalf("cannot find foo[0]")
	}
	if !Exists(data, "foo", "0", "baz") {
		t.Fatalf("cannot find foo[0].baz")
	}

	if Exists(data, "foobar") {
		t.Fatalf("found unexpected foobar")
	}
	if Exists(data, "foo", "1") {
		t.Fatalf("found unexpected foo[1]")
	}
	if Exists(data, "foo", "0", "234") {
		t.Fatalf("found unexpected foo[0][234]")
	}
	if Exists(data, "foo", "bar") {
		t.Fatalf("found unexpected foo.bar")
	}

	if Exists([]byte(`invalid JSON`), "foo", "bar") {
		t.Fatalf("Exists returned true on invalid json")
	}
}

func TestParse(t *testing.T) {
	v, err := Parse(`{"foo": "bar"}`)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	str := v.String()
	if str != `{"foo":"bar"}` {
		t.Fatalf("unexpected value parsed: %q; want %q", str, `{"foo":"bar"}`)
	}
}

func TestParseBytes(t *testing.T) {
	v, err := ParseBytes([]byte(`{"foo": "bar"}`))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	str := v.String()
	if str != `{"foo":"bar"}` {
		t.Fatalf("unexpected value parsed: %q; want %q", str, `{"foo":"bar"}`)
	}
}

func TestMustParse(t *testing.T) {
	s := `{"foo":"bar"}`
	v := MustParse(s)
	str := v.String()
	if str != s {
		t.Fatalf("unexpected value parsed; %q; want %q", str, s)
	}

	v = MustParseBytes([]byte(s))
	if str != s {
		t.Fatalf("unexpected value parsed; %q; want %q", str, s)
	}

	if !causesPanic(func() { v = MustParse(`[`) }) {
		t.Fatalf("expected MustParse to panic")
	}

	if !causesPanic(func() { v = MustParseBytes([]byte(`[`)) }) {
		t.Fatalf("expected MustParse to panic")
	}
}

func causesPanic(fn func()) (p bool) {
	defer func() {
		if r := recover(); r != nil {
			p = true
		}
	}()
	fn()
	return
}
