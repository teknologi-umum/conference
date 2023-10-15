# Backend

## Developing

To develop the backend, just create integration test and make sure it's running on CI.

But, if you want to run it locally, copy the `configuration.example.yml` into `configuration.yml`,
then start the Docker Compose file on the root directory using `docker compose up -d postgres mailcrab`.

Your `configuration.yml` file should be similar to:

```yaml
feature_flags:
  registration_closed: false

environment: local

database:
  host: localhost
  port: 5432
  user: conference
  password: VeryStrongPassword
  database: conference

port: 8080

mailer:
  hostname: localhost
  port: 1025
  from: administrator@localhost
  password:

blob_url: file:///tmp/teknologi-umum-conference

signature:
  public_key: 2bb6b9b9e1d9e12bfdd4196bfba6a081156ac...
  private_key: 48d0ca64011fec1cb23b21820e9f7e880843e71f236b7f8decfe3568f...

validate_payment_key: 24326124313024514d56324d782f4a7342446f36363653784b324175657341...
```

Generate the `signature.public_key` and `signature.private_key` using this simple Go script:

```go
package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
)

func main() {
	pub, priv, _ := ed25519.GenerateKey(nil)
	fmt.Println(hex.EncodeToString(pub))
	fmt.Println(hex.EncodeToString(priv))
}
```

Generate the `validate_payment_key` using this simple Go script:

```go
package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	passphrase := "test"
	value, _ := bcrypt.GenerateFromPassword([]byte(passphrase), 10)
	fmt.Printf("%x", value)

}
```