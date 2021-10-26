package libs

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
)

func CheckAESKey(key_length int) bool {
	valid_key_lengths := []int{16, 24, 32}
	for i := 0; i < len(valid_key_lengths); i++ {
		if key_length == valid_key_lengths[i] {
			return true
		}
	}
	return false
}

func AESEncryptDecrypt(encrypt bool, filepath string, key string, chunk_size int64, outpath string, start_offset int64, end_offset int64) {

	infile, err := os.Open(filepath)
	CheckErrors(err)
	defer infile.Close()

	outfile, err := os.OpenFile(outpath, os.O_RDWR|os.O_CREATE, FilePerms)
	CheckErrors(err)
	defer outfile.Close()

	fs, err := os.Stat(filepath)
	CheckErrors(err)

	// The key should be 16 bytes (AES-128), 24 bytes (AES-192) or
	// 32 bytes (AES-256)
	key_bytes := []byte(key)
	key_length := len(key_bytes)
	valid_key := CheckAESKey(key_length)
	if !valid_key {
		fmt.Println("[!] The AES key should be 16 bytes (AES-128), 24 bytes (AES-192) or 32 bytes (AES-256) !")
		os.Exit(1)
	}

	var data_size int64 = 0
	var current_offset int64 = 0
	file_size := fs.Size()
	block, err := aes.NewCipher(key_bytes)
	CheckErrors(err)

	iv := make([]byte, block.BlockSize())
	iv_size := int64(len(iv))
	if encrypt {
		// Generate a new IV
		_, err := io.ReadFull(rand.Reader, iv)
		CheckErrors(err)
	} else {
		// Grab the current IV at the end of file
		data_size = file_size - iv_size
		_, err = infile.ReadAt(iv, data_size)
		CheckErrors(err)
		file_size = data_size
	}

	// The buffer size must be multiple of 16 bytes
	if chunk_size < 16 {
		chunk_size = 16
	} else {
		chunk_size = chunk_size - (chunk_size % 16)
	}

	stream := cipher.NewCTR(block, iv)
	buf := make([]byte, chunk_size)

	for {
		n, err := infile.Read(buf)
		if n > 0 {
			if !encrypt && current_offset == data_size {
				break
			}

			if current_offset+chunk_size > file_size {
				chunk_size = file_size - current_offset
			}

			bytes := make([]byte, chunk_size)
			if current_offset >= start_offset && current_offset <= end_offset {
				stream.XORKeyStream(bytes[:chunk_size], buf[:chunk_size])
			} else {
				bytes = buf[:chunk_size]
			}
			outfile.Write(bytes)
			current_offset += chunk_size
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Read %d bytes: %v", n, err)
			break
		}
	}
	if encrypt {
		// Append the IV at the end of the file
		outfile.Write(iv)
	}
}
