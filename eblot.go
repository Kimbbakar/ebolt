package eblot

import (
	"encoding/json"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

var (
	bucketDefault = []byte("default")
	boltFile      = "bbolt.db"
)

func getBoltClient(readOnly bool) *bolt.DB {
	option := &bolt.Options{
		Timeout:  time.Minute,
		ReadOnly: readOnly,
	}

	db, err := bolt.Open(boltFile, 0600, option)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func closeConnection(db *bolt.DB) {
	db.Close()
}

type EBoltClient struct {
	bucketName *string
}

func GetEbolt(bucketName *string) EBoltClient {
	c := EBoltClient{bucketName}
	c.Init()
	return c
}

type cachePayload struct {
	Value     interface{}
	CreatedAt time.Time
	Exp       *time.Time
}

func (p cachePayload) isExpired() bool {
	if p.Exp != nil {
		now := time.Now()
		return now.After(*p.Exp)
	}

	return false
}

func (c *EBoltClient) getBucketName() []byte {
	bucketName := bucketDefault
	if c.bucketName != nil {
		bucketName = []byte(*c.bucketName)
	}
	return bucketName
}

func (c *EBoltClient) Init() {
	db := getBoltClient(false)
	defer closeConnection(db)
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(c.getBucketName())
		if err != nil {
			logrus.Fatalf("[Ebolt/Init] Create bucket: %v", err.Error())
		}
		return err
	})

	go c.Sweep()
}

func (c *EBoltClient) Sweep() {
	db := getBoltClient(false)
	defer closeConnection(db)

	keyToDelete := []string{}

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(c.getBucketName())
		cr := bucket.Cursor()

		for key, value := cr.First(); key != nil; key, value = cr.Next() {
			shouldDelete := false
			if value == nil {
				shouldDelete = true
			} else {
				payload := cachePayload{}
				err := json.Unmarshal(value, &payload)
				if err != nil || payload.isExpired() {
					shouldDelete = true
				} else if payload.Exp != nil {
					go c.Expire(string(key), payload.Exp.Sub(time.Now()))
				}
			}

			if shouldDelete {
				keyToDelete = append(keyToDelete, string(key))
			}
		}
		return nil
	})

	if err != nil {
		logrus.Error("[Ebolt/Swap] Error: ", err.Error())
	}

	go c.DeleteMany(keyToDelete)
}

func (c *EBoltClient) Get(key string) interface{} {
	db := getBoltClient(true)
	defer closeConnection(db)
	payload := cachePayload{}

	err := db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(c.getBucketName())

		value := bucket.Get([]byte(key))
		if value != nil {
			json.Unmarshal(value, &payload)
		}
		return nil
	})
	if err != nil {
		logrus.Error("[Ebolt/Get] Error: ", err.Error())
	}

	if payload.isExpired() {
		payload.Value = nil
		go c.Delete(key)
	}

	return payload.Value
}

func (c *EBoltClient) Put(key string, value interface{}, ttl *time.Duration) error {
	db := getBoltClient(false)
	defer closeConnection(db)

	err := db.Batch(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(c.getBucketName())
		payload := cachePayload{Value: value}
		if ttl != nil {
			exp := time.Now().Add(*ttl)
			payload.Exp = &exp
			go c.Expire(key, *ttl)
		}

		byteValue, _ := json.Marshal(payload)
		return bucket.Put([]byte(key), byteValue)
	})
	if err != nil {
		logrus.Error("[Ebolt/Put] Error: ", err.Error())
	}

	return err
}

func (c *EBoltClient) Expire(key string, ttl time.Duration) error {
	<-time.After(ttl)
	return c.Delete(key)
}

func (c *EBoltClient) DeleteMany(keys []string) error {
	for _, key := range keys {
		if err := c.Delete(key); err != nil {
			return err
		}
	}
	return nil
}

func (c *EBoltClient) Delete(key string) error {
	db := getBoltClient(false)
	defer closeConnection(db)

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.getBucketName())
		return b.Delete([]byte(key))
	})
	if err != nil {
		logrus.Error("[Ebolt/Delete] Error: ", err.Error())
	}

	return err
}
