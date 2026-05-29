package utils

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"Tiktok/internal/config"

	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/downloader"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
)

const qiniuPrefix = "qiniu://"

func UploadToQiNiu(ctx context.Context, reader io.Reader, objectName string, fileName string) (string, error) {
	creds := credentials.NewCredentials(config.Cfg.QiNiu.AccessKey, config.Cfg.QiNiu.SecretKey)

	var opts http_client.Options
	opts.Credentials = creds
	if err := opts.SetBucketHosts(http_client.DefaultBucketHosts()); err != nil {
		return "", err
	}

	upManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: opts,
	})

	err := upManager.UploadReader(ctx, reader, &uploader.ObjectOptions{
		BucketName: config.Cfg.QiNiu.Bucket,
		ObjectName: &objectName,
		FileName:   fileName,
	}, nil)
	if err != nil {
		return "", err
	}

	return qiniuPrefix + objectName, nil
}

func SignQiNiuURL(key string) string {
	if len(key) < len(qiniuPrefix) || key[:len(qiniuPrefix)] != qiniuPrefix {
		return key
	}
	objectName := key[len(qiniuPrefix):]

	creds := credentials.NewCredentials(config.Cfg.QiNiu.AccessKey, config.Cfg.QiNiu.SecretKey)
	signer := downloader.NewCredentialsSigner(creds)

	rawURL := fmt.Sprintf("%s/%s", config.Cfg.QiNiu.Domain, objectName)
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return key
	}

	err = signer.Sign(context.Background(), parsed, &downloader.SignOptions{
		TTL: 7 * 24 * time.Hour,
	})
	if err != nil {
		return key
	}

	return parsed.String()
}
