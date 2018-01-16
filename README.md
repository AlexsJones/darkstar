# darkstar

Command and control program in golang

- Rotates tls keys per run.
- TCP hole punching
- Arbitrary code execution (WIP)

## Usage
`darkstar -mode=server -operation=scavange`
`darkstar -mode=client -serverhostaddress=0.0.0.0`



| Client        | Direction     | Server        |
| ------------- | ------------- | ------------- |
| Message       | ->            |               |
|               | <-            | Operation mode|
| Work          |               |               |
| Message       | ->            |               |
