package model

// C groups methods for accessing config entries.
var C = &configStore{}

type configStore struct{}

func (z *configStore) Get(tx Transaction) *Config {
	config := &Config{}
	bData := tx.Bucket(bucketData, entityConfig)
	if data := bData.Get([]byte(idConfig)); data != nil {
		config.decode(data)
	}
	return config
}

func (z *configStore) Put(tx Transaction, config *Config) error {
	return save(tx, entityConfig, config)
}
