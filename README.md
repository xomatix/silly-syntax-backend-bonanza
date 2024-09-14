# silly-syntax-backend-bonanza

- [Instalation and building](#instalation)
  - [frontend](#frontend---admin-dashboard)
  - [backend](#frontend---admin-dashboard)

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
