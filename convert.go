// Copyright 2016 HenryLee. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Package apiware provides a tools which can bind the http/fasthttp request parameters to the structure and validate.
package apiware

import (
	"fmt"
	"reflect"
	"strconv"
)

// Type conversions for request parameters.
//
// convertAssign copies to dest the value in src, converting it if possible.
// An error is returned if the copy would result in loss of information.
// dest should be a pointer type.
func convertAssign(dest reflect.Value, src []string) (err error) {
	if len(src) == 0 {
		return nil
	}

	dest = reflect.Indirect(dest)
	if !dest.CanSet() {
		return fmt.Errorf("[apiware] %s can not be setted", dest.Type().Name())
	}

	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("%v", p)
		}
	}()

	switch dest.Interface().(type) {
	case string:
		dest.Set(reflect.ValueOf(src[0]))
		return nil

	case []string:
		dest.Set(reflect.ValueOf(src))
		return nil

	case []byte:
		dest.Set(reflect.ValueOf([]byte(src[0])))
		return nil

	case [][]byte:
		b := make([][]byte, 0, len(src))
		for _, s := range src {
			b = append(b, []byte(s))
		}
		dest.Set(reflect.ValueOf(b))
		return nil

	case bool:
		bol, err := strconv.ParseBool(src[0])
		if err != nil {
			return err
		}
		dest.Set(reflect.ValueOf(bol))
		return nil

	case []bool:
		b := make([]bool, 0, len(src))
		for _, s := range src {
			bol, err := strconv.ParseBool(s)
			if err != nil {
				return err
			}
			b = append(b, bol)
		}
		dest.Set(reflect.ValueOf(b))
		return nil
	}

	switch dest.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64, err := strconv.ParseInt(src[0], 10, dest.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("[apiware] converting type %T (%q) to a %s: %v", src, src[0], dest.Kind(), err)
		}
		dest.SetInt(i64)
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u64, err := strconv.ParseUint(src[0], 10, dest.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("[apiware] converting type %T (%q) to a %s: %v", src, src[0], dest.Kind(), err)
		}
		dest.SetUint(u64)
		return nil

	case reflect.Float32, reflect.Float64:
		f64, err := strconv.ParseFloat(src[0], dest.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("[apiware] converting type %T (%q) to a %s: %v", src, src[0], dest.Kind(), err)
		}
		dest.SetFloat(f64)
		return nil

	case reflect.Slice:
		switch dest.Type().Elem().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			for _, s := range src {
				i64, err := strconv.ParseInt(s, 10, dest.Type().Bits())
				if err != nil {
					err = strconvErr(err)
					return fmt.Errorf("[apiware] converting type %T (%q) to a %s: %v", src, s, dest.Kind(), err)
				}
				dest.Set(reflect.Append(dest, reflect.ValueOf(i64)))
			}
			return nil

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			for _, s := range src {
				u64, err := strconv.ParseUint(s, 10, dest.Type().Bits())
				if err != nil {
					err = strconvErr(err)
					return fmt.Errorf("[apiware] converting type %T (%q) to a %s: %v", src, s, dest.Kind(), err)
				}
				dest.Set(reflect.Append(dest, reflect.ValueOf(u64)))
			}
			return nil

		case reflect.Float32, reflect.Float64:
			for _, s := range src {
				f64, err := strconv.ParseFloat(s, dest.Type().Bits())
				if err != nil {
					err = strconvErr(err)
					return fmt.Errorf("[apiware] converting type %T (%q) to a %s: %v", src, s, dest.Kind(), err)
				}
				dest.Set(reflect.Append(dest, reflect.ValueOf(f64)))
			}
			return nil
		}
	}

	return fmt.Errorf("[apiware] unsupported storing type %T into type %s", src, dest.Kind())
}

func strconvErr(err error) error {
	if ne, ok := err.(*strconv.NumError); ok {
		return ne.Err
	}
	return err
}
