# API

## Incidents

Turbulence API server can perform several types of tasks against a subset of instances in several deployments. Types of tasks and selection of instances is represented by an incident.

See 'Incident Types' section below for details on which task types are available.

Instance selection is based on `Deployments` (array of hashes; required) configuration. An incident may affect one or deployments from which one or more jobs can be selected. For each job selected, its instances are filtered via `Indices` (array of ints) or `Limit` (string) keys. Limit value may be one of the following:

- `10`: kill 10 *random* instances
- `10%`: kill random 10% of all instances
- `5-10%`: kill random 5% to 10% of all instances

Endpoints:

- `POST /api/v1/incidents`
- `GET /api/v1/incidents`
- `GET /api/v1/incidents/:id`

```
curl -vvv -k -X POST https://user:pass@10.244.8.2:8080/api/v1/incidents -H 'Accept: application/json' -d@example.json
```

Create request:

```json
{
	"Tasks": [{
		"Type": "stress",
		"Timeout": "10m",
		
		"NumCPUWorkers": 1
	}],

	"Deployments": [{
		"Name": "cf",
		"Jobs": [{
			"Name": "*_z1",
			"Limit": "10-50%"
		},{
			"Name": "api_z2",
			"Indices": [0,1,2,3]
		}]
	},{
		"Name": "etcd",
		"Jobs": [{
			"Name": "node_z1",
			"Limit": "10-50%"
		}]
	}]
}
```

Response:

```json
{
  "ID": "d77adc3b-1de4-4e12-4bee-b325adfbecbd",

  "Tasks": [ ... ],
  "Deployments": [ ... ],

  "ExecutionStartedAt": "0001-01-01T00:00:00Z",
  "ExecutionCompletedAt": "",

  "Events": null
}
```

---
## Scheduled Incidents

Turbulence API server can create new incidents based on a schedule. `Schedule` (string; required) can be specified in cron format or with one of the shorthands:

- `@yearly`
- `@monthly`
- `@weekly`
- `@daily`
- `@hourly`
- `@every X` where X is value accepted by the [golang's time.ParseDuration](http://golang.org/pkg/time/#ParseDuration)

`Incident` (hash; required) is specified in the exactly the same way as when creating a single incident.

Endpoints:

- `POST /api/v1/scheduled_incidents`
- `GET /api/v1/scheduled_incidents`
- `GET /api/v1/scheduled_incidents/:id`
- `DELETE /api/v1/scheduled_incidents/:id`

```
curl -vvv -k -X POST https://user:pass@10.244.8.2:8080/api/v1/scheduled_incidents -H 'Accept: application/json' -d@example.json
```

Create request:

```json
{
	"Schedule": "@every 1m",
	"Incident": { ... }
}
```

Response:

```json
{
  "ID": "bf43eed7-91c7-4983-5895-44b9a18a5461",

  "Schedule": "@every 1m",
  "Incident": { ... }
}
```

---
## Incident Tasks

Currently there are four support task types that can be included in an incident. All but the `kill` task requires for a `Timeout` key to be set so that the task can complete.

### Kill

Deletes the VM associated with an instance. Turbulence API server uses collocated CPI job directly. It does not use the Director to the delete the VM.

Example:

```json
{
	"Type": "kill"
}
```

### Stress

Stresses different subsystems on the VM associated with an instance.

One or more of the following configurations must be selected:

- CPU
  - set `NumCPUWorkers` (int; required)

- IO
  - set `NumIOWorkers` (int; required)

- RAM
  - set `NumMemoryWorkers` (int; required)
  - set `MemoryWorkerBytes` (string; required). Must be suffixed with B,K,M,G.

- HDD
  - set `NumHDDWorkers` (int; required)
  - set `HDDWorkerBytes` (string; required). Must be suffixed with B,K,M,G.

Example:

```json
{
	"Type": "stress",
	"Timeout": "10m", // Times may be suffixed with s,m,h,d,y

	"NumCPUWorkers": 1,

	"NumIOWorkers": 1,

	"NumMemoryWorkers": 1,
	"MemoryWorkerBytes": "10K"
}
```

### Firewall

Blocks incoming and outgoing traffic from the VM associated with an instance. Useful for simulating network partitions. By default BOSH Agent and SSH on the VM will continue to operate.

Currently iptables is used for dropping packets from INPUT and OUTPUT chains.

Optionally specify:

- set `BlockBOSHAgent` (bool) to false to block access to the BOSH Agent

Example:

```json
{
	"Type": "firewall",
	"Timeout": "10m" // Times may be suffixed with ms,s,m,h
}
```

### Control Network

Controls network quality on the VM associated with an instance. Does not affect `lo0`.

Currently [tc](http://www.lartc.org/manpages/tc.txt) is used to control package delay and loss.

One or both of the following configurations must be selected:

- packet delay
  - set `Delay` (string; required). Must be suffixed with `ms`.
  - set `DelayVariation` (string; optional). Must be suffixed with `ms`. Default is `10ms`.

- packet loss
  - set `Loss` (string; required). Must be suffixed with `%`.
  - set `LossCorrelation` (string; optional). Must be suffixed with `%`. Default is `75%`.

Example:

```json
{
	"Type": "control-net",
	"Timeout": "10m", // Times may be suffixed with ms,s,m,h

	"Delay": "50ms"
}
```
