# pss

Simple but secure CLI password manager written in Go.

## Building

* Install latest version of Go from [golang.org](https://golang.org/).
* Make sure you have $GOPATH setup.
* Add $GOPATH/bin to your $PATH.
* `go get github.com/yamnikov-oleg/pss`
* PROFIT

## Usage

Before using pss you must initialize password storage with `pss init`.

It will promt for *master password*, which will be used to encrypt your storage. Please, make a strong master password, but remember it! If you lose it, you will *lose all your passwords*.

```
$ pss init
Enter master-password:
Repeat:
Successfully created new password storage at /home/oleg/.pss_storage/storage
```

You can add new password record to the storage with `pss insert` or `pss i`:

```
$ pss i github.com yamnikov-oleg
Master-password:
Password for yamnikov-oleg at github.com:
Repeat:
Successfully saved the record
```

To view all the passwords you have ever saved, type `pss list` or `pss ls`:

```
$ pss ls
Master-password:
Contents of the storage:
           yamnikov-oleg @ github.com
                 workjob @ gmail.com
       workjob@gmail.com @ twitter.com
```

To retrieve a password and use it for logging in or something else, type `pss retrieve` or `pss r`:

```
$ pss r twitter
Master-password:
Password for:
       workjob@gmail.com @ twitter.com
has been copied to your clipboard.
Use Ctrl-V or 'Paste' command to use it.
```

Notice that the password is copied directly into the clipboard and not showed on the screen. Therefore, your account won't be compromised by your terminal history.

You can always view all the available commands by calling `pss` without arguments:

```
$ pss
Usage: pss <command>

Available commands:

dump        Dump storage contents into json file.
            Usage: dump <path>

gen         Generate random password of given length (default 12).
            Alises: g
            Usage: gen [length]

init        Make a new password storage.
            Usage: init

insert      Put a new password record into the storage.
            Alises: ins, i
            Usage: insert <website> <username>

list        Print list of saved records, filtred by optional search query.
            Alises: ls, l
            Usage: list [search query]

load        Load dump from json file into the storage.
            Usage: load <path>

modify      Modify a record in the storage.
            Alises: mod, m, update, upd, u
            Usage: modify <website> [username]

remove      Remove specific record from storage.
            Alises: rm
            Usage: remove <website> [username]

retrieve    Load a password from storage to clipboard
            Alises: r, checkout, co
            Usage: retrieve <website> [username]

```

## Encryption

You password storage is encrypted using [AES-256](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) cryptography argorithm with [cipher feedback mode](https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#Cipher_Feedback_.28CFB.29). Master password is also hashed with [SHA256](https://en.wikipedia.org/wiki/SHA-2) before using as a cipher key.

## Security

You can pass your password storage file whereever you want, but with caution. While it's not practically possible to decrypt it without knowing the *encryption key*, *master password* is usually much shorter and predictable, so it would be enough to bruteforce all the short passwords than 32-byte keys to break your storage.

But with a strong master password it's much safer.

If an attacker has bruteforce rate of one million combinations per second, he would crack the 10-symbol password in 18 thousand years.
12-symbol - in 74 million years.
20-symbol - in 10^22 years.
