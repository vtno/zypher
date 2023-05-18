# Zypher

A dead simple AES encryption / decryption CLI and Go library.

It works perfectly when you just want to encrypt `.env` file and check it into the codebase like [Rails](https://edgeguides.rubyonrails.org/security.html#environmental-security) does.

## Installation

- requires go 1.20.4

```shell
go install github.com/vtno/zypher/cmd/zypher@latest
```

## Usage

There are 3 subcommands. `encrypt`, `decrypt` which have similar options and `keygen` which is a utility command to generate a `zypher.key` file.

```shell
# available options:
# -k, --key       A key to be used on encryption / decryption. The length should be:
#                  16 bytes for AES-128
#                  24 bytes for AES-192
#                  32 bytes for AES-256
#
# -f, --file      A path to file to be encrypted / decrypted
# -o, --out       A path to output file
# -kf, --key-file A path to key file. Default: zypher.key

# encrypting/decrypting from arg to stdout
zypher encrypt -k <AES-KEY> input-to-be-encrypt
zypher decrypt -k <AES-KEY> base64-of-the-encrypted-text

# encrypting/decrypting from file input to stdout
zypher encrypt -k <AES-KEY> -f input.txt
zypher decrypt -k <AES-KEY> -f input.txt.enc

# encrypting/decrypting from file input to another file
zypher encrypt -k <AES-KEY> -f input.txt -o input.txt.enc
zypher decrypt -k <AES-KEY> -f input.txt.enc -o input.txt

# this works too!
zypher encrypt -k <AES-KEY> -f input.txt > input.txt.enc
zypher decrypt -k <AES-KEY> -f input.txt.enc > input.txt

# the key can also be set as ZYPHER_KEY env
export ZYPHER_KEY=<AES-KEY>
zypher encrypt -f input.txt > input.txt.enc
zypher decrypt -f input.txt.enc > input.txt

# the key is automatically lookup on zypher.key file
# and could be overridden with --key-file flag
zypher encrypt -kf your-own.key -f input.txt -o input.txt.enc
zypher decrypt -kf /some/path/your-own.key -f input.txt.enc -o input.txt

# generate zypher.key easily with keygen command
zypher keygen
```
