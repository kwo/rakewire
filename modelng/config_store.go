package modelng

// C groups methods for accessing config entries.
var C = &configStore{}

type configStore struct{}

func (z *configStore) Delete(id string, tx Transaction) error {
	return delete(entityConfig, id, tx)
}

func (z *configStore) GetByID(id string, tx Transaction) *Config {
	bData := tx.Bucket(bucketData, entityConfig)
	if data := bData.Get([]byte(id)); data != nil {
		config := &Config{}
		if err := config.decode(data); err == nil {
			return config
		}
	}
	return nil
}

func (z *configStore) New(name, value string) *Config {
	return &Config{Name: name, Value: value}
}

func (z *configStore) Save(config *Config, tx Transaction) error {
	return save(entityConfig, config, tx)
}
