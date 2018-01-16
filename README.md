# darkstar

Command and control program in golang

- Rotates tls keys per run.
- TCP hole punching
- Arbitrary code execution (WIP)

# Requirements

A sqlite3 database for the server mode
e.g. `sqlite3 /usr/local/share/darkstar.db`
```
sqlite> CREATE TABLE CAPTURES(uuid varchar(45), ip varchar(45), system varchar(45));
sqlite> .exit
```


## Usage
`darkstar -mode=server -operation=scavange -serverdbpath=/usr/local/share/darkstar.db`
`darkstar -mode=client -serverhostaddress=0.0.0.0`



| Client        | Direction     | Server        |
| ------------- | ------------- | ------------- |
| Message       | ->            |               |
|               | <-            | Operation mode|
| Work          |               |               |
| Message       | ->            |               |
