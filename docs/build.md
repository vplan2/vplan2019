# How To Build

## Requirements

- [dep](https://github.com/golang/dep)  
  *Dependency management system for go*
- [git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)  
  *For cloning the repository and necessary for go get*
- [go compiler toolchain](https://golang.org/doc/install)  
  *goang compiler and other tools for go develoment*
- [GNU make](https://www.gnu.org/software/make/)  
  *Tool supporting building, release and testing*

Important to notice is, that scripts and the Makefile are created to be executed on Unix-like systems. If you are using Windows, cross-compile the binary by using WSL or by using [git-bash](https://gitforwindows.org/).

---

## Compiling

There are different options for compiling the binaries:

### make

Just use the `Makefile` in the root of the cloned repository:
```
$ make release
```

Now, the binary will be compiled as same as the frontend web files and they will be placed under a new generated directory named `./release`.

#### Cross compiling

By defining the envoirement variables `GOOS` and `GOARCH`, you can specify for whch system and architecture you want to compile the binaries, not depending on which system you are compiling on atually. [Here](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63) you can find a gist with all `GOOS` and `GOARCH` options you can use *(or use [this fork](https://gist.github.com/zekroTJA/4fa6c30606d9e702ccc3b84433dbf131) if the source is no more existent)*.

### Manual compilation

> This manual is just for compiling on a linux like system or shell.

1. Set up your go envoirement by setting up a gopath worktree and the `GOPATH` envoirement variable:  
   ```
   $ mkdir -p ~/go/src/github.com/zekroTJA
   $ cd ~/go/src/github.com/zekroTJA
   $ export GOPATH=$PWD
   ```

2. Cloning the repository into the worktree:  
   ```
   $ git clone https://github.com/zekroTJA/vplan2019.git .
   ```

3. Compiling frontend with zola:  
   ```
   $ cd ./web
   $ zola -c ../config/frontend.release.toml build
   $ cd ..

fqdn = "https://zekro.de:8080"
author = "industrieschule.de"
   ```

4. Before starting to build the backend files, you need to define some variables which will be passed to the binary via go compiler's ldflags:  
   ```
   $ LD="github.com/zekroTJA/vplan2019/internal/ldflags"
   $ TAG=$(git describe --tags)
   $ COMMIT=$(git rev-parse HEAD)
   $ GOVERS=$(go version | sed -e 's/ /_/g')
   ```

5. Now, compile the binary with following command:  
   ```
   $ go build -v -o ./release/vplan2_server -ldflags "\
       -X ${LD}.AppVersion=${TAG} \
       -X ${LD}.AppCommit=${COMMIT} \
       -X ${LD}.GoVersion=${GOVERS} \
       -X ${LD}.Release=TRUE" \
       ./cmd/server
   ```

6. Move the compiled web files to the `./release` directory:  
   ```
   $ mv -R ./web/public ./release/web
   ```

7. After that, you can try testing your binary by running `vplan2_server` in the release directory.

---

© 2019 Justin Trommler, Richard Heidenreich, Ringo Hoffmann  
Covered by MIT Licence.