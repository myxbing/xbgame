package packet

import (
	"errors"
	"math"
)
// 包结构
type Packet struct {
	offset uint
	data []byte
}

//获取包原始数据
func (this *Packet)Data()[]byte{
	return this.data
}

//获取包大小
func (this *Packet)Size()uint{
	return uint(len(this.data))
}

//获取当前包位置偏移量
func (this *Packet)Position()uint{
	return this.offset
}

//向后改变偏移量位置
func (this *Packet)Seek(offset uint){
	this.offset += offset
}

//======================包读取器========================
//读取一个字节，并下移一个位置
func (this *Packet)ReadByte() (byte, error){
	if this.offset >= this.Size() {
		return byte(0),errors.New("Out of packet size!")
	}

	cur := this.data[this.offset]
	this.offset++
	return cur,nil
}

//读取一个布尔值1:true，其他为false
func (this *Packet)ReadBool() (bool, error){
	b, err := this.ReadByte()
	if b == byte(1){
		return true, err
	}
	return false, err
}

//读取一个uint16值
func (this *Packet)ReadU16() (uint16, error){
	if this.offset+2 >= this.Size() {
		return uint16(0),errors.New("Out of packet size!")
	}

	buf := this.data[this.offset : this.offset+2]
	ret := uint16(buf[0])<<8 | uint16(buf[1])
	this.offset += 2
	return ret, nil
}

//读取一个int16值
func (this *Packet)Read16() (int16, error){
	val, err :=  this.ReadU16()
	return int16(val), err
}

//读取一个数据包，此包以uint16为包头大小
func (this *Packet)ReadBytes() ([]byte, error){
	if this.offset+2 >= this.Size() {
		return nil,errors.New("Out of packet size!")
	}
	size,_ :=this.ReadU16()
	if this.offset+uint(size) >= this.Size() {
		return nil,errors.New("Out of packet size!")
	}

	ret := this.data[this.offset : this.offset+uint(size)]
	this.offset += uint(size)
	return ret, nil
}

//读取一个字符串，此包以uint16为包头大小
func (this *Packet)ReadString() (string, error){
	buf, err :=  this.ReadBytes()
	if err != nil{
		return "", err
	}

	return string(buf), nil
}

//读取一个24位uint32值
func (this *Packet)ReadU24() (uint32, error){
	if this.offset+3 >= this.Size() {
		return uint32(0),errors.New("Out of packet size!")
	}

	buf := this.data[this.offset : this.offset+3]
	ret := uint32(buf[0])<<16 | uint32(buf[1])<<8 | uint32(buf[2])
	this.offset += 3
	return ret, nil
}

//读取一个24位int32值
func (this *Packet)Read24() (int32, error){
	val, err :=  this.ReadU24()
	return int32(val), err
}

//读取一个32位uint32值
func (this *Packet)ReadU32() (uint32, error){
	if this.offset+4 >= this.Size() {
		return uint32(0),errors.New("Out of packet size!")
	}

	buf := this.data[this.offset : this.offset+4]
	ret := uint32(buf[0])<<24 |  uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3])
	this.offset += 4
	return ret, nil
}

//读取一个32位int32值
func (this *Packet)Read32() (int32, error){
	val, err :=  this.ReadU32()
	return int32(val), err
}

//读取一个64位uint64值
func (this *Packet)ReadU64() (uint64, error){
	if this.offset+8 >= this.Size() {
		return uint64(0),errors.New("Out of packet size!")
	}

	buf := this.data[this.offset : this.offset+8]
	var ret uint64 = 0
	for i, v := range buf {
		ret |= uint64(v) << uint((7-i)*8)
	}
	this.offset += 8
	return ret, nil
}

//读取一个64位int64值
func (this *Packet)Read64() (int64, error){
	val, err :=  this.ReadU64()
	return int64(val), err
}

//读取一个32位float32值
func (this *Packet)ReadF32() (float32, error){
	val, err :=  this.ReadU32()
	if err != nil{
		return float32(0), err
	}
	return math.Float32frombits(val), nil
}

//读取一个64位float64值
func (this *Packet)ReadF64() (float64, error){
	val, err :=  this.ReadU64()
	if err != nil{
		return float64(0), err
	}
	return math.Float64frombits(val), nil
}

//======================包写入器========================
//以0填充的size个字节数据
func (this*Packet)Zeros(size int){
	zeros := make([]byte, size)
	this.data = append(this.data, zeros...)
}
//写入布尔值
func (this *Packet)WriteBool(val bool){
	if val{
		this.WriteByte(byte(1))
	}else{
		this.WriteByte(byte(0))
	}
}
//写入字节
func (this *Packet)WriteByte(val byte){
	this.data = append(this.data, val)
}
//写入字节数组，以uint16长度开头
func (this *Packet)WriteBytes(val []byte){
	this.WriteU16(uint16(len(val)))
	this.Write(val)
}

//直接写入二进制数据
func (this *Packet)Write(val []byte){
	this.data = append(this.data, val...)
}

//写入uint16数值
func (this *Packet) WriteU16(val uint16) {
	buf := make([]byte, 2)
	buf[0] = byte(val >> 8)
	buf[1] = byte(val)
	this.Write(buf)
}

//写入字符串，以uint16长度开头
func (this *Packet)WriteString(val string){
	this.WriteBytes([]byte(val))
}

//写入int16数值
func (this *Packet) Write16(val int16) {
	this.WriteU16(uint16(val))
}

//写入24位的uint32数值
func (this *Packet) WriteU24(val uint32) {
	buf := make([]byte, 3)
	buf[0] = byte(val >> 16)
	buf[1] = byte(val >> 8)
	buf[2] = byte(val)
	this.Write(buf)
}

//写入24位的int32数值
func (this *Packet) Write24(val int32) {
	this.WriteU24(uint32(val))
}

//写入32位的uint32数值
func (this *Packet) WriteU32(val uint32) {
	buf := make([]byte, 4)
	buf[0] = byte(val >> 24)
	buf[1] = byte(val >> 16)
	buf[2] = byte(val >> 8)
	buf[3] = byte(val)
	this.Write(buf)
}

//写入32位的int32数值
func (this *Packet) Write32(val int32) {
	this.WriteU32(uint32(val))
}

//写入64位的uint64数值
func (this *Packet) WriteU64(val uint64) {
	buf := make([]byte, 8)
	for i := range buf {
		buf[i] = byte(val >> uint((7-i)*8))
	}

	this.Write(buf)
}

//写入64位的int64数值
func (this *Packet) Write64(val int64) {
	this.WriteU64(uint64(val))
}

//写入32位float32数值
func (this *Packet) WriteF32(val float32) {
	v := math.Float32bits(val)
	this.WriteU32(v)
}

//写入64位float64数值
func (this *Packet) WriteF64(val float64) {
	v := math.Float64bits(val)
	this.WriteU64(v)
}

//=========================实例化=========================
//实例化一个读取实例
func NewReader(data []byte) *Packet {
	return &Packet{offset: 0, data: data}
}

//实例化一个写入实例
func NewWriter() *Packet {
	pkt := &Packet{offset: 0}
	pkt.data = make([]byte, 0, 128)
	return pkt
}
