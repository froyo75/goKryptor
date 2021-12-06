package libs

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/lazybeaver/entropy"
)

const FilePerms = 0755
const DirPerms = 0755
const DefaultEncFileExt = ".enc"
const DefaultPlainFileExt = ".plain"

func CheckErrors(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func EncodeBase64(data []byte) string {
	b64encoded := base64.StdEncoding.EncodeToString(data)
	return b64encoded
}

func DecodeBase64(data string) []byte {
	b64decoded, err := base64.StdEncoding.DecodeString(data)
	CheckErrors(err)
	return b64decoded
}

func GetRealPath(path string) string {
	abspath, err := filepath.Abs(path)
	CheckErrors(err)
	return abspath
}

func CreateDirectory(dirpath string) {
	err := os.Mkdir(dirpath, DirPerms)
	CheckErrors(err)
}

func IsDirectory(filepath string) bool {
	fi, err := os.Stat(filepath)
	CheckErrors(err)
	return fi.IsDir()
}

func RemoveFile(filepath string) {
	err := os.Remove(filepath)
	CheckErrors(err)
}

func RenameFile(inpath string, outpath string) {
	err := os.Rename(inpath, outpath)
	CheckErrors(err)
}

func GenOutPath(encrypt bool, inpath string, outpath string, fileext string) string {
	outfilepath := ""

	if len(fileext) == 0 {
		if encrypt {
			fileext = DefaultEncFileExt
		} else {
			fileext = DefaultPlainFileExt
		}
	}

	if len(outpath) > 0 {
		absoutpath := GetRealPath(outpath)
		endpath := outpath[len(outpath)-1:]
		if endpath == "/" || endpath == "\\" {
			if !IsDirectory(absoutpath) {
				CreateDirectory(absoutpath)
			}
			inpath = absoutpath + "/" + path.Base(inpath)
		}
	}

	if !encrypt {
		findext := strings.HasSuffix(inpath, fileext)
		if findext {
			outfilepath = strings.TrimSuffix(inpath, fileext)
		} else {
			if fileext != DefaultPlainFileExt {
				fmt.Printf("[!] File extension '%s' not found for the specified file '%s', using the default file extension '%s' instead...\n", fileext, inpath, DefaultPlainFileExt)
			}
			fileext = DefaultPlainFileExt
			outfilepath = inpath + fileext
		}
	} else {
		outfilepath = inpath + fileext
	}

	return outfilepath
}

func CalculateOffset(filepath string, start_offset int64, end_offset int64) (int64, int64) {
	fs, err := os.Stat(filepath)
	CheckErrors(err)
	file_size := fs.Size()

	if start_offset < 0 {
		start_offset = file_size + start_offset
	}

	if end_offset < 0 {
		end_offset = file_size + end_offset
	}

	if start_offset > end_offset && end_offset != 0 {
		fmt.Println("[!] 'start-offset' must be less than 'end-offset' !")
		os.Exit(1)
	} else if start_offset > file_size || end_offset > file_size {
		fmt.Printf("[!] The maximum offset limit for '%s' is %d !\n", filepath, file_size)
		os.Exit(1)
	} else if (start_offset == 0 && end_offset == 0) || (start_offset > 0 && end_offset == 0) || (start_offset < 0 && end_offset == 0) {
		end_offset = file_size
	}

	return start_offset, end_offset
}

func CalculateAndDisplayEntropy(filepath string) {
	infile, err := os.Open(filepath)
	CheckErrors(err)
	defer infile.Close()

	reader := bufio.NewReader(infile)
	estimator := entropy.NewShannonEstimator()
	_, err = io.Copy(estimator, reader)
	CheckErrors(err)

	fmt.Printf("[*] Entropy of '%s' => %f\n", filepath, estimator.Value())
}

func Crypt(encrypt bool, algo string, key string, currentpath string, chunk_size int64, entropy bool, outpath string, fileext string, start_offset int64, end_offset int64) {

	action := ""
	if encrypt {
		action = "Encrypting"
	} else {
		action = "Decrypting"
	}

	err := filepath.Walk(currentpath, func(inpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if IsDirectory(inpath) {
			fmt.Printf("[+] Processing folder '%s'\n", inpath)
			return nil
		} else {
			var start int64 = 0
			var end int64 = 0	
			outfilepath := GenOutPath(encrypt, inpath, outpath, fileext)
			start, end = CalculateOffset(inpath, start_offset, end_offset)
			fmt.Printf("[+] '%s' file '%s' using the '%s' algorithm...\n", action, inpath, algo)

			switch algo {
			case "aes":
				AESEncryptDecrypt(encrypt, inpath, key, chunk_size, outfilepath, start, end)
				if entropy {
					CalculateAndDisplayEntropy(inpath)
					CalculateAndDisplayEntropy(outfilepath)
				}
				RemoveFile(inpath)
				fmt.Printf("[+] Done.\n")
				return nil
			case "xor":
				XOREncryptDecrypt(inpath, key, chunk_size, outfilepath, start, end)
				if entropy {
					CalculateAndDisplayEntropy(inpath)
					CalculateAndDisplayEntropy(outfilepath)
				}
				RemoveFile(inpath)
				fmt.Printf("[+] Done.\n")
				return nil
			case "none":
				outfilepath := GenOutPath(encrypt, inpath, outpath, fileext)
				fmt.Printf("[+] Renaming '%s' file using the '%s' file extension...\n", inpath, fileext)
				RenameFile(inpath, outfilepath)
				return nil
			default:
				return nil
			}
		}
	})
	if err != nil {
		CheckErrors(err)
	}
}
