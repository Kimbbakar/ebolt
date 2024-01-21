package ebolt_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kimbbakar/ebolt"
)

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

		t.Run("Overwrite value", func(t *testing.T) {
			ttl1 := time.Second * 5
			c := ebolt.GetEbolt(bucketName)

			c.Put("test_key", "test_value", &ttl1)
			value := c.Get("test_key")
			assert.NotNil(t, value)
			assert.Equal(t, "test_value", value)

			ttl2 := time.Second * 8
			c.Put("test_key", "new_value", &ttl2)
			time.Sleep(ttl1)

			value = c.Get("test_key")
			assert.Equal(t, value, "new_value")

			time.Sleep(ttl2)
			value = c.Get("test_key")
			assert.Equal(t, value, nil)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("String Value", func(t *testing.T) {
			c := ebolt.GetEbolt(bucketName)

			ttl := time.Minute * 10
			c.Put("test_key", "test_value", &ttl)
			value := c.Get("test_key")
			assert.NotNil(t, value)
			assert.Equal(t, "test_value", value)

			c.Delete("test_key")
			value = c.Get("test_key")
			assert.Nil(t, value)
		})
	})
}
