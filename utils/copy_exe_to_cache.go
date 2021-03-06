package utils

import (
	"fmt"
	"io"
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/afero"
)

// CopyExeToCache copies over the nonmem executable to the cache from the run directory
func CopyExeToCache(
	fs afero.Fs,
	baseDir string,
	modelDir string,
	cacheDir string,
	nmNameInCache string,
) error {

	fullCacheDirPath := filepath.Join(baseDir, cacheDir)
	// this should always be there how the package is being used currently as the time
	// this function is called is after the modelDir is successfully created, but given
	// there may be other uses, to be safe will also check again.
	fullModelDirPath := filepath.Join(baseDir, modelDir)

	ok, err := DirExists(fullCacheDirPath, fs)
	if !ok || err != nil {
		//TODO: change these exits to instead just return an error probably
		log.Printf("issue with cache directory at: %s, will not save executable to cache. ERR: %s, ok: %v", cacheDir, err, ok)
		return err
	}
	//check that modelDir exists to copy nonmem executable into
	ok, err = DirExists(fullModelDirPath, fs)
	if !ok || err != nil {
		//TODO: change these exits to instead just return an error probably
		log.Printf("issue with model directory at: %s, will not save executable to cache. ERR: %s, ok: %v", cacheDir, err, ok)
		return err
	}
	// check nmNameInCache is in in cache
	// copy file from cache to modelDir

	newCacheFileLocation := filepath.Join(
		fullCacheDirPath,
		nmNameInCache,
	)
	cacheFile, err := fs.Create(newCacheFileLocation)
	if err != nil {
		return fmt.Errorf("error copying file: (%s)", err)
	}
	defer cacheFile.Close()

	var nonmemExecutableName string
	if runtime.GOOS == "windows" {
		nonmemExecutableName = "nonmem.exe"
	} else {
		nonmemExecutableName = "nonmem"
	}
	exeLocation := filepath.Join(
		fullModelDirPath,
		nonmemExecutableName,
	)
	exeFile, err := fs.Open(exeLocation)
	if err != nil {
		return fmt.Errorf("error with nonmem exe file: (%s)", err)
	}
	defer exeFile.Close()

	_, err = io.Copy(cacheFile, exeFile)
	if err != nil {
		return fmt.Errorf("error copying to new file: (%s)", err)
	}
	return nil
}
