package greflect

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type Person struct {
	Name    string
	Age     int
	Friends []string
}

func TestDeepCopy(t *testing.T) {
	var a *Person
	assert.Nil(t, DeepCopy(a))
	// Test case 1: Copying nil should return nil
	assert.Nil(t, DeepCopy(nil))

	// Test case 2: Copying primitive types should return the same value
	assert.Equal(t, 42, DeepCopy(42))
	assert.Equal(t, "hello", DeepCopy("hello"))
	assert.Equal(t, true, DeepCopy(true))
	assert.Equal(t, 3.14, DeepCopy(3.14))

	// Test case 3: Copying slices should return a new slice with the same elements
	slice := []string{"apple", "banana", "cherry"}
	copySlice := DeepCopy(slice).([]string)
	assert.Equal(t, slice, copySlice)
	assert.False(t, &slice[0] == &copySlice[0]) // Ensure the addresses are different

	// Test case 4: Copying maps should return a new map with the same key-value pairs
	m := map[string]string{"key1": "value1", "key2": "value2"}
	copyMap := DeepCopy(m).(map[string]string)
	assert.Equal(t, m, copyMap)
	assert.True(t, reflect.DeepEqual(m["key1"], copyMap["key1"])) // Ensure the values are equal

	// Test case 5: Copying structs should return a new struct with the same values
	p1 := Person{Name: "Alice", Age: 30, Friends: []string{"Bob", "Charlie"}}
	p2 := DeepCopy(p1).(Person)
	assert.Equal(t, p1, p2)
	assert.False(t, &p1.Friends[0] == &p2.Friends[0]) // Ensure the addresses are different

	// Test case 6: Copying pointer types should return a new pointer with a deep copied value
	ptr := &slice
	copyPtr := DeepCopy(ptr).(*[]string)
	assert.Equal(t, *ptr, *copyPtr)
	assert.False(t, &(*ptr)[0] == &(*copyPtr)[0]) // Ensure the addresses are different

	// Test case 7: Copying nested structures should return a new structure with deep copied values
	nested := map[string]*Person{"person1": &p1}
	copyNested := DeepCopy(nested).(map[string]*Person)
	assert.Equal(t, nested, copyNested)
	assert.False(t, &nested["person1"].Friends[0] == &copyNested["person1"].Friends[0]) // Ensure the addresses are different
}
