package fastjson

import "testing"

func TestValueGetStringImproved(t *testing.T) {
	data := []byte(`{"a":"x","b":1,"c":1.5,"d":true,"e":null,"f":{"g":"h"},"i":["j",{"k":"l"}],"m":"","n":"{\"o\":\"p\"}","esc":"a\"b\\c\n","neg":-2,"deep":"{\"p\":{\"q\":\"r\"}}","arr":[{"z":"1"},{"z":"2"}]}`)

	v, err := ParseBytes(data)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	paths := [][]string{
		{"a"}, {"b"}, {"c"}, {"d"}, {"e"}, {"f"}, {"f", "g"},
		{"i"}, {"i", "0"}, {"i", "1"}, {"i", "1", "k"}, {"i", "-1", "k"},
		{"m"}, {"n"}, {"n", "o"}, {"esc"}, {"neg"}, {"deep", "p", "q"},
		{"arr", "1", "z"}, {"zz"}, {"f", "zz"}, {"i", "5", "k"}, {"missing", "path"},
	}

	for _, p := range paths {
		want := GetStringImproved(data, p...)
		got := v.GetStringImproved(p...)
		if want != got {
			t.Errorf("path %v: want %q, got %q", p, want, got)
		}
	}

	// 复用同一 Value 多次取值, 结果应保持稳定
	for i := 0; i < 3; i++ {
		if got := v.GetStringImproved("f", "g"); got != "h" {
			t.Fatalf("reuse #%d: want h, got %q", i, got)
		}
	}
}

func TestValueGetStringImprovedNil(t *testing.T) {
	var v *Value
	if got := v.GetStringImproved("a"); got != "" {
		t.Fatalf("nil receiver: want empty, got %q", got)
	}
}
