
# goKryptor

***`goKryptor` is a small and portable cryptographic tool for encrypting and decrypting files. This tool supports XOR and AES-CTR (Advanced Encryption Standard in Counter Mode) encryption algorithm. The '--start-offset' or '--end-offset' flag allows you to encrypt/decrypt a portion of file or the whole of it.***

```shell
Usage: goKryptor -a aes -k STRING -c 1024 <command> <file or folder>

A tool to encrypt/decrypt files by @froyo75

Flags:
  -h, --help               Show context-sensitive help.
  -a, --algo=STRING        Algorithm used to encrypt or decrypt data (xor,aes). The none value simulate encryption/decryption by simply renaming the files with the file extension specified.
  -k, --key=STRING         Key used to encrypt or decrypt data.
  -c, --chunk-size=1024    Chunk size used to encrypt or decrypt data (default: 1024 bytes).
  -o, --out-path=STRING    Specify the output path to write encrypted data. A specific file can be specified or a folder path (must end with '/' or '\\').
      --start-offset=0     Encrypt/Decrypt data FROM a specific offset (by using negative number, the encryption/decryption starts from the end of file).
      --end-offset=0       Encrypt/Decrypt data TO a specific offset (by using negative number, the encryption/decryption starts from the end of file).
  -x, --file-ext=STRING    Specify the file extension to be used for encrypting/decrypting files (default: .enc or .plain).
  -e, --entropy            Display and calculate the Shannon Entropy of the file.
  -v, --version            Print version information and quit.

Commands:
  encrypt --key=STRING --chunk-size=1024 <path>
    Encrypt file or all items in folder.

  decrypt --key=STRING --chunk-size=1024 <path>
    Decrypt file or all items in folder.

Run "goKryptor <command> --help" for more information on a command.
```
