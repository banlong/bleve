package storage

import (
	"errors"
	"log"
	"time"
	"bleve/externals/boltdb"
	"bleve/tools"
	"bleve/config"
)


var (
	openDBs map[string]*bolt.DB
)

func init() {
	openDBs = make(map[string]*bolt.DB)

	//Location for *.db files
	tools.CreateDir(config.DB_DIR)

}

// Save bin data to Bolt
// - bucket: name of the bucket, constant as above
// - key: key of this data, this is the mediaId
// - data: binary data to save
// - dbName: mediaId
func Put(dbName string, bucket string, key string, data []byte) error {
	var err error
	if openDBs[dbName] == nil {
		openDBs[dbName] = new(bolt.DB)
		openDBs[dbName], err = bolt.Open(config.DB_DIR+dbName+".db", config.READWRITE, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			log.Printf("-- Put() cannot open %s \n", dbName)
		}
	}

	return openDBs[dbName].Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Println(err.Error())
			return err
		}

		if err = b.Put([]byte(key), data); err != nil {
			return err
		}
		log.Printf("-- saved %s to %s.db \n", key, dbName )
		return nil
	})
}

//Save a file into Bolt db.
func PutFile(dbName string, bucket string, key string, file string) error {
	//targetDB := openDBs[dbName]
	var err error
	if openDBs[dbName] == nil {
		openDBs[dbName] = new(bolt.DB)
		openDBs[dbName], err = bolt.Open(config.DB_DIR + dbName + ".db", config.READWRITE, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			log.Printf("-- PutFile() cannot open %s \n", dbName)
		}
	}

	data, err := tools.GetBytes(file)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	result := openDBs[dbName].Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Printf("-- create bucket %s failed! \n", bucket)
			return err
		}

		if err = b.Put([]byte(key), data); err != nil {
			//log.Printf("-- saved %s to %s failed! \n", file, dbName)
			log.Printf("-- saved %s failed \n", key)
			return err
		}
		//log.Printf("-- saved %s to %s.db \n", file, dbName )
		log.Printf("-- saved %s \n", key)
		return nil
	})

	return result
}

//Get binary val from Bolt
func Get(dbName string, bucket string, key string) (data []byte, err error) {
	//targetDB := openDBs[dbName]
	if openDBs[dbName] == nil {
		openDBs[dbName] = new(bolt.DB)
		openDBs[dbName], err = bolt.Open(config.DB_DIR+dbName+".db", config.READWRITE, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			log.Printf("-- %s not found\n", dbName+".db")
			return nil, errors.New("-- Movie resource file not found")
		}
	}

	err = openDBs[dbName].View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		//Bucket exist??
		if b == nil {
			return nil
		}

		val := b.Get([]byte(key))
		if len(val) > 0 {
			data = make([]byte, len(val))
			copy(data, val)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return data, err
}

//Delete the data using key
func Delete(dbName string, bucket string, key string) error {

	var err error
	if openDBs[dbName] == nil {
		openDBs[dbName] = new(bolt.DB)
		openDBs[dbName], err = bolt.Open(config.DB_DIR+dbName+".db", config.READWRITE, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			log.Printf("-- cannot open %s \n", dbName)
		}
	}

	// Delete the element using it key
	err = openDBs[dbName].Update(
		func(tx *bolt.Tx) error {
			return tx.Bucket([]byte(bucket)).Delete([]byte(key))
		})
	if err != nil {
		log.Printf("-- remove %s failed\n", dbName)
		return err
	}
	return nil
}
