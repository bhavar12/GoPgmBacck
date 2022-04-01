package main

type MyHashMap struct {
	data []int
}

/** Initialize your data structure here. */
func Constructor() MyHashMap {
	return MyHashMap{data: []int{}}
}

/** value will always be non-negative. */
func (this *MyHashMap) Put(key int, value int) {
	this.data[key] = value
}

/** Returns the value to which the specified key is mapped, or -1 if this map contains no mapping for the key */
func (this *MyHashMap) Get(key int) int {

	data, ok := this.data[key]
	if !ok {
		return -1
	}
	return data
}

/** Removes the mapping of the specified value key if this map contains a mapping for the key */
func (this *MyHashMap) Remove(key int) {
	this.data[key] = 1, false
	delete(this.data, key)
}
func main() {
	//Your MyHashMap object will be instantiated and called as such:
	obj := Constructor()
	obj.Put(1, 15)
	param_2 := obj.Get(1)
	obj.Remove(1)
}
