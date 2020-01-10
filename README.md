# powder-monkey
managing dynomite instances

Powder monkey is a skilled technician who engages in blasting work and handles explosives carefully.
So does `powder-monkey` deal with [dynomite](https://github.com/Netflix/dynomite/) to managing the
* state of dynomite
* bootstrap the dynomite backend for warmup
* backup/restore dynomite backends

## Limitations
The is currently no support for other dynomite backends then redis.

## Build
```
git clone https://github.com/sapcc/powder-monkey.git
cd powder-monkey
make
```

### Update dependencies
```
make vendor
```