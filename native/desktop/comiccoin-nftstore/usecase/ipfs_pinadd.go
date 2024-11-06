package usecase

import (
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
	pkg_domain "github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type IPFSPinAddUseCase struct {
	logger   *slog.Logger
	ipfsRepo pkg_domain.IPFSRepository
}

func NewIPFSPinAddUseCase(logger *slog.Logger, r1 pkg_domain.IPFSRepository) *IPFSPinAddUseCase {
	return &IPFSPinAddUseCase{logger, r1}
}

func (uc *IPFSPinAddUseCase) Execute(multipartFile multipart.File) (string, error) {
	//
	// STEP 1:
	// Validation.
	//

	e := make(map[string]string)

	if multipartFile == nil {
		e["multipartFile"] = "missing value"
	}
	if len(e) != 0 {
		return "", httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2:
	// Convert from `multipart.File` to `[]byte`.
	//

	tmpFilepath := os.TempDir() + "/" + "test.txt"

	file, err := saveMultipartFileToDisk(multipartFile, tmpFilepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	//
	// STEP 3:
	// Execute submitting to IPFS.
	//

	cid, err := uc.ipfsRepo.AddViaFile(file, true)
	if err != nil {
		return "", err
	}

	uc.logger.Debug("data",
		slog.Any("cid", cid),
		slog.Any("tmpFilepath", tmpFilepath))

	return cid, nil
}

// Convert multipart.File to *os.File
func saveMultipartFileToDisk(src multipart.File, destPath string) (*os.File, error) {
	// Create a new file on the filesystem
	destFile, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	// Copy the contents of the multipart file to the destination file
	if _, err := io.Copy(destFile, src); err != nil {
		destFile.Close() // Make sure to close the file on error
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// Close the multipart file (source file)
	if err := src.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart file: %w", err)
	}

	// Return the newly created file on the filesystem
	return destFile, nil
}
