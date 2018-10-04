# GoCrypt
## Work In Progress :)
<p align="center">
<img style="float: center;" src="https://goreportcard.com/badge/github.com/TD4B/GoCrypt">
<img style="float: center;" src="https://img.shields.io/badge/License-MIT-yellow.svg">

`When your afraid of forgetting a bunch of secrets.`<p>
`Keep one in you're head.`<p>
`With AES256 bit encryption, expedient they said!`

# Description
This is a CLI tool for encyrpting and uploading files to the IPFS blockchain network. The tool also generates unique file signatures and records them in a local database for one to keep track of the files. Because of the way AES256 encyrption is implemented a second chain is needed to keep track and record metadata about the content/information of a file. This would be the "Public" or non-confidential metadata assosciated with a file so a client can identify it's validaity before attempting to decrypt the true file structure. Of course the data will be nearly impossible to decrypt unless the client has the private key (AES 256).

# Planned / Additional Features
- Better error handling.
- File handling.
- Intraplanary FIlesystem Interface (IPFS) to upload and store data securley.
- CLI / Management Interface to keep track of you're files. 
- API Interface to easily view storage metadata information and query attributes of specific files.


