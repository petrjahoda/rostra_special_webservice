# Rostra Special WebService


## Installation
* use docker image from https://cloud.docker.com/r/petrjahoda/rostra_special_webservice
* use linux, mac or windows version and make it run like a service

## Description
Go webservice that enables operators to start and end their work


## Todo 
- [ ] save ok, nok and cycle to zapsi
- [ ] write back to syteline
    - [ ] at start
    - [ ] at transfer (OK one record, every NOK one record)
    - [ ] at end (OK one record, every NOK one record)
- [ ] add conditions
    - [ ] check pair_part
    - [ ] check only_amount_transfer
    - [ ] check open_orders versus requested_orders
- [ ] show running orders in table

www.zapsi.eu © 2020
