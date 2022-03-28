# rd(register, discover)

## Start

- **start register**
    - `make -C rd/register`
    - `cd rd && make register`
    - `cd rd/register && go run .`
- `make -C rd/discover` **start discover**
    - `make -C rd/discover`
    - `cd rd && make discover`
    - `cd rd/discover && go run .`

## Start with params

```bash
cd register && go run . [-endpoints <http://127.0.0.1:2379>] [-registerName <registerName>] [-serviceKey <server>] \
[-port <port>] [-nodeName <nodeName>] [-internal <10s>]

cd discover && go run . [-endpoints <http://127.0.0.1:2379>] [-discoverName <discoverName>] [-scheme <scheme>] \
[-serviceKey <serviceKey>] [-nodeName <nodeName>] [-internal <10s>]
```
