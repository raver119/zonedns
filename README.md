# ZonedNS

### What is it?

ZonedNS is a plugin for a [CoreDNS](https://coredns.io) that provides authoritative resolving using MySQL as a backend.

### How to use?

1. Follow CoreDNS build instructions, and add `zonedns:github.com/raver119/zonedns/plugin` to `plugin.cfg`
2. Provide MySQL credentials, i.e. via ENV in your K8S deployment description.
   ```bash
   DB_HOST="localhost"
   DB_USERNAME="xxxx"
   DB_PASSWORD="yyyyy"
   DB_NAME="database"
   ```
2. Include the API to your application by importing `zonedns:github.com/raver119/zonedns/api`
3. Use `ZoneStorage` or `ZoneReader` implementation to fetch/add/remove/update records/domains/zones in your app

### Requirements:

- CoreDNS
- MySQL 8.x+
- Golang 1.13+