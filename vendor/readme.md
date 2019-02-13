### `vendor`

Dieser Ordner enthält sämtliche Dateien und Repositories von 3rd-Party-Packages. In diesem Fall werden diese in Form von Git-Submodules angelegt und verwaltet. Somit wird bei `go get` zuerst in dem `vendor` Verzeichnis nachgeschaut, bevor die repositories online bezogen werden, was das Compiling auch offline-tauglich machen würde.

Beim clonen der Repository muss dabei beachtet werden, dass nur mit dem Flag `--recursive` alle Submodules mit gecloned wirden:
```
$ git clone git@github.com:zekroTJA/vplan2019.git --recursive
```

Das clonen der Submodules ist auch im Nachhinein möglich mit:
```
$ git submodule init
$ git submodule update
```

Möchte man ein Submodule hinzufügen, so wird folgender Command verwendet:
```
$ git submodule add github.com/gorilla/mux vendor/github.com/gorilla/mux
```
