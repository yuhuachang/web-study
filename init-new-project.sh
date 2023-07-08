#!/usr/bin/bash

PROJ=$1

if [ -z "$PROJ" ]; then
    echo "Usage: $0 <project-name>"
    exit 1
fi

# Create a new project directory
mkdir $PROJ

# Setup a new Go project
cd $PROJ
go mod init factory95.com/$PROJ

# Create main.go
cat <<__EOF__ > main.go
package main

import (
  "fmt"
)

func main() {
  fmt.Println("Hello, World!")
}
__EOF__
