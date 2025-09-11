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
		c := NewCache(3)
		c.Set("p", 100)
		c.Set("q", 200)
		c.Set("r", 300)
		c.Set("s", 400)

		_, ok := c.Get("p")
		require.False(t, ok)
		val, ok := c.Get("q")
		require.True(t, ok)
		require.Equal(t, 200, val)
		val, ok = c.Get("r")
		require.True(t, ok)
		require.Equal(t, 300, val)
		val, ok = c.Get("s")
		require.True(t, ok)
		require.Equal(t, 400, val)

		c = NewCache(3)
		c.Set("m", 11)
		c.Set("n", 22)
		c.Set("o", 33)

		c.Get("m")
		c.Get("o")

		c.Set("t", 44)

		_, ok = c.Get("n")
		require.False(t, ok)
		val, ok = c.Get("m")
		require.True(t, ok)
		require.Equal(t, 11, val)
		val, ok = c.Get("o")
		require.True(t, ok)
		require.Equal(t, 33, val)
		val, ok = c.Get("t")
		require.True(t, ok)
		require.Equal(t, 44, val)
	})
}

func TestCacheMultithreading(_ *testing.T) {
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
