package marshaler

import (
	"errors"
	"testing"
	"time"

	"github.com/joseph-bing-han/gocache/store"
	mocksCache "github.com/joseph-bing-han/gocache/test/mocks/cache"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack"
)

type testCacheValue struct {
	Hello string
}

func TestNew(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cache := mocksCache.NewMockCacheInterface(ctrl)

	// When
	marshaler := New(cache)

	// Then
	assert.IsType(t, new(Marshaler), marshaler)
	assert.Equal(t, cache, marshaler.cache)
}

func TestGetWhenStoreReturnsSliceOfBytes(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cacheValue := &testCacheValue{
		Hello: "world",
	}

	cacheValueBytes, err := msgpack.Marshal(cacheValue)
	if err != nil {
		assert.Error(t, err)
	}

	cache := mocksCache.NewMockCacheInterface(ctrl)
	cache.EXPECT().Get("my-key").Return(cacheValueBytes, nil)

	marshaler := New(cache)

	// When
	value, err := marshaler.Get("my-key", new(testCacheValue))

	// Then
	assert.Nil(t, err)
	assert.Equal(t, cacheValue, value)
}

func TestGetWhenStoreReturnsString(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cacheValue := &testCacheValue{
		Hello: "world",
	}

	cacheValueBytes, err := msgpack.Marshal(cacheValue)
	if err != nil {
		assert.Error(t, err)
	}

	cache := mocksCache.NewMockCacheInterface(ctrl)
	cache.EXPECT().Get("my-key").Return(string(cacheValueBytes), nil)

	marshaler := New(cache)

	// When
	value, err := marshaler.Get("my-key", new(testCacheValue))

	// Then
	assert.Nil(t, err)
	assert.Equal(t, cacheValue, value)
}

func TestGetWhenUnmarshalingError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cache := mocksCache.NewMockCacheInterface(ctrl)
	cache.EXPECT().Get("my-key").Return("unknown-string", nil)

	marshaler := New(cache)

	// When
	value, err := marshaler.Get("my-key", new(testCacheValue))

	// Then
	assert.NotNil(t, err)
	assert.Nil(t, value)
}

func TestGetWhenNotFoundInStore(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := errors.New("Unable to find item in store")

	cache := mocksCache.NewMockCacheInterface(ctrl)
	cache.EXPECT().Get("my-key").Return(nil, expectedErr)

	marshaler := New(cache)

	// When
	value, err := marshaler.Get("my-key", new(testCacheValue))

	// Then
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, value)
}

func TestSetWhenStruct(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cacheValue := &testCacheValue{
		Hello: "world",
	}

	options := &store.Options{
		Expiration: 5 * time.Second,
	}

	cache := mocksCache.NewMockCacheInterface(ctrl)
	cache.EXPECT().Set("my-key", []byte{0x81, 0xa5, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0xa5, 0x77, 0x6f, 0x72, 0x6c, 0x64}, options).Return(nil)

	marshaler := New(cache)

	// When
	err := marshaler.Set("my-key", cacheValue, options)

	// Then
	assert.Nil(t, err)
}

func TestSetWhenString(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cacheValue := "test"

	options := &store.Options{
		Expiration: 5 * time.Second,
	}

	cache := mocksCache.NewMockCacheInterface(ctrl)
	cache.EXPECT().Set("my-key", []byte{0xa4, 0x74, 0x65, 0x73, 0x74}, options).Return(nil)

	marshaler := New(cache)

	// When
	err := marshaler.Set("my-key", cacheValue, options)

	// Then
	assert.Nil(t, err)
}

func TestSetWhenError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cacheValue := "test"

	options := &store.Options{
		Expiration: 5 * time.Second,
	}

	expectedErr := errors.New("An unexpected error occurred")

	cache := mocksCache.NewMockCacheInterface(ctrl)
	cache.EXPECT().Set("my-key", []byte{0xa4, 0x74, 0x65, 0x73, 0x74}, options).Return(expectedErr)

	marshaler := New(cache)

	// When
	err := marshaler.Set("my-key", cacheValue, options)

	// Then
	assert.Equal(t, expectedErr, err)
}

func TestDelete(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cache := mocksCache.NewMockCacheInterface(ctrl)
	cache.EXPECT().Delete("my-key").Return(nil)

	marshaler := New(cache)

	// When
	err := marshaler.Delete("my-key")

	// Then
	assert.Nil(t, err)
}

func TestDeleteWhenError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := errors.New("Unable to delete key")

	cache := mocksCache.NewMockCacheInterface(ctrl)
	cache.EXPECT().Delete("my-key").Return(expectedErr)

	marshaler := New(cache)

	// When
	err := marshaler.Delete("my-key")

	// Then
	assert.Equal(t, expectedErr, err)
}

func TestInvalidate(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	options := store.InvalidateOptions{
		Tags: []string{"tag1"},
	}

	cache := mocksCache.NewMockCacheInterface(ctrl)
	cache.EXPECT().Invalidate(options).Return(nil)

	marshaler := New(cache)

	// When
	err := marshaler.Invalidate(options)

	// Then
	assert.Nil(t, err)
}

func TestInvalidatingWhenError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	options := store.InvalidateOptions{
		Tags: []string{"tag1"},
	}

	expectedErr := errors.New("Unexpected error when invalidating data")

	cache := mocksCache.NewMockCacheInterface(ctrl)
	cache.EXPECT().Invalidate(options).Return(expectedErr)

	marshaler := New(cache)

	// When
	err := marshaler.Invalidate(options)

	// Then
	assert.Equal(t, expectedErr, err)
}

func TestClear(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cache := mocksCache.NewMockCacheInterface(ctrl)
	cache.EXPECT().Clear().Return(nil)

	marshaler := New(cache)

	// When
	err := marshaler.Clear()

	// Then
	assert.Nil(t, err)
}

func TestClearWhenError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := errors.New("An unexpected error occurred")

	cache := mocksCache.NewMockCacheInterface(ctrl)
	cache.EXPECT().Clear().Return(expectedErr)

	marshaler := New(cache)

	// When
	err := marshaler.Clear()

	// Then
	assert.Equal(t, expectedErr, err)
}
