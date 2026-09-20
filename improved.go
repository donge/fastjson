package fastjson

// GetValImproved returns the value found by the improved lookup, starting from
// an already-parsed Value.
//
// Unlike Value.Get it applies the "improved" semantics used by
// GetValImproved(data, ...):
//   - it falls back to shorter key prefixes when the full path is missing;
//   - when an intermediate value is a JSON-encoded string, it drills into it.
//
// Callers that look up many paths in the same document can parse it once and
// reuse the resulting Value across calls, avoiding repeated full parses.
func (v *Value) GetValImproved(endIndex int, keys ...string) (ok bool, res *Value) {
	for i := len(keys); i > 0; i-- {
		x := v.Get(keys[:i]...)
		if x == nil {
			continue
		}
		if i == endIndex {
			return true, x
		}
		if x.Type() == TypeString {
			json := jsonStr2Json(x.GetStringBytes())
			if len(json) > 0 {
				sub, err := ParseBytes(json)
				if err != nil {
					return false, nil
				}
				return sub.GetValImproved(endIndex-i, keys[i:endIndex]...)
			}
		}
	}
	return false, nil
}

// GetStringImproved is the Value-based equivalent of
// GetStringImproved(data, keys...). It reuses the receiver instead of parsing.
func (v *Value) GetStringImproved(keys ...string) string {
	ok, x := v.GetValImproved(len(keys), keys...)
	if !ok || x == nil {
		return ""
	}
	var bs []byte
	switch x.Type() {
	case TypeString:
		bs = x.GetStringBytes()
	default:
		bs = x.MarshalTo(nil)
	}
	return string(bs)
}
