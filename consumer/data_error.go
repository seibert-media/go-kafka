// Copyright (c) 2019 //SEIBERT/MEDIA GmbH All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consumer

type HasCause interface {
	Cause() error
}

type HasData interface {
	Data() map[string]string
}

type DataError interface {
	error
	HasData
	HasCause
}

func AddDataToError(err error, data map[string]string) DataError {
	return &dataError{
		err:  err,
		data: data,
	}
}

type dataError struct {
	err  error
	data map[string]string
}

func (d *dataError) Cause() error {
	return d.err
}

func (d *dataError) Error() string {
	return d.err.Error()
}

func (d *dataError) Data() map[string]string {
	return d.data
}

func DataFromError(err error) map[string]string {
	data := make(map[string]string)
	for err != nil {
		hasData, ok := err.(HasData)
		if ok {
			for k, v := range hasData.Data() {
				data[k] = v
			}
		}
		hasCause, ok := err.(HasCause)
		if !ok {
			break
		}
		err = hasCause.Cause()
	}
	return data
}
