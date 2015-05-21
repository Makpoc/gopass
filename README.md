# gopass
Inspired by [supergenpass-lib](https://github.com/chriszarate/supergenpass-lib) this is an attempt to make a non-complete password generator based on master password and domain.

Usage
=====
    $ go run main.go -master-file masterphrase.file -domain google.com -password-length 16
    Your password for google.com is: !@^#CZQGMGVHY3a9

    $ go run main.go -master super-mega-secret-master-phrase -domain github.com -password-length 16 -additional-info rev3
    Your password for github.com is: 06r68L1RMlyN)(*$
