package data


type MyValue struct {
	value string
}

func (my *MyValue) WriteAnswer(name string, value interface{}) error {
	my.value = value.(string)
	return nil
}

type DynamicTypeData = map[string]MyValue

func (dtd *DynamicTypeData)WriteValue(key string,value interface{}, type string){
	switch type {
	case "string":
	case "bool":
		
	}
}