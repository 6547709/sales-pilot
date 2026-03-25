package storage

import (
	"context"
	"errors"
	"io"
)

// BlobStore 附件存储抽象，当前未挂接业务；未来可替换为 S3/MinIO 实现
type BlobStore interface {
	Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
	Get(ctx context.Context, key string) (io.ReadCloser, string, error)
	Delete(ctx context.Context, key string) error
}

// LocalStore 本地磁盘占位实现（可扩展）
type LocalStore struct {
	Root string
}

func (LocalStore) Put(_ context.Context, _ string, _ io.Reader, _ int64, _ string) error {
	return nil
}

func (LocalStore) Get(_ context.Context, _ string) (io.ReadCloser, string, error) {
	return nil, "", errors.New("未实现")
}

func (LocalStore) Delete(_ context.Context, _ string) error {
	return nil
}
