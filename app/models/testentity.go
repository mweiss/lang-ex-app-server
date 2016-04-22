// models/testentity.go
package models

import  "time" // if you need/want

type TestEntity struct {          // example user fields
    Id                    int64
    Name                  string
    EncryptedPassword    []byte
    Password              string      `sql:"-"`
    CreatedAt             time.Time
    UpdatedAt             time.Time
    DeletedAt             time.Time     // for soft delete
}