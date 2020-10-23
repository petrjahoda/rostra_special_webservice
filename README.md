[![developed_using](https://img.shields.io/badge/developed%20using-Jetbrains%20Goland-lightgrey)](https://www.jetbrains.com/go/)
<br/>
![GitHub](https://img.shields.io/github/license/petrjahoda/rostra_special_webservice)
[![GitHub last commit](https://img.shields.io/github/last-commit/petrjahoda/rostra_special_webservice)](https://github.com/petrjahoda/rostra_special_webservice/commits/master)
[![GitHub issues](https://img.shields.io/github/issues/petrjahoda/rostra_special_webservice)](https://github.com/petrjahoda/rostra_special_webservice/issues)
<br/>
![GitHub language count](https://img.shields.io/github/languages/count/petrjahoda/rostra_special_webservice)
![GitHub top language](https://img.shields.io/github/languages/top/petrjahoda/rostra_special_webservice)
![GitHub repo size](https://img.shields.io/github/repo-size/petrjahoda/rostra_special_webservice)
<br/>
[![Docker Pulls](https://img.shields.io/docker/pulls/petrjahoda/rostra_special_webservice)](https://hub.docker.com/r/petrjahoda/rostra_special_webservice)
[![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/petrjahoda/rostra_special_webservice?sort=date)](https://hub.docker.com/r/petrjahoda/rostra_special_webservice/tags)
<br/>
[![developed_using](https://img.shields.io/badge/database-MySQL-red)](https://www.mysql.com) [![developed_using](https://img.shields.io/badge/database-SQL_Server-red)](https://www.microsoft.com/en-us/sql-server) [![developed_using](https://img.shields.io/badge/runtime-Docker-red)](https://www.docker.com)


# Rostra Special WebService

![Example](/gif/example.gif)

## Installation
* use docker image from https://cloud.docker.com/r/petrjahoda/rostra_special_webservice

## Description
Go webservice with js frontend that enables operators to start and end their work described [here](logika.pdf)
Service communicates with Zapsi MySQL database and Syteline SQL Server database.
Numeral checks are made in the background before an action (database changes) is made.

## Information
### User check button action
- user.js
- table.js (for displaying data for actual user)
- user.go
### Order check button action
- order.js
- order.go
### Operation check button action
- operation.js
- operation.go
### Workplace check button action
- workplace.js
- workplace.go
### Count check button action
- count.js
- count.go
### Start order button action
- start.js
- start.go
### End order button action
- end.js
- end.go

## Changelog
Updated [here](CHANGELOG.md)

©2020 Petr Jahoda
