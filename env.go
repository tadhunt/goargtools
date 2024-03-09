package goargtools

import (
	"fmt"
	"log"
	"os"
	"reflect"
)

var Debug = false

func EnvParse(args any) error {
	val := reflect.ValueOf(args)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	} else {
		return fmt.Errorf("bad args: not a pointer")
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("bad args: not a struct pointer")
	}

	for i := 0; i < val.NumField(); i++ {
		ftype := val.Type().Field(i)
		fval := val.Field(i)
		switch fval.Kind() {
		case reflect.String:
			tag, found := ftype.Tag.Lookup("env")
			if !found {
				if Debug {
					log.Printf("field: %s, no tag", ftype.Name)
				}
				break
			}


			eval, found := os.LookupEnv(tag)
			if !found {
				if Debug {
					log.Printf("field: %s, tag: %s: unset", ftype.Name, tag)
				}
				break
			}
			fval.SetString(eval)
			if Debug {
				log.Printf("field: %s, tag: %s newval %s", ftype.Name, tag, eval)
			}
		case reflect.Struct:
			EnvParse(fval.Addr().Interface())
		default:
			return fmt.Errorf("%s: unsupported kind %v", ftype.Name, fval.Kind())
		}
	}

	return nil
}
