```      
                                             _/      
    _/_/_/  _/  _/_/  _/    _/  _/_/_/    _/_/_/_/   
 _/        _/_/      _/    _/  _/    _/    _/        
_/        _/        _/    _/  _/    _/    _/         
 _/_/_/  _/          _/_/_/  _/_/_/        _/_/      
                        _/  _/                       
                     _/_/  _/                      
```
Simple CLI tool to encrypt / decrypt text.

# Installation

```sh
go get github.com/btb55/crypt
go install
```

# Usage example

## Encrypt:
```sh
crypt -e -k mykey 'secret text'
```
```sh
// output
-------------------------------------------
26vBzKOIOpg06qGluKzJuFmBcX2_Rt7tSMRUF3-LHo0
-------------------------------------------
```

## Decrypt:
```sh
crypt -d -k mykey 26vBzKOIOpg06qGluKzJuFmBcX2_Rt7tSMRUF3-LHo0
```
```sh
// output
-----------
secret text
-----------
```
