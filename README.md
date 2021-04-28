# stegasaurus
Steganography in Go. Inspired heavily from [Steganography](https://github.com/auyer/steganography) by [auyer](https://github.com/auyer).
Stegasaurus encrypts any file with AES 256 and embeds the encrypted file into a PNG, as well as extracts and decrypts. Stegasaurus can also simply encrypt
and decrypt a file.

### Requirements
Tested with `go1.15.6`

### Installing
1. `git clone https://github.com/markditsworth/stegasaurus`
2. `cd stegasaurus && make build`
3. `mv stegasarus /path/to/directory/in/$PATH`

### Usage
Encrypt and embed a file: `stegasaurus --embed -i <file_to_encrypt> -p <passkey_to_use> --image <PNG_file_to_use> -o <output_filename_for_encoded_image>`

Decrypt from an encoded image: `stegasaurus --extract --image <encoded_image> -p <passkey> -o <output_filename>`

Encrypt a file: `stegasaurus --encrypt -i <file_to_encrypt> -o <output_filename> -p <passkey_to_use>`

Decrypt a file:  `stegasaurus --decrypt -i <file_to_decrypt> -o <output_filename> -p <passkey>`

### Coming Soon
- Analyzing images to look for encoded information
- JPEG support
- Tests
