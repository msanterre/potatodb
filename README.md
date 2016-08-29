## Potato DB

Code from the [Let's make a database series](http://alpacalunchbox.com/)

The code at this point is at: `1-basic`


## Running

```
# While in directory

go get
go run main.rb
```

## Usage

#### GET
```
curl http://localhost:5050/db/potato
```

#### SET
```
curl -X POST http://localhost:5050/db/potato -d "This is data"
```

#### DEL
```
curl -X DELETE http://localhost:5050/db/potato
```
