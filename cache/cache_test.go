package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_store_and_get(t *testing.T) {
	t.Run("Test get and set values", func(t *testing.T) {
		cases := []struct {
			key      string
			value    any
			expected any
		}{
			{
				"int",
				1,
				1,
			},
			{
				"string",
				"string",
				"string",
			},
		}
		for _, tt := range cases {
			kv := &cache[any]{
				datas:      make(map[string]*data[any]),
				ttlDefault: DEFAULT_TTL_CACHE * time.Millisecond,
			}
			kv.Set(tt.key, tt.value)
			actual, _ := kv.Get(tt.key)

			assert.Equal(t, actual, tt.expected)
		}
	})
}

func Test_get_instance(t *testing.T) {
	t.Run("Test get instance", func(t *testing.T) {
		kv := GetInstance()
		assert.NotNil(t, kv)
	})
}

func Test_get_not_found(t *testing.T) {
	t.Run("Test Not Found", func(t *testing.T) {
		kv := &cache[any]{
			datas:      make(map[string]*data[any]),
			ttlDefault: DEFAULT_TTL_CACHE * time.Millisecond,
		}
		_, found := kv.Get("key")
		assert.False(t, found)
	})
}

func Test_get_not_expired(t *testing.T) {
	t.Run("Test Not Found", func(t *testing.T) {
		kv := &cache[any]{
			datas:      make(map[string]*data[any]),
			ttlDefault: DEFAULT_TTL_CACHE * time.Millisecond,
		}
		kv.SetWithTTL("key", 1, 0)
		_, found := kv.Get("key")
		assert.True(t, found)
	})
}
