# Renderman (Render Manager)

Dynamic server-side rendering using headless Chrome for modern javascript single page application (SPA)

## Running Renderman Go Version

- Install dart `brew install go`
- Install dependencies `go mod download`
- Running :
  - Running chrome headless `go run ./remote`
  - Running renderman `go run ./renderman/main.go`

## Running Renderman Dart Version

- Install dart `brew tap dart-lang/dart && brew install dart`
- Install dependencies `dart pub get`
- Running :
  - Running chrome headless `go run ./remote`
  - Running renderman `dart run ./renderman_dart/main.dart`
