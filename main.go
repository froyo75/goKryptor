package main

import (
	"fmt"
	"goKryptor/libs"
	"time"

	"github.com/alecthomas/kong"
)

type Context struct {
	Debug bool
}
type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(vars["version"])
	app.Exit(0)
	return nil
}

var cli struct {
	Algo        string      `flag:"algo" short:"a" enum:"xor,aes,none" help:"Algorithm used to encrypt or decrypt data (xor,aes). The none value simulate encryption/decryption by simply renaming the files with the file extension specified."`
	Key         string      `flag:"key" short:"k" required:"" help:"Key used to encrypt or decrypt data."`
	ChunkSize   int64       `flag:"chunk_size" short:"c" required:"" default:"1024" help:"Chunk size used to encrypt or decrypt data (default: 1024 bytes)."`
	OutPath     string      `optional:"" short:"o" help:"Specify the output path to write encrypted data. A specific file can be specified or a folder path (must end with '/' or '\\')."`
	StartOffset int64       `optional:"" default:"0" help:"Encrypt/Decrypt data FROM a specific offset (by using negative number, the encryption/decryption starts from the end of file)."`
	EndOffset   int64       `optional:"" default:"0" help:"Encrypt/Decrypt data TO a specific offset (by using negative number, the encryption/decryption starts from the end of file)."`
	FileExt     string      `optional:"" short:"x" help:"Specify the file extension to be used for encrypting/decrypting files (default: .enc or .plain)."`
	Entropy     bool        `optional:"" short:"e" help:"Display and calculate the Shannon Entropy of the file."`
	Version     VersionFlag `name:"version" short:"v" help:"Print version information and quit."`

	Encrypt struct {
		Path string `arg:"" help:"Encrypt file or all items in folder." type:"path"`
	} `cmd:"" help:"Encrypt file or all items in folder."`

	Decrypt struct {
		Path string `arg:"" help:"Decrypt file or all items in folder." type:"path"`
	} `cmd:"" help:"Decrypt file or all items in folder."`
}

func main() {
	start := time.Now()

	ctx := kong.Parse(&cli, kong.Name("goKryptor"), kong.Description("A tool to encrypt/decrypt files by @froyo75"), kong.UsageOnError(), kong.Vars{"version": "1.0"})
	switch ctx.Command() {
	case "encrypt <path>":
		libs.Crypt(true, cli.Algo, cli.Key, cli.Encrypt.Path, cli.ChunkSize, cli.Entropy, cli.OutPath, cli.FileExt, cli.StartOffset, cli.EndOffset)
	case "decrypt <path>":
		libs.Crypt(false, cli.Algo, cli.Key, cli.Decrypt.Path, cli.ChunkSize, cli.Entropy, cli.OutPath, cli.FileExt, cli.StartOffset, cli.EndOffset)
	}

	elapsed := time.Since(start)
	fmt.Printf("Execution took %s\n", elapsed)
}
