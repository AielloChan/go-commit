package model

type Store map[string]string

var store = Store{}

func GetInstance() *Store {
	return &store
}

func (store *Store) SaveValue(key string, value string) {
	(*store)[key] = value
}

func (store *Store) GetValue(key string) string {
	if value, ok := (*store)[key]; ok {
		return value
	}
	return ""
}
