Func main(arg) calls handlers(arg) from the handlers package
|
|-- If arg == "test"
|      |
|      V
|   initHandlers reads .test.env and sets dbAccessor accordingly
|      |
|      V
|   StartServerTest() reads .test.env and calls DBAccess() from the DB package
|                |
|                V
|              Migrate() from the DB package with the correct access
|
|-- If arg != "test"
|      |
|      V
|   initHandlers reads .env and sets dbAccessor accordingly
|      |
|      V
|   StartServer() reads .env and calls DBAccess() from the DB package
|             |
|             V
|          Migrate() from the DB package with the correct access
|
|    handlers.JwtSecret(arg) reads .env and sets jwtSecret (key to read tokens)
|
|---|
    |
    V
MainRouter from the router package initiates a router calling routers v1, v2, and v3 also from the router package
    |
    V
HttpServer added and starts listening on the correct port according to .env
    |
    V
Any server call hits MainRouter, which directs to v1, v2, or v3
    |
    V
Handlers manage the functions called by HTTP requests correctly using the utils, res packages, and functions like RespondWithError and RespondWithJSON from the handlers package
