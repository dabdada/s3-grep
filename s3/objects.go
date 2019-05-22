package s3

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/dabdada/s3-grep/config"
)

// Object provides an interface for Objects in a S3 Bucket
type s3Object interface {
	GetObjectContent (*config.AWSSession, string) (string, int64, error)
}

// Object provides an Object with a Key
type Object struct {
	Object s3Object
	Key    string
}

// NewObject is a constructor for Objects
func NewObject(key string) Object {
	return Object{Key: key}
}

// ListObjects lists all objects in the specified bucket
func ListObjects(svc s3iface.S3API, bucketName string) ([]Object, error) {
	var objects []Object
	err := svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		for _, obj := range p.Contents {
			objects = append(objects, NewObject(*obj.Key))
		}
		return true
	})
	if err != nil {
		return []Object{}, err
	}

	return objects, nil
}

// GetObjectContent loads the content of a S3 object key into a buffer
func (o Object) GetObjectContent(session *config.AWSSession, bucketName string) ([]byte, int64, error) {
	if o.Key == "" {
		return []byte{}, 0, errors.New("Object has no Key")
	}
	buff := &aws.WriteAtBuffer{}
	downloader := s3manager.NewDownloader(session.Session)
	numBytes, err := downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(o.Key),
	})

	return buff.Bytes(), numBytes, err
}
