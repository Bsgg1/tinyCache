package tinyCache

// ByteView 存储真实的缓存值，b是只读的，所有通过byteslice返回一个copy的值
type ByteView struct {
	b []byte
}

func (b ByteView) Len() int {
	return len(b.b)
}
func (b *ByteView) ByteSlice() []byte {
	return copyBytes(b.b)
}

func (b *ByteView) String() string {
	return string(b.b)
}
func copyBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
