# silly-syntax-backend-bonanza

- [Instalation and building](#instalation)
  - [frontend](#frontend---admin-dashboard)
  - [backend](#frontend---admin-dashboard)
- [Core](#core)
  - [Overview](#overview)
  - [Key Features](#key-features)
    - [Automatic Route Generation](#automatic-route-generation)
    - [GUI for Collections](#key-features)
    - [Views System](#views-system)
    - [Permission System](#permission-system)
- [Core - How to use it](#core---how-to-use-it)
  - [Authentication](#authentication)
  - [Collections Routes](#collections-routes)
    - [Create](#create)
    - [Read](#read)
    - [Update](#update)
    - [Delete](#delete)
- [Plugins](#plugins)
  - [Description](#description-1)
  - [Key features](#key-features-1)
  - [How to Create a Plugin](#key-features)

## Instalation

### Frontend - admin dashboard

Open admin dashboard directory `cd admin-dashboard/`

Install dependencies `npm install`

Running and testing `npm run dev`

Building `npm run build` this should create directory **./dist**

### Backend - golang core

Open backend / core directory `cd backend/`

Running and testing `go run .`

Building to create minimal build without **admin dashboard** skip to step 2

1. Before this step make sure to finish steps in [frontend instalation section](#frontend---admin-dashboard). Run `statik -src="../admin-dashboard/dist"` this should generate statik folder and file
1. If you want to make minimal build comment line in main.go `_ "github.com/xomatix/silly-syntax-backend-bonanza/statik"` and delete statik directory
1. After that use golang build command `go build -o silly-syntax-backend-bonanza.exe .`

## Core

### Overview

`silly-syntax-backend-bonanza` Backend-as-a-Service (BaaS) System is designed to simplify backend development by automatically generating CRUD (Create, Read, Update, Delete) routes for your applications. With a focus on ease of use and flexibility, this system empowers developers to build robust APIs effortlessly.

### Key Features

#### Automatic Route Generation

The system automatically creates RESTful routes for all collections, allowing you to focus on your application's logic rather than boilerplate code. This includes endpoints for creating, reading, updating, and deleting records.

#### GUI for Collections

The intuitive graphical user interface (GUI) allows users to create and manage tables, referred to as collections. Users can easily define fields, data types, and relationships between collections without writing any backend code or logic.

#### Views System

Our views system enables users to generate SQL-based reports effortlessly. You can create customized views that aggregate and present data from multiple collections, allowing for insightful reporting and analysis.

#### Permission System

The built-in permission system utilizes macros and user context to ensure secure access control. Administrators can define permissions at a granular level, controlling which users can access or modify specific collections and views based on their roles and contexts.

#### Conclusion

This BaaS system provides a comprehensive solution for backend development, combining automatic route generation, an intuitive GUI, powerful reporting capabilities, and a robust permission system. It's designed to help developers focus on building features while ensuring security and flexibility.

## Core - how to use it

### Authentication

Login route `{{base_url}}/api/auth/login`

Send **POST** HTTP request with body

Basic user "admin"

Login `admin` Password `admin`

```
{
    "login": "string", // required
    "password": "string" // required
}
```

Response

```
{
    "message": "string",
    "success": bool,
    "data": "token"
}
```

To all of the collection routes needs authenticated user

Just that set header `Authorization` with token from login

### Collections Routes

### Create

route `{{base_url}}/api/collection/update`

Send **POST** HTTP request with body

```
{
    "collectionName": "string", // required
    "Values":{
        "Key": "value" // value must be same datatype as configured column
    } // required
}
```

Response

```
{
    "message": "string",
    "success": bool,
    "data": 0 // id of created record
}
```

### Read

route `{{base_url}}/api/collection/list`

Send **POST** HTTP request with body

```
{
    "collectionName": "string", // required
    "filter": string,
    "size": number,
    "limit": number,
    "ID": number[] // if provided all filters will be ommited and will return ony selected ids
}
```

Response

```
{
    "message": "string",
    "success": bool,
    "data": [] // list of found in database objects
}
```

### Update

route `{{base_url}}/api/collection/update`

Send **POST** HTTP request with body

```
{
    "collectionName": "string", // required
    "ID":0, // required
    "Values":{
        "Key": "value" // value must be same datatype as configured column
    } // required
}
```

Response

```
{
    "message": "string",
    "success": bool,
    "data": ""
}
```

### Delete

path `{{base_url}}/api/collection/delete`

Send **POST** HTTP request with body

```
{
    "collectionName": "string", // required
    "ID":0 // required
}
```

Response

```
{
    "message": "string",
    "success": bool,
    "data": ""
}
```

## Plugins

### Description

This project uses a plugin system to dynamically extend the functionality of the application without modifying the core codebase. The plugins are compiled as shared objects (.so files) and loaded at runtime.

Each plugin follows a standard structure and implements an InitPlugin function. This function initializes the plugin with required resources (such as HTTP routes, or utility functions) and returns a set of operations or handlers that the core system can invoke.

### Key Features

- **Dynamic Loading**: Plugins are loaded at runtime from the ./plugins directory, allowing seamless extensibility.
- **Extensible API**: Plugins can register new API routes or add business logic without modifying the main application.
- **Isolation**: Each plugin is isolated in its own module, minimizing interference with the core logic.
- **Easy Integration**: A flexible interface that allows plugins to access shared resources like HTTP routers, database query functions, etc.
- **Triggers**: Plugins can interfere HTTP functions like block update or run some calculations

### How to Create a Plugin

Implement an InitPlugin(context interface{}) function.
Register routes or operations using the provided context (e.g., HTTP router or database connections).
Compile the plugin as a shared object using go build -buildmode=plugin -o example_plugin.so.

```

```
