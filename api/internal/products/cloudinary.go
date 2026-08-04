package products

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const defaultCloudinaryFolder = "oyecouture"

// CloudinaryUploader sends admin product images to Cloudinary with signed uploads.
type CloudinaryUploader struct {
	cloudName  string
	apiKey     string
	apiSecret  string
	folder     string
	httpClient *http.Client
}

// NewCloudinaryUploader creates a signed Cloudinary image uploader.
func NewCloudinaryUploader(cloudName string, apiKey string, apiSecret string, folder string) *CloudinaryUploader {
	return &CloudinaryUploader{
		cloudName: strings.TrimSpace(cloudName),
		apiKey:    strings.TrimSpace(apiKey),
		apiSecret: strings.TrimSpace(apiSecret),
		folder:    strings.Trim(strings.TrimSpace(valueOrDefault(folder, defaultCloudinaryFolder)), "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IsConfigured reports whether uploads can be attempted.
func (uploader *CloudinaryUploader) IsConfigured() bool {
	return uploader != nil &&
		uploader.cloudName != "" &&
		uploader.apiKey != "" &&
		uploader.apiSecret != ""
}

// Upload sends one image to Cloudinary and returns the resulting public_id.
func (uploader *CloudinaryUploader) Upload(
	requestContext context.Context,
	productSlug string,
	upload AdminImageUpload,
) (string, error) {
	if !uploader.IsConfigured() {
		return "", ErrImageUploadUnavailable
	}

	publicID := uploader.publicID(productSlug, upload.FileName)
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	parameters := map[string]string{
		"public_id": publicID,
		"timestamp": timestamp,
	}

	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)
	if err := writeCloudinaryFields(multipartWriter, parameters, uploader); err != nil {
		return "", err
	}
	fileWriter, err := multipartWriter.CreateFormFile("file", upload.FileName)
	if err != nil {
		return "", fmt.Errorf("create cloudinary upload file: %w", err)
	}
	if _, err := io.Copy(fileWriter, bytes.NewReader(upload.Content)); err != nil {
		return "", fmt.Errorf("write cloudinary upload file: %w", err)
	}
	if err := multipartWriter.Close(); err != nil {
		return "", fmt.Errorf("close cloudinary upload body: %w", err)
	}

	request, err := http.NewRequestWithContext(
		requestContext,
		http.MethodPost,
		fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", uploader.cloudName),
		&requestBody,
	)
	if err != nil {
		return "", fmt.Errorf("build cloudinary upload request: %w", err)
	}
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	response, err := uploader.httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("upload image to cloudinary: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		responseBody, _ := io.ReadAll(io.LimitReader(response.Body, 2048))
		return "", fmt.Errorf("cloudinary upload failed: status %d: %s", response.StatusCode, responseBody)
	}

	var uploadResponse struct {
		PublicID string `json:"public_id"`
	}
	if err := json.NewDecoder(response.Body).Decode(&uploadResponse); err != nil {
		return "", fmt.Errorf("decode cloudinary upload response: %w", err)
	}
	if strings.TrimSpace(uploadResponse.PublicID) == "" {
		return "", fmt.Errorf("cloudinary upload response missing public_id")
	}

	return uploadResponse.PublicID, nil
}

func (uploader *CloudinaryUploader) publicID(productSlug string, fileName string) string {
	extension := strings.TrimPrefix(strings.ToLower(filepath.Ext(fileName)), ".")
	if extension == "" {
		extension = "image"
	}

	return fmt.Sprintf("%s/%s-%d-%s", uploader.folder, productSlug, time.Now().UnixNano(), extension)
}

func writeCloudinaryFields(
	multipartWriter *multipart.Writer,
	parameters map[string]string,
	uploader *CloudinaryUploader,
) error {
	for fieldName, fieldValue := range parameters {
		if err := multipartWriter.WriteField(fieldName, fieldValue); err != nil {
			return fmt.Errorf("write cloudinary field %s: %w", fieldName, err)
		}
	}
	if err := multipartWriter.WriteField("api_key", uploader.apiKey); err != nil {
		return fmt.Errorf("write cloudinary api key: %w", err)
	}
	if err := multipartWriter.WriteField("signature", cloudinarySignature(parameters, uploader.apiSecret)); err != nil {
		return fmt.Errorf("write cloudinary signature: %w", err)
	}

	return nil
}

func cloudinarySignature(parameters map[string]string, apiSecret string) string {
	parameterNames := make([]string, 0, len(parameters))
	for parameterName := range parameters {
		parameterNames = append(parameterNames, parameterName)
	}
	sort.Strings(parameterNames)

	signatureParts := make([]string, 0, len(parameterNames))
	for _, parameterName := range parameterNames {
		signatureParts = append(signatureParts, parameterName+"="+parameters[parameterName])
	}

	hash := sha1.Sum([]byte(strings.Join(signatureParts, "&") + apiSecret))
	return hex.EncodeToString(hash[:])
}

func valueOrDefault(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}

	return value
}
