# SyncthingUpgradeServer
## use:

mkdir dl

copy sha256sum.txt.asc sha1sum.txt.asc syncthing-windows-amd64-v1.26.0-rc.2.zip syncthing-linux-amd64-v1.26.0-rc.2.tar.gz  ...  into "dl"

run


## params:
  - string
        file dir (default "./dl")
  -listen string
        http listen (default "0.0.0.0:8080")
  -url string
        download url (default "http://127.0.0.1:8080")
