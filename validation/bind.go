package validation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mqdvi-dp/go-common/errs"
	"github.com/mqdvi-dp/go-common/tracer"
)

const (
	param   string = "param"
	query   string = "query"
	header  string = "header"
	jsonTag string = "json"
	form    string = "form"
	ctxx    string = "context"
	skip    string = "skiptag"
	// equal     string = "="
	// semicolon string = ";"
	// comma     string = ","
)

func (v *validation) BindJSONAndValidate(ctx context.Context, bytes []byte, dest interface{}) error {
	trace, ctx := tracer.StartTraceWithContext(ctx, "Validation:BindJSON")
	defer trace.Finish()

	var err error
	// should pointer can bind
	t := reflect.TypeOf(dest)
	switch t.Kind() {
	case reflect.Ptr:
		t = t.Elem()

		switch t.Kind() {
		case reflect.Struct:
		default:
			err = fmt.Errorf("types not allowed. please change it to struct")
			tracer.SetError(ctx, err)
			return err
		}
	default:
		err = fmt.Errorf("should be pointer")
		tracer.SetError(ctx, err)
		return err
	}

	err = json.Unmarshal(bytes, dest)
	if err != nil {
		tracer.SetError(ctx, err)
		return err
	}

	err = v.Validate(ctx, dest)
	if err != nil {
		trace.SetError(err)
		return err
	}

	return nil
}

func (v *validation) BindAndValidate(c *gin.Context, dest interface{}) error {
	return v.BindAndValidateWithContext(context.Background(), c, dest)
}

func (v *validation) BindAndValidateWithContext(ctx context.Context, c *gin.Context, dest interface{}) error {
	trace, ctx := tracer.StartTraceWithContext(ctx, "Validation:Bind")
	defer trace.Finish()

	err := v.bind(ctx, c, dest)
	if err != nil {
		trace.SetError(err)
		return err
	}

	err = v.Validate(ctx, dest)
	if err != nil {
		trace.SetError(err)
		return err
	}

	return nil
}

func (v *validation) bind(ctx context.Context, c *gin.Context, dest interface{}, isAlreadyUnmarshal ...bool) error {
	alreadyUnmarshal := false // default is false
	if len(isAlreadyUnmarshal) > 0 {
		alreadyUnmarshal = isAlreadyUnmarshal[0]
	}

	// get value type
	vof := reflect.ValueOf(dest)
	tof := reflect.TypeOf(dest)

	switch vof.Type().Kind() {
	case reflect.Ptr: // when type is pointer, we should get the element of value
		if vof.IsNil() { // if the pointer is nil, we should set it to a new instance of the same type
			if !vof.CanSet() { // if the pointer is not settable, we should set it to a new instance of the same type
				vof = reflect.New(vof.Type()).Elem()
			}
			// set the value of pointer to a new instance of the same type
			vof.Set(reflect.New(vof.Type().Elem()))
		}

		vof = vof.Elem()
		tof = tof.Elem()
	case reflect.Slice: // when type is slice/array, recursive
		for i := 0; i < vof.Len(); i++ {
			err := v.bind(ctx, c, vof.Index(i).Interface(), alreadyUnmarshal)
			if err != nil {
				return err
			}
		}

		// finish the binding request
		return nil
	}
	// loop the struct field
	for i := 0; i < tof.NumField(); i++ {
		// get the struct field with i as position
		tf := tof.Field(i)
		// if the field is not exported (private), we should continue
		if !tf.IsExported() {
			continue
		}
		vf := vof.Field(i)

		// check if the field is skipped
		hasSkipTagJson := false
		if val, ok := tf.Tag.Lookup(skip); ok {
			if val == jsonTag {
				hasSkipTagJson = true
			}

			// if the tag value is boolean and true, we should continue
			// to check the next field. no need to bind the field
			if ok, err := strconv.ParseBool(val); ok && err == nil {
				continue
			}
		}

		// validate struct tag, is there has `json` tag? if yes, we need to unmarshal
		if val, ok := tf.Tag.Lookup(jsonTag); ok && !hasSkipTagJson {
			if alreadyUnmarshal {
				continue
			}
			// if the struct tag is `json:"-"`, we should continue
			if val == "-" {
				goto nonjson
			}

			if vf.IsZero() {
				body, err := io.ReadAll(c.Request.Body)
				defer c.Request.Body.Close()
				if err != nil {
					tracer.SetError(ctx, err)

					return errs.NewErrorWithCodeErr(err, errs.BAD_REQUEST)
				}

				// put the body again
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
				// unmarshal
				err = json.Unmarshal(body, dest)
				if err != nil {
					tracer.SetError(ctx, err)

					return errs.NewErrorWithCodeErr(err, errs.BAD_REQUEST)
				}

				alreadyUnmarshal = true
				continue
			}
		}

	nonjson:
		switch tf.Type.Kind() {
		case reflect.Ptr:
			switch tf.Type.Elem().Kind() {
			case reflect.Struct:
				err := v.bind(ctx, c, vf.Interface(), alreadyUnmarshal)
				if err != nil {
					return err
				}

				continue
			}
		case reflect.Struct:
			err := v.bind(ctx, c, vf.Interface(), alreadyUnmarshal)
			if err != nil {
				return err
			}

			continue
		case reflect.Slice:
			for j := 0; j < vf.Len(); j++ {
				err := v.bind(ctx, c, vf.Index(j).Interface(), alreadyUnmarshal)
				if err != nil {
					return err
				}

				continue
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Float32, reflect.Float64, reflect.String, reflect.Bool: // currently, we didn't support reflect.Map
		default:
			continue
		}

		var value string
		// when struct tag is url param
		if key, ok := tf.Tag.Lookup(param); ok {
			value = c.Param(key)
		}

		// when struct tag is query
		if key, ok := tf.Tag.Lookup(query); ok {
			value = c.Query(key)
		}

		// when struct tag is header
		if key, ok := tf.Tag.Lookup(header); ok {
			value = c.GetHeader(key)
		}

		// when struct tag is form post data (form-data or x-www-form-urlencoded)
		if key, ok := tf.Tag.Lookup(form); ok {
			value = c.PostForm(key)
		}

		// when struct tag is context
		if key, ok := tf.Tag.Lookup(ctxx); ok {
			if val, ok := c.Get(key); ok {
				switch v := val.(type) {
				case string:
					value = v
				case int:
					value = strconv.Itoa(v)
				case int64:
					value = strconv.FormatInt(v, 10)
				case float64:
					value = strconv.FormatFloat(v, 'f', -1, 64)
				case bool:
					value = strconv.FormatBool(v)
				}
			}
		}

		if value != "" {
			err := v.assignValue(value, tf, vf)
			if err != nil {
				tracer.SetError(ctx, err)

				return errs.NewErrorWithCodeErr(err, errs.BAD_REQUEST)
			}
		}
	}

	return nil
}
