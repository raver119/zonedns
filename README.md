# ZonedNS

### What is it?

ZonedNS is a plugin for a [CoreDNS](https://coredns.io) that provides authoritative resolving using MySQL as a backend.
The only real difference from other plugins is specific to A/AAAA mapping: it assumes you have certain "deployment zones", 
with fixed set of A/AAAA records, and domain names can be mapped to these zones ONLY.    

### How to use?

1. Follow CoreDNS build instructions, and add `zonedns:github.com/raver119/zonedns` to `plugin.cfg`
2. Provide MySQL credentials, i.e. via ENV in your K8S deployment description.
   ```bash
   DB_HOST="localhost"
   DB_USERNAME="xxxx"
   DB_PASSWORD="yyyyy"
   DB_NAME="database"
   ```
3. Add the plugin to your Corefile. I.e.:
```
.:1053 {
   cache 3600
   zonedns
}
```

In your management application:

1. Import the API package into your application: `zonedns:github.com/raver119/zonedns-go-api`
2. Use `ZoneStorage` or `ZoneReader` implementation to fetch/add/remove/update records/domains/zones in your app


### Requirements:

- CoreDNS
- MySQL 8.x+
- Golang 1.13+