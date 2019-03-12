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


## Compiling

There are different options for compiling the binaries:


### make

Just use the `Makefile` in the root of the cloned repository:
```
$ make release
```

Now, the binary will be compiled as same as the frontend web files and they will be placed under a new generated directory named `./release`.


### Bash script

In order to use the bash script, you need to prepare your envoirement for compiling:

1. Create the GOPATH directory three like following:  
   ```
   $ mkdir -p ./gopath/src/github.com/zekroTJA/vplan2019
   ```
2. Copy all files of the repository into the created work directory:  
   ```
   $ cp -R ./* ./gopath/src/github.com/zekroTJA/vplan2019
   ```
3. Now, set the `gopath` directory as `GOPATH` envoirement variable:  
   ```
   $ export GOPATH=$PWD/gopath
   ```
4. Now, enter your work tree and execute the build script in `scripts/build.sh`:  
   ```
   $ cd ./gopath/src/github.com/zekroTJA/vplan2019
   $ bash ./scripts/build.sh
   ```

Now, the generated binary is located in the created `./bin` directory.


### Instalaltion

To install the server application, you need to move the generated binary **and** the `web` folder from the repository to a location where you want. It is very important that the `web` folder is directly in the same directory as the binary. Else, you need to specify the `web` folders location with the `-web` start flag.

If the application is started first time, a preset config file will be generated with the name `config.yml`. You can also specify an other config location with the `c` flag.