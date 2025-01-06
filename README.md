# go-license-management
An open-source license management written in Go 

---
### Introduction
This is my personal attempt to develop a License Management Service in pure Golang.
Currently, this repository is a work-in-progress (WIP), the design and implementation are based on what I know and suggestions from the internet,
things may be unconventional and hopefully, will get better as I finish it.

This service will provide APIs for license management and validation. 
However, what to do with this information is entirely up to the client applications.

---
### Build and Installation
This repository requires PostgreSQL as the backend database. 
Support for other databases are not yet implemented at this time.

If you want to build this service as a Docker image and push to the local registry on your computer,
I have included an example docker-compose file. Run the following command to set up a local registry:
```shell
docker compose up -d local-registry
```
You should have a up-and-running local registry on port 5000 in no time.

To build this repository as a Docker image, clone this repository to your local computer and run:
```shell
make build
```

To push this repository to the local registry, run:
```shell
make push-local
```

If you want to build this service as is, run the good old:
```shell
go build .
```

Otherwise, if you just want to test it out, you can also do:
```shell
go run .
```

---
### Environmental Variable and Config

This service support 2 ways to specify environmental variables, using ```conf.toml``` and ```.env```

| Using Env          | Using config.toml  | Default Value  | Description                               |
|--------------------|--------------------|----------------|-------------------------------------------|
| SERVER__MODE       | [server]mode       | debug          | Server mode                               |
| SERVER__HTTP_PORT  | [server]http_port  | 8888           | Port to listen                            |
| POSTGRES__HOST     | [postgres]host     | 127.0.0.1      | IP/Hostname of the postgres db            |
| POSTGRES__PORT     | [postgres]port     | N/A            | Port of the postgres db                   |
| POSTGRES__USERNAME | [postgres]username | N/A            | postgres username to use                  |
| POSTGRES__PASSWORD | [postgres]password | N/A            | postgres port to use                      |
| POSTGRES__DATABASE | [postgres]database | licenses       | database name, must be created beforehand |

---
### Authorization and Permissions
Authentication with the server is handled through Json Web Token (JWT). The default token lifespan duration is 1 hour.

Authorization (Permissions) is handled using Casbin. The model configuration is as follows:
```text
[request_definition]
r = dom, sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = g(r.dom, r.sub, p.sub) && r.obj == p.obj && r.act == p.act || r.sub == "superadmin"
```

for more information about Casbin, refer to [Casbin](https://casbin.org/docs/overview/)

---
### Tenant


### Account


### Product


### Policy


### License


### Machine

---
### Roadmap
- [x] Tenant APIs
- [x] Account APIs
- [x] Product APIs
- [x] Policy APIs
- [x] Entitlement APIs
- [x] License APIs
- [x] Machine APIs
- [x] Authentication APIs
- [x] Token APIs
- [x] Product APIs
- [ ] Other Database Supports

---
### Supported License Scheme
*  [x] **Ed25519**
*  [x] **RSA2048**

---
### Supported License Types
| License Type            | Expiration Date | Activation Limits | Feature Limits | Offline Support |
|-------------------------|-----------------|-------------------|----------------|-----------------|
| **Perpetual License**   | No              | Optional          | Optional       | Yes             |
| **Timed License**       | Yes             | Optional          | Optional       | Yes             |
| **Floating License**    | Optional        | > 0               | Optional       | Yes             |
| **Node-locked License** | Optional        | 1                 | Optional       | Yes             |
| **Feature License**     | Optional        | Optional          | Yes            | Yes             |


---
### Authentication and Authorization
**Todo**

---
### API Response
For every successful request, a Content-Digest header is included. Client application should verify this 
hash using sha256 algorithm.
```text
sha256=bc08aa0cbd668c66d1a40e447a64cf887824670c7d098f75fcd3d8e0280b158f
```






