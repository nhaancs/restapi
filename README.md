# RESTAPI

## Getting started
- Setup environment: `make envup`
- Run migrations: `make migrate`
- Run the service: `make run`
- Check other commands in the `Makefile` file

## Project structure
```
|middleware: application middlewares
|migration: database migrations
|pkg: helper packages (like 3rd-party packages in vendor/)
|module
|--product
   |--producttransport: receive and respond client requests 
   |--productbusiness: handle the business logic
   |--productstore: handle database logic
   |--productmodel: contain models in the module
|--user
    |--usertransport
    |--userbusiness
    |--userstore
    |--usermodel
```
![project-structure](./project-structure.jpg)