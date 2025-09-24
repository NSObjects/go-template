package code

import (
	"errors"
	"testing"

	errorslib "github.com/marmotedu/errors"
	"github.com/stretchr/testify/assert"
)

func TestLookup(t *testing.T) {
	coder, ok := Lookup(ErrDatabase)
	if assert.True(t, ok) {
		assert.Equal(t, ErrDatabase, coder.Code())
		assert.Equal(t, 500, coder.HTTPStatus())
		assert.Equal(t, "Database error", coder.String())
	}

	_, ok = Lookup(-1)
	assert.False(t, ok)
}

func TestHTTPStatus(t *testing.T) {
	assert.Equal(t, 200, HTTPStatus(0))
	assert.Equal(t, 400, HTTPStatus(ErrBadRequest))
	assert.Equal(t, 401, HTTPStatus(ErrUnauthorized))
	assert.Equal(t, 404, HTTPStatus(ErrUserNotFound))
	assert.Equal(t, 500, HTTPStatus(ErrDatabase))
	assert.Equal(t, 500, HTTPStatus(999999))
}

func TestWrapOrNewHelpers(t *testing.T) {
	err := WrapDatabaseError(errors.New("boom"), "query")
	assert.NotNil(t, err)

	coder := errorslib.ParseCoder(err)
	if assert.NotNil(t, coder) {
		assert.Equal(t, ErrDatabase, coder.Code())
	}

	err = WrapDatabaseError(nil, "query")
	assert.Nil(t, err)

	err = WrapError(nil, ErrForbidden, "access denied")
	assert.NotNil(t, err)
	coder = errorslib.ParseCoder(err)
	if assert.NotNil(t, coder) {
		assert.Equal(t, ErrForbidden, coder.Code())
	}
}
