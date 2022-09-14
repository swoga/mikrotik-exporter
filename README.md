# mikrotik-exporter
`mikrotik-exporter` is a Prometheus exporter written in Go with the goal to export all possible metrics from MikroTik devices.  
It is not predetermined which metrics are collected, you can create your own modules.  
Some modules are shipped with the program, see [here](/dist/modules).  

## Probing
Targets can be probed by requesting:
<pre>http://localhost:9436/probe?<b>target=xxx</b></pre>
The modules defined at the target configuration can be overwritten via the query string:
<pre>http://localhost:9436/probe?target=xxx&<b>modules=interface,health</b></pre>
For troubleshooting there are also two log levels available:
<pre>http://localhost:9436/probe?target=xxx&<b>debug=1</b></pre>
<pre>http://localhost:9436/probe?target=xxx&<b>trace=1</b></pre>

## Command line flags
<pre>
--config.file=config.yml
--debug
--trace
</pre>

## Docker image
The docker image is available on Docker Hub, Quay.io and GitHub.

<pre>
docker pull swoga/mikrotik-exporter
docker pull quay.io/swoga/mikrotik-exporter
docker pull ghcr.io/swoga/mikrotik-exporter
</pre>

You just need to map your config file into the container at `/etc/mikrotik-exporter/config.yml`
<pre>
docker run -v config.yml:/etc/mikrotik-exporter/config.yml swoga/mikrotik-exporter
</pre>

## Configuration

`mikrotik-exporter` can reload its configuration files at runtime via `SIGHUP` or by sending a request to `/-/reload`.  

### main config
```yaml
[ listen: <string> | default = :9436 ]
[ metrics_path: <string> | default = /metrics ]
[ probe_path: <string> | default = /probe ]
[ reload_path: <string> | default = /-/reload ]

[ namespace: <string> | default = mikrotik ]
[ username: <string> ]
[ password: <string> ]

config_files:
  [ - <string> ... | default = ./conf.d/* ]

[ connection_cleanup_interval: <int> | default = 60 ]
[ connection_use_timeout: <int> | default = 300 ]

targets:
  [ - <target> ... ]
modules:
  [ - <module> ... ]
module_extensions:
  [ - <module_extension> ... ]
```

### conf.d
```yaml
targets:
  [ - <target> ... ]
modules:
  [ - <module> ... ]
module_extensions:
  [ - <module_extension> ... ]
```

### `<target>`
```yaml
name: <string>
address: <string>
[ username: <string> | default = main.username ]
[ password: <string> | default = main.password ]
[ timeout: <int> | default = 10 ]
[ queue: <int> | default = 1000 ]
variables:
  [ <string>: <string> ]
modules:
  [ - <string> ... ]
```

### `<module>`
```yaml
name: <string>
commands:
  [ - <command> ... ]
```

### `<template_string>`
fields of this type support value substitution with values of parent variables  
syntax: `{name_of_variable}`

### `<command>`
```yaml
command: <template_string>
[ timeout: <int> | default = 10 ]
[ prefix: <string> ]

metrics:
  [ - <metric> ... ]
labels:
  [ - <label/variable> ... ]
variables:
  [ - <label/variable> ... ]

sub_commands:
  [ - <command> ... ]
```

### `<param>` base for `<metric>` and `<label/variable>`
```yaml
# either param_name or value must be set
[ param_name: <string> ]
# static value for this param
[ value: <template_string> ]
# value used if not found in API response
[ default: <template_string> ]
# only relevant for param_type = datetime
[ datetime_type: tonow / fromnow / timestamp | default = fromnow ]
# only relevant for param_type = bool
[ negate: <bool> ]

# remapping is stopped after the first match in remap_values or remap_values_re
# remapping to null, stops further processing of this parameter
remap_values:
  [ <string>: <string> / null ]
remap_values_re:
  [ <regex>: <string> / null ]
```

### `<metric>`
```yaml
# derives from param
<param>
[ param_type: int / bool / timespan / datetime | default = int]

# either metric_name or param_name must be set
[ metric_name: <string> | default = param_name ]
metric_type: counter / gauge
[ help: <string> ]

labels:
  [ - <label/variable> ]
```

### `<label/variable>`
```yaml
# derives from param
<param>
[ param_type: string / int / bool / timespan / datetime | default = int]

# either label_name or param_name must be set
[ label_name: <string> | default = param_name ]
```

### `<module_extension>`
module extensions are matched by `name`
```yaml
name: <string>
commands: 
  [ - <command_extension> ... ]
```

### `<command_extension>`
command extensions are matched by `command`
```yaml
command: <string>

metrics:
  [ - <metric_extension> ... ]
labels:
  [ - <label/variable_extension> ... ]
variables:
  [ - <label/variable_extension> ... ]

sub_commands:
  [ - <command_extension> ... ]
```

### `<metric_extension>`
metric extensions are matched by `metric_name`
```yaml
# derives from metric
<metric>

extension_action: add / overwrite / remove
```

### `<label/variable_extension>`
label/variable extensions are matched by `label_name`
```yaml
# derives from label/variable
<label/variable>

extension_action: add / overwrite / remove
```
