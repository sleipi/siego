# Siego
-------

Siego is an http load tester and benchmarking utility (inspiered by [JoeDog/siege](https://github.com/JoeDog/siege)) 

## Features

* perform request in paralell
* support all Http Methods (`--method POST`)
* add request headers (`--header "Authorization: Bearer ***, X-Other-Header: Foo`)

## How-To

build it
````
go build
````

run it
````
./main --help
````

use it
````
./main --method POST --target http://127.0.0.1 -c 5
````