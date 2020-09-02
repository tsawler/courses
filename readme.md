[![Go Report Card](https://goreportcard.com/badge/github.com/tsawler/goblender-client-sample)](https://goreportcard.com/report/github.com/tsawler/goblender-client-sample)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/tsawler/goblender/master/LICENSE)
[![Version](https://img.shields.io/badge/goversion-1.14.x-blue.svg)](https://golang.org)
<a href="https://golang.org"><img src="https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat-square" alt="Built with GoLang"></a> 


# GoBlender client

Sample project for client specific code for GoBlender.

## Setup

First, install [goblender](https://github.com/tsawler/goblender).

All client specific code lives in `./client/clienthandlers`, and is its own 
git repository. When working in JetBrains Goland, you must create 
(or clone) a git repository in this location, and then add the directory
in Preferences -> Version Control.

## Updating on server
Change  `update.sh` in GoBlender root folder so as to enable git pull of client:

```shell script
# uncomment if using custom client code
#cd ./client/clienthandlers
#git pull
#cd ../..

# run migrations for pg
# soda migrate -c ../../migrations-pg/database.yml

#run client migrations for mariadb
# soda migrate -c ../../database.yml
```

After changing, it should look like this (assuming you want to run postgres migrations):

```shell script
# uncomment if using custom client code
cd ./client/clienthandlers
git pull
cd ../..

# run migrations for pg
soda migrate -c ../../migrations-pg/database.yml

#run client migrations for mariadb
# soda migrate -c ../../database.yml
```


## Using custom templates

Inside of `clientviews` there are two folders: `public` and `private`. If you wish to use the base templates
from goBlender to create templates, do it like this:

For public pages:

```
{{template "base" .}}

{{define "title"}}Some title{{end}}

{{define "body"}}
    <p>Put whatever you want here</a>
{{end}}
```

For admin pages:

```
{{template "admin-base" .}}

{{define "title"}}Some Title - vMaintain Admin{{end}}

{{define "admin-title"}}Some title{{end}}
{{define "content-title"}}Some title{{end}}

{{define "content"}}
    <p>Some content</p>
{{end}}
```

**Note**

You can override anything in the base templates, or specific pages/partials, but putting a file in 
`client/clientviews/public`, `client/clientviews/public/partials`, `client/clientviews/admin`, or 
`client/clientviews/admin/partials`.

## Client Specific Migrations

Migrations live in `client/migrations`. To run them, add the -c flag to soda, e.g.:

To generate Postgres migrations:
~~~
cd client/clienthandlers
soda -c ../../migrations-pg/database.yml generate fizz SomeMigrationName
~~~

To run Postgres migrations:
~~~
cd client/clienthandlers
soda -c ../../migrations-pg/database.yml migrate
~~~

To generate MariaDB/MySQL migrations:
~~~
cd client/clienthandlers
soda -c ../../database.yml generate fizz SomeMigration
~~~

To run MariaDB/MySQL migrations:
~~~
cd client/clienthandlers
soda -c ../../database.yml migrate
~~~

## Middleware

Add custom middleware to `./client/clienthandlers/client-middleware.go`, e.g.:

```go
// SomeMiddleware is sample middleware
func SomeMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ok := true
        // perform some logic to set ok
        
        if ok {
            next.ServeHTTP(w, r)
        } else {
         helpers.ClientError(w, http.StatusUnauthorized)
        }
    })
}
```