# GoCrypt
**Work In Progress**
<p align="center">
<img src="https://github.com/td4b/GoCrypt/blob/master/image/key_icon-01.png" width="100" height="100">
</p>
<p align="center">
<img style="float: center;" src="https://goreportcard.com/badge/github.com/TD4B/GoCrypt">
<img style="float: center;" src="https://img.shields.io/badge/License-MIT-yellow.svg">

`When your afraid of forgetting a bunch of secrets.`<p>
`Keep one in you're head.`<p>
`With AES256 bit encryption, expedient they said!`
  
# Description
This is a CLI tool for encyrpting and uploading files to the IPFS blockchain network. The tool also generates unique file signatures and records them in a local database for one to keep track of the files. The reason behind this is due to the nature of AES256, at some point potentially a second block chain could be leveraged to track (file uniqueness). This would be the "Public" or non-confidential metadata assosciated with a file. The data will be nearly impossible to decrypt unless the client has the private key since the encryption method is symmetric (AES256).

# MicroServices Based Network

    +--------------------------+    +------+
    |Postgres (Metadata Service+----+      |
    +--------------------------+    |      |
    +--------------------------+    |      |
    |API                       +----+      |
    +--------------------------+    | RTR  |
    +--------------------------+    |      |
    |Encryption Service        +----+      |
    +--------------------------+    |      |
    +--------------------------+    |      |
    |Bulk Storage Node (Go|IPFS+----+      |
    +--------------------------+    +------+

# Done.
- Migrated IPFS Node to Container System.
- Migrated Postgres Services to Container Service.
- Updated golang routines to support Postgres Container Service.

# TODO.

@@ GoCrypt @@
- Migrate additional Services over to Containers.
- Need to implement multiple API's & MicroServices to suppor this.

@@ IPFS @@
- Configure Data Caching.
- Configure Private Node Network (SSL Verification). 

@@ Database @@
- Configure SSL Security.
- Add in Database security via Generated ENV userid & password variables when DB is spun up.
- Enable Database Replication & Configure Cluster Network.

# Planned / Additional Features
- Convert each service into a microservice!
- File handling / API. Would like to implement the ability to upload the data from the Client side.
- Client Side Mutual Authentication to allow connectivity to encryption/storage service. 
- API Interface to easily view storage metadata information and query attributes of specific files.
