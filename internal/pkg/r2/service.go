package r2

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"serica-go/internal/conf"
	"serica-go/internal/pkg/httputil"
)

type Service struct {
	cfg conf.R2
}

func NewService(cfg *conf.Bootstrap) *Service {
	return &Service{cfg: cfg.R2}
}

func (s *Service) UploadFile(ctx context.Context, r2Key, localFilePath string) error {
	data, err := os.ReadFile(localFilePath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	return s.UploadBuffer(ctx, r2Key, data)
}

func (s *Service) UploadBuffer(ctx context.Context, r2Key string, data []byte) error {
	endpoint := s.cfg.Endpoint
	if endpoint == "" {
		endpoint = fmt.Sprintf("https://%s.r2.cloudflarestorage.com", s.cfg.AccountID)
	}

	url := fmt.Sprintf("%s/%s/%s", endpoint, s.cfg.BucketName, r2Key)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	client := httputil.NewHTTPClient(30 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("upload request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
