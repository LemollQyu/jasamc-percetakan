package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type LocalStorage struct {
	BaseDir string
	Url     string
}

func NewLocalStorage(baseDir, url string) *LocalStorage {
	return &LocalStorage{
		BaseDir: baseDir,
		Url:     url,
	}
}

// helper untuk validasi ekstensi
func isValidExt(filename string, allowedExts []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, a := range allowedExts {
		if ext == a {
			return true
		}
	}
	return false
}

// helper generic untuk upload file
func (s *LocalStorage) uploadFile(ctx context.Context, file *multipart.FileHeader, folder string, allowedExts []string, maxSize int64) (string, error) {
	if file.Size > maxSize {
		return "", fmt.Errorf("ukuran file terlalu besar")
	}

	if !isValidExt(file.Filename, allowedExts) {
		return "", fmt.Errorf("ekstensi file tidak diperbolehkan")
	}

	dir := filepath.Join(s.BaseDir, folder)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))
	fullPath := filepath.Join(dir, filename)

	c := ctx.(*gin.Context)
	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		return "", err
	}

	publicURL := s.Url + "/static/" + folder + "/" + filename
	return publicURL, nil
}

// Ikon Layanan (Limit lebih ketat untuk keamanan dan efisiensi)
func (s *LocalStorage) UploadServiceIcon(ctx context.Context, file *multipart.FileHeader) (string, error) {

	const maxIconSize = 500 * 1024
	return s.uploadFile(ctx, file, "jasa/icon-service", []string{".svg", ".png"}, maxIconSize)
}

// Upload thumbnail
func (s *LocalStorage) UploadServiceThumbnail(ctx context.Context, file *multipart.FileHeader) (string, error) {

	const maxThumbnailSize = 10 * 1024 * 1024
	return s.uploadFile(ctx, file, "jasa/thumbnail-service", []string{".jpg", ".jpeg", ".png", ".webp"}, maxThumbnailSize)
}

// Upload gallery
func (s *LocalStorage) UploadServiceGallery(ctx context.Context, file *multipart.FileHeader) (string, error) {

	const maxGallerySize = 10 * 1024 * 1024
	return s.uploadFile(ctx, file, "jasa/gallery-service", []string{".jpg", ".jpeg", ".png", ".webp"}, maxGallerySize)
}

// delete storage yang dituju
func (s *LocalStorage) DeleteMediaByURL(
	ctx context.Context,
	mediaURL string,
) error {

	fmt.Println("==== DELETE MEDIA DEBUG ====")
	fmt.Println("MEDIA URL :", mediaURL)
	fmt.Println("BASE DIR  :", s.BaseDir)

	// 1️⃣ Pastikan URL mengandung /static/
	idx := strings.Index(mediaURL, "/static/")
	if idx == -1 {
		return fmt.Errorf("invalid media url")
	}

	// 2️⃣ Ambil path relatif setelah /static/
	relativePath := mediaURL[idx+len("/static/"):]
	// contoh: jasa/gallery-service/123.webp

	if relativePath == "" {
		return fmt.Errorf("empty media path")
	}

	// 3️⃣ Gabungkan dengan BaseDir
	fullPath := filepath.Join(s.BaseDir, relativePath)

	// 4️⃣ Security check (hindari path traversal)
	cleanBase := filepath.Clean(s.BaseDir)
	cleanPath := filepath.Clean(fullPath)

	fmt.Println("CLEAN BASE:", cleanBase)
	fmt.Println("CLEAN PATH:", cleanPath)

	if !strings.HasPrefix(cleanPath, cleanBase) {
		return fmt.Errorf("invalid file path")
	}

	// 5️⃣ Cek file ada atau tidak
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		// file sudah tidak ada → anggap sukses
		return nil
	}

	// 6️⃣ Hapus file
	if err := os.Remove(cleanPath); err != nil {
		return err
	}

	return nil
}

// Upload category icon
func (s *LocalStorage) UploadCategoryIcon(ctx context.Context, file *multipart.FileHeader) (string, error) {

	if file.Size > 2*1024*1024 {
		return "", fmt.Errorf("ukuran icon maksimal 2MB")
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".png" && ext != ".svg" {
		return "", fmt.Errorf("icon hanya png dan svg")
	}

	dir := filepath.Join(s.BaseDir, "jasa", "icon-category")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	fullPath := filepath.Join(dir, filename)

	// ctx dari gin
	c := ctx.(*gin.Context)
	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		return "", err
	}

	// public URL
	return s.Url + "/static/jasa/icon-category/" + filename, nil
}
