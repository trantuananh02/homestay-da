package mysql

// import (
// 	"fmt"
// 	"reflect"
// )

// func GetListColumn(model interface{}, ignoreColumn []string, datetimeColumn []string) []string {
// 	t := reflect.TypeOf(model)
// 	if t.Kind() == reflect.Ptr {
// 		t = t.Elem()
// 	}
// 	result := make([]string, 0)
// 	for i := 0; i < t.NumField(); i++ {
// 		jsonTag := t.Field(i).Tag.Get("db")

// 		if jsonTag != "" && !utils.StringSliceContains(ignoreColumn, jsonTag) {
// 			if utils.StringSliceContains(datetimeColumn, jsonTag) {
// 				result = append(result, fmt.Sprintf("CAST(IFNULL(%s,'') AS char) as %s", jsonTag, jsonTag))
// 			} else {
// 				result = append(result, jsonTag)
// 			}
// 		}
// 	}
// 	return result
// }

// func GetListValues(model interface{}, ignoreColumn []string, datetimeColumn []string) []interface{} {
// 	v := reflect.ValueOf(model)
// 	if v.Kind() == reflect.Ptr {
// 		v = v.Elem()
// 	}
// 	t := reflect.TypeOf(model)
// 	if t.Kind() == reflect.Ptr {
// 		t = t.Elem()
// 	}
// 	result := make([]interface{}, 0)
// 	for i := 0; i < v.NumField(); i++ {
// 		value := v.Field(i)
// 		jsonTag := t.Field(i).Tag.Get("db")
// 		if jsonTag != "" && !utils.StringSliceContains(ignoreColumn, jsonTag) {
// 			if utils.StringSliceContains(datetimeColumn, jsonTag) && value.String() == "" {
// 				result = append(result, nil)
// 			} else {
// 				result = append(result, value.Interface())
// 			}
// 		}
// 	}
// 	return result
// }

// func GetListColumnAndValueForUpdate(old interface{}, new interface{}, datetimeColumn []string) ([]string, []interface{}) {
// 	oldV := reflect.ValueOf(old)
// 	if oldV.Kind() == reflect.Ptr {
// 		oldV = oldV.Elem()
// 	}
// 	newV := reflect.ValueOf(new)
// 	if newV.Kind() == reflect.Ptr {
// 		newV = newV.Elem()
// 	}
// 	t := reflect.TypeOf(old)
// 	if t.Kind() == reflect.Ptr {
// 		t = t.Elem()
// 	}
// 	listColumn := make([]string, 0)
// 	listValue := make([]interface{}, 0)
// 	for i := 0; i < t.NumField(); i++ {
// 		dbTag := t.Field(i).Tag.Get("db")
// 		oldValue := oldV.Field(i)
// 		newValue := newV.Field(i)
// 		if oldValue != newValue && dbTag != "" {
// 			listColumn = append(listColumn, dbTag)
// 			if utils.StringSliceContains(datetimeColumn, dbTag) && newValue.String() == "" {
// 				listValue = append(listValue, nil)
// 			} else {
// 				listValue = append(listValue, newValue.Interface())
// 			}
// 		}
// 	}
// 	return listColumn, listValue
// }

// func GetListValuesUpdate(model interface{}, updateColumns []string, datetimeColumns []string) []interface{} {
// 	v := reflect.ValueOf(model)
// 	if v.Kind() == reflect.Ptr {
// 		v = v.Elem()
// 	}
// 	t := reflect.TypeOf(model)
// 	if t.Kind() == reflect.Ptr {
// 		t = t.Elem()
// 	}
// 	result := make([]interface{}, 0)
// 	mapColumnsAndValues := make(map[string]interface{})
// 	for i := 0; i < v.NumField(); i++ {
// 		value := v.Field(i)
// 		jsonTag := t.Field(i).Tag.Get("db")
// 		if jsonTag != "" {
// 			if utils.StringSliceContains(datetimeColumns, jsonTag) && value.String() == "" {
// 				mapColumnsAndValues[jsonTag] = nil
// 			} else {
// 				mapColumnsAndValues[jsonTag] = value.Interface()
// 			}
// 		}
// 	}
// 	for _, column := range updateColumns {
// 		if value, isExist := mapColumnsAndValues[column]; isExist {
// 			result = append(result, value)
// 		} else {
// 			result = append(result, nil)
// 		}
// 	}
// 	return result
// }
