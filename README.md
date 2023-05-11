# Zypher

A dead simple AES encryption / decryption CLI and Go library.

It works perfectly when you just want to encrypt `.env` file and check it into the codebase like [Rails](https://edgeguides.rubyonrails.org/security.html#environmental-security) does.

## Usage

There are only 2 subcommands `encrypt` and `decrypt`, both have similar options.

```shell
# available options:
# -k, --key    A key to be used on encryption / decryption. The length should be:
#                16 bytes for AES-128
#                24 bytes for AES-192
#                32 bytes for AES-256
#
# -f, --file  A path to file to be encrypted / decrypted
# -o, --out   A path to output file

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
```
