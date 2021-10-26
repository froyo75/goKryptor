package libs

import (
	"fmt"
	"io"
	"os"
)

func XOREncryptDecrypt(filepath string, key string, chunk_size int64, outpath string, start_offset int64, end_offset int64) {

	infile, err := os.Open(filepath)
	CheckErrors(err)
	defer infile.Close()

	outfile, err := os.OpenFile(outpath, os.O_RDWR|os.O_CREATE, FilePerms)
	CheckErrors(err)
	defer outfile.Close()

	fs, err := os.Stat(filepath)
	CheckErrors(err)

	var current_offset int64 = 0
	var current_chunk int64 = 0
	file_size := fs.Size()
	total_chunk := file_size / chunk_size
	max_chunk_size := total_chunk * chunk_size
	remainder_bytes := file_size % chunk_size
	buf := make([]byte, chunk_size)
	key_length := len(key)

	if chunk_size > file_size {
		chunk_size = file_size
	}

	for {
		n, err := infile.Read(buf)
		if n > 0 {
			if current_chunk < max_chunk_size {
				current_chunk += chunk_size
			} else if remainder_bytes > 0 {
				chunk_size = remainder_bytes
			}

			bytes := make([]byte, chunk_size)
			for i := range bytes {
				if current_offset >= start_offset && current_offset <= end_offset {
					bytes[i] = buf[i] ^ key[i%key_length]
				} else {
					bytes[i] = buf[i]
				}
				current_offset++
			}
			outfile.Write(bytes)
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("[!] Read %d bytes: %v\n", n, err)
			break
		}
	}
}
