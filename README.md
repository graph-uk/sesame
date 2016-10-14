# Sesame

Windows app, that opens port using win firewall

You can:
  - Open ports for your IP using CURL request
  - Do the same, using HTML form.



## How to open ports with sesame
By curl request:
```sh
$ curl -H "Content-type: application/json" --data "{\"Token\":\"<secretToken>\"}" -X POST <serverIP>:<port>
```
By html form:
  - Open http://<serverIP>:<port>/letmein<URLPostfix>
  - Submit token.


## Installation
```sh
$ go get github.com/graph-uk/sesame
$ go install
```
Run `sesame.exe` as admin in the directory you want to store the config. Default config will be created.

If you run it without admin creds - app will cannot edit firewall rules, and show message about.

## Configure
`config.json` values:

| Name   | Description  |
|----------|:-------------:|
| MasterPort |  Port for sesame commands | 
| SlavePorts |    Ports, that will be opened by request to MasterPort   |
| DurationMinutes | Duration, while slave ports are open |
| Token |  Secret sequence to prevent unauthorized commands | 
| ARR_IP_HeaderField |    Name of request's field contains client IP (for reverse-proxy)   |
| URLPostfix | Part form's URL. Form is disabled with empty parameter. |
Edit config as you need, and restart sesame.
Example:
```json
{
  "MasterPort":3139,
  "SlavePorts":"21,22,31-443",
  "DurationMinutes":240,
  "Token":"bxc9v87b6s9df87vy09a8sd7csado978yOKJHCKSUDT",
  "ARR_IP_HeaderField":"X-Forwarded-For",
  "URLPostfix":"87vy09a8sd7csadoSDcxvlkzjs93"
}
```

## HTTPS
Also you able setup IIS reverse proxy, and call sesame through https. Then edit ARR_IP_HeaderField in config, if needed. 

## Autostart
Easily way to start sesame automatically - [nssm].

[nssm]: <https://nssm.cc/download>