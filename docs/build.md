# How To Build

## Requirements

- [git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [go compiler toolchain](https://golang.org/doc/install)
- [gcc](https://gcc.gnu.org/)
- [dep](https://github.com/golang/dep)

If you are on windows, you can install the [TDM-GCC](http://tdm-gcc.tdragon.net/download) toolchain, which includes gcc and make. Also, you should install [gitbash](https://gitforwindows.org/) and use it for executing bash scripts and the Makefile.


## Compiling

There are different options for compiling the binaries:


### make

Just use the `Makefile` in the root of the cloned repository:
```
$ make
```

Then, you will have the compiled binary file in the root directory of the repository.


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