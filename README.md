# TZeroHeroGun

## Version 1

Gun app version one. Should communicate via bluetooth with "lead" device
to settle on `t0` (gun). Once it's defined it calculates `t -1` (set) and
`t -2` (on your marks). 

###
`ssh -o PreferredAuthentications=password  cuau@192.168.1.93`


Service UUID `fb7f3ba1-93ab-4eed-9e5f-6197aead8e07`
Characteristic UUID `e4eaaaf2-347d-4f5e-b3f3-8f6e491f3a11`

### To Run 
- `export GOARCH=arm64`
- `go build -o tzeroherogun main.go`
- `sudo setcap 'cap_net_raw,cap_net_admin+eip' ./tzeroherogun`

### Simple benchmarks

#### Standard approach

mean: 25.527203ms
median: 50.691µs
min: 0s
max: 267.75235ms
