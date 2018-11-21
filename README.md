Gyotaku
=======
This repository contains the CTF problem used in SECCON Beginners NEXT 2017 Tokyo.

Official Writeup is sold at [BOOTH](https://dragonuniversity.booth.pm/items/1055860). (Japanese)

## Setup
```bash
git clone https://github.com/rekkusu/gyotaku.git
cd gyotaku
docker-compose build
docker-compose up -d
```

Then, open `http://[docker's host]:8080/`.

## Problem Statement
Steal admin's gyotakus.

Notice: ModSecurity and OWASP Core Rule Set are enabled.

