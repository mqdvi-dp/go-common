# go-common

Common modules for PVG microservices

## 1. Specific Requirement

- Get token access for gitlab *Setting  -> Access Token*
- Setting git config --global like this [tutorial](https://medium.com/cloud-native-the-gathering/go-modules-with-private-git-repositories-dfe795068db4)
-

  ```bash
  git config --global url."https://${user}:${token}@git.infra.pvg.im".insteadOf "https://git.infra.pvg.im"
  ```

## 2. Installing libs

```bash
export GOSUMDB=off

go get git.infra.pvg.im/go-common

```

## 3. Using this lib in go

```go
package main

import (
  "git.infra.pvg.im/go-common/config"
  "git.infra.pvg.im/go-common/env"
)

func main() {
 config.New(serviceName)
 env.GetString("DSN_MASTER")
}

```
