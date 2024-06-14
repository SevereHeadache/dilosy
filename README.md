# DiLoSy
**Di**stributed **Lo**gging **Sy**stem

## Install
1. Create _config.yaml_ (see _config.yaml.example_)
2. Add private keys to _./keys_
3. `go build -C ./src -o main`
4. `./src/main`

## Docker
1. Create _.env_ (see _.env.example_):
    - Set **PORT** - port to the container
2. Create _config.yaml_ (see _config.yaml.example_)
    - Set **port: 80**
3. Add private keys to _./keys_
4. `docker compose up -d`
