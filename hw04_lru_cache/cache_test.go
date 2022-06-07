package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(5)

		values := []struct {
			key   Key
			value string
		}{
			{key: "aaa", value: "test_value_1"},
			{key: "bbb", value: "test_value_2"},
			{key: "ccc", value: "test_value_3"},
		}

		for _, v := range values {
			c.Set(v.key, v.value)

			val, exist := c.Get(v.key)
			require.True(t, exist)
			require.Equal(t, v.value, val)
		}

		c.Clear()

		for _, v := range values {
			val, exist := c.Get(v.key)
			require.False(t, exist)
			require.Nil(t, val)
		}
	})

	t.Run("full cache", func(t *testing.T) {
		c := NewCache(2)

		c.Set("1", 1)
		c.Set("2", 2)
		c.Set("3", 3)

		val, exist := c.Get("3")
		require.True(t, exist)
		require.Equal(t, 3, val)

		val, exist = c.Get("2")
		require.True(t, exist)
		require.Equal(t, 2, val)

		val, exist = c.Get("1")
		require.False(t, exist)
		require.Nil(t, val)
	})

	t.Run("full cache removes less used element", func(t *testing.T) {
		c := NewCache(3)

		c.Set("1", 1)
		c.Set("2", 2)
		c.Set("3", 3)

		c.Get("1")
		c.Get("2")
		c.Get("3")

		c.Get("1")
		c.Get("3")

		c.Set("4", 44)

		val, exist := c.Get("2")
		require.False(t, exist)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
