package ebolt_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kimbbakar/ebolt"
)

var wg sync.WaitGroup
var count [10000]int

func TestEBolt(t *testing.T) {
	bucketName := "test_bucket"
	ebolt.InitEbolt(&bucketName)

	t.Run("Put and Get", func(t *testing.T) {

		t.Run("String Value", func(t *testing.T) {
			c := ebolt.GetEbolt(bucketName)
			c.Put("test_key", "test_value", nil)
			value := c.Get("test_key")
			assert.NotNil(t, value)
			assert.Equal(t, "test_value", value)
		})

		t.Run("Map Value", func(t *testing.T) {
			value := map[string]interface{}{
				"test_key": "value",
			}

			c := ebolt.GetEbolt(bucketName)
			c.Put("test_key", value, nil)
			response := c.Get("test_key")
			assert.NotNil(t, response)
			assert.Equal(t, value, response)
		})

		t.Run("Number Value", func(t *testing.T) {
			c := ebolt.GetEbolt(bucketName)
			c.Put("test_key", 101, nil)
			response := c.Get("test_key")
			assert.NotNil(t, response)
			assert.Equal(t, 101, int(response.(float64)))
		})
	})

	t.Run("Put and Get with expiry", func(t *testing.T) {
		t.Run("String Value", func(t *testing.T) {
			ttl := time.Second * 5
			c := ebolt.GetEbolt(bucketName)
			c.Put("test_key", "test_value", &ttl)
			value := c.Get("test_key")
			assert.NotNil(t, value)
			assert.Equal(t, "test_value", value)

			time.Sleep(ttl)
			value = c.Get("test_key")
			assert.Equal(t, value, nil)
		})
	})
}
