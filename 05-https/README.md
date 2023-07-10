# Use `mkcert` to create certificate for localhost
https://github.com/FiloSottile/mkcert

1. Build `mkcert`
```
git clone https://github.com/FiloSottile/mkcert && cd mkcert
go build -ldflags "-X main.Version=$(git describe --tags)"
```

2. Create a local CA
```
mkcert -install
```

3. Create certificate for localhost
```
mkcert localhost
```
- Certificate: `localhost.pem`
- Private key: `localhost-key.pem`
