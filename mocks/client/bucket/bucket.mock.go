package bucket

import (
	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) Upload(file []byte, filename string) (url string, key string, err error) {
	args := c.Called(file, filename)

	if args.Get(0) != nil {
		url = args.Get(0).(string)
	}

	if args.Get(1) != nil {
		key = args.Get(1).(string)
	}

	return url, key, args.Error(2)
}

func (c *ClientMock) Delete(in string) error {
	args := c.Called(in)

	return args.Error(0)
}
