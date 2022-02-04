package utils

import "sync"

const maxCap = 1 << 11 // 2 kB

var PacketBytePool = sync.Pool{
	New: func() interface{} {
		return new(PacketByte)
	},
}

type PacketByte struct {
	data []byte
}

// Release если cap() больше maxCap то лучше его не ложить обратно в пул
// а дождаться пока GC его уничтожит,
// использование packetByte с cap() большого размера не эффективно
func (b *PacketByte) Release() {
	if cap(b.data) <= maxCap {
		b.data = b.data[:0]
		PacketBytePool.Put(b)
	}
}

func (b *PacketByte) Free() {
	b.data = b.data[:0]
}

func GetPacketByte() (b *PacketByte) {
	return PacketBytePool.Get().(*PacketByte)
}

// GetData получение массива байт из packetByte
func (b *PacketByte) GetData() []byte {
	cl := make([]byte, len(b.data))
	_ = copy(cl, b.data)
	return cl
}

// SetData копирует массив байт в packetByte
func (b *PacketByte) SetData(v []byte) {
	cl := make([]byte, len(v))
	b.data = cl
	copy(b.data, v)
}
