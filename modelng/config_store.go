package modelng

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

func (z *configStore) Put(config *Config, tx Transaction) error {
	return save(entityConfig, config, tx)
}
