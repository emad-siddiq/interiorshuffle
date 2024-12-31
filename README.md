# Sample Containerized Application

A property search backend and default Next.js frontend with Redis/Postgresql support containerized with Docker. 

## Run

Make sure you have Docker running and run
```
chmod +x launch.sh
./launch.sh
```
This should launch the backend on `localhost:8080` and the frontend on `localhost:3000`

## Test

To run backend tests, cd into backend/test and run `go test -v`
