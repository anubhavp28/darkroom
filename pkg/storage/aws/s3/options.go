package s3

import "***REMOVED***/darkroom/core/pkg/storage"

// Option represents the S3 storage options
type Option func(s *Storage)

// WithBucketName sets the bucket name
func WithBucketName(name string) Option {
	return func(s *Storage) {
		s.bucketName = name
	}
}

// WithBucketRegion sets the bucket region
func WithBucketRegion(region string) Option {
	return func(s *Storage) {
		s.bucketRegion = region
	}
}

// WithAccessKey sets the bucket accessKey
func WithAccessKey(accessKey string) Option {
	return func(s *Storage) {
		s.accessKey = accessKey
	}
}

// WithSecretKey sets the bucket secretKey
func WithSecretKey(secretKey string) Option {
	return func(s *Storage) {
		s.secretKey = secretKey
	}
}

// WithHystrixCommand sets the bucket hystrixCmd
func WithHystrixCommand(hytrixCmd storage.HystrixCommand) Option {
	return func(s *Storage) {
		s.hystrixCmd = hytrixCmd
	}
}
