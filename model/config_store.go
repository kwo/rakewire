package model

// C groups methods for accessing config entries.
var C = &configStore{}

type configStore struct{}

func (z *configStore) Get(tx Transaction) *Configuration {
	config := z.New()
	bData := tx.Bucket(bucketData, entityConfig)
	if data := bData.Get([]byte(idConfig)); data != nil {
		config.decode(data)
	}
	return config
}

func (z *configStore) New() *Configuration {
	return &Configuration{
		Values: make(map[string]string),
	}
}

func (z *configStore) Put(tx Transaction, config *Configuration) error {
	return saveObject(tx, entityConfig, config)
}
