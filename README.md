# mikrotik-exporter
`mikrotik-exporter` is a Prometheus exporter written in Go with the goal to export all possible metrics from MikroTik devices.  
It is not predetermined which metrics are collected, you can create your own modules.  
Some modules are shipped with the program, see [here](/blob/master/dist/modules).  

## Probing
Targets can be probed by requesting:
<pre>`http://localhost:9436/probe?<b>target=xxx</b>`</pre>
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