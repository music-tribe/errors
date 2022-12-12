# errors
Standardized error reporting for music tribe golang projects

## Installation
To use this package in your go program, open your terminal and run the command... 
```
go get github.com/music-tribe/errors
```

## In use
To init a new storage error...
```golang
import (
  "github.com/music-tribe/errors"
  "github.com/music-tribe/uuid"

  "some/local/path/database"
)

func (svc *service)someMethod(id) error {
  if err := svc.db.Get(id); err != nil {
    if err == database.NotFoundError {
      return errors.NewCloudError(404, "add your own error message here")
    }
    return errors.NewCloudError(500, err.Error())
  }
}
```

