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
Authentication with the server is handled through Json Web Token (JWT). The token lifespan duration is hard coded to 1 hour.

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
Tenants are globally unique, in the context of this repository, 
you can think of tenant as 

Tenant-related APIs can only be authorized with superadmin's permissions.


### Account
An account represent an entity (user, admin) with permissions to communicate with the licensing server
and perform authorized actions such as creating a new product, checking out a license, or validate license.
At this time, an account can be one of the three roles: ```superadmin```, ```admin```, or ```user```.


### Product
A product is essentially any software or application that you want to license.
Each tenant can have multiple products where each of them can have multiple attributes 
such as supported platforms, product code, etc. Any policy, license created must 
be associated with a product. A single product can have multiple associated policies and licenses.

### Policy
A policy is a set of rules that specify how a license should behave for a product. 
It controls the scopes and limits of licenses issued.

### License
License represents the rights to use the defined product. A license must be associated with a policy.
It allows you to enforce your licensing models. Online and offline verifications are done through license.

### Machine
Machine represents a server or computer on which the license is activated. 
They are used to track license activations and enforce licensing rules.

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
### Use Cases
TODO

---
### API Request
Any request made to the licensing server is assigned a request id that can be used to track the execution
and tracing if needed.

---
### API Response

ALl returned response follow a common format, as shown below:
Example Response
```json
{
    "request_id": "0cc7505f-b62f-439e-add8-4b7ce7b1cb90",
    "code": "00000",
    "message": "OK",
    "server_time": 1736306590,
    "data": {
        "access": "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzdXBlcmFkbWluIiwiZXhwIjoxNzM2MzEwMTkwLCJpYXQiOjE3MzYzMDY1OTAsImlzcyI6ImdvLWxpY2Vuc2UtbWFuYWdlbWVudCIsIm5iZiI6MTczNjMwNjU5MCwicGVybWlzc2lvbnMiOlsibGljZW5zZS5kZWxldGUiLCJsaWNlbnNlLnJldm9rZSIsImxpY2Vuc2UtdXNhZ2UuZGVjcmVtZW50IiwiYWRtaW4uY3JlYXRlIiwiYWRtaW4udXBkYXRlIiwidXNlcl9wYXNzd29yZC51cGRhdGUiLCJlbnRpdGxlbWVudC51cGRhdGUiLCJwb2xpY3kuY3JlYXRlIiwibGljZW5zZS11c2FnZS5yZXNldCIsIm1hY2hpbmUuZGVsZXRlIiwidGVuYW50LmNyZWF0ZSIsImFkbWluLmRlbGV0ZSIsInVzZXIuYmFuIiwibGljZW5zZS10b2tlbnMuZ2VuZXJhdGUiLCJtYWNoaW5lLmNyZWF0ZSIsImVudGl0bGVtZW50LmRlbGV0ZSIsImVudGl0bGVtZW50LnJlYWQiLCJwb2xpY3kucmVhZCIsInBvbGljeV9lbnRpdGxlbWVudHMuZGV0YWNoIiwibGljZW5zZS51cGRhdGUiLCJ0ZW5hbnQuZGVsZXRlIiwidXNlci5kZWxldGUiLCJsaWNlbnNlLmNoZWNrLWluIiwibGljZW5zZS5jaGVjay1vdXQiLCJtYWNoaW5lLWhlYXJ0YmVhdC5yZXNldCIsInBvbGljeS51cGRhdGUiLCJsaWNlbnNlLnJlbmV3IiwibWFjaGluZS51cGRhdGUiLCJwcm9kdWN0LmRlbGV0ZSIsImxpY2Vuc2UtdXNlcnMuYXR0YWNoIiwidGVuYW50LnJlYWQiLCJ1c2VyLnVuYmFuIiwicHJvZHVjdC5yZWFkIiwicG9saWN5X2VudGl0bGVtZW50cy5hdHRhY2giLCJsaWNlbnNlLnJlaW5zdGF0ZSIsInByb2R1Y3QudXBkYXRlIiwibGljZW5zZS5yZWFkIiwiYWRtaW4ucmVhZCIsInVzZXJfcGFzc3dvcmQucmVzZXQiLCJ0ZW5hbnQudXBkYXRlIiwidXNlci5yZWFkIiwicG9saWN5LmRlbGV0ZSIsImxpY2Vuc2UudmFsaWRhdGUiLCJsaWNlbnNlLWVudGl0bGVtZW50cy5hdHRhY2giLCJlbnRpdGxlbWVudC5jcmVhdGUiLCJsaWNlbnNlLnN1c3BlbmQiLCJsaWNlbnNlLWVudGl0bGVtZW50cy5kZXRhY2giLCJtYWNoaW5lLmNoZWNrLW91dCIsInVzZXIuY3JlYXRlIiwidXNlci51cGRhdGUiLCJsaWNlbnNlLmNyZWF0ZSIsIm1hY2hpbmUtaGVhcnRiZWF0LnBpbmciLCJwcm9kdWN0LmNyZWF0ZSIsInByb2R1Y3RfdG9rZW5zLmdlbmVyYXRlIiwibWFjaGluZS5yZWFkIiwibGljZW5zZS11c2FnZS5pbmNyZW1lbnQiLCJsaWNlbnNlLXBvbGljeS51cGRhdGUiLCJsaWNlbnNlLXVzZXJzLmRldGFjaCJdLCJzdGF0dXMiOiJhY3RpdmUiLCJzdWIiOiJzdXBlcmFkbWluIiwidGVuYW50IjoiKiJ9.PCSut7oHJAEgjFWCCKG8FJNptJEYFXOgC7AVzmWTBxx6MdseHMkdbMnVuXgDZhE96-TABDECvlOjfQd4KUUxAg",
        "expire_at": 1736310190
    }
}
```
Every response includes a ```request_id``` in uuid v4 format, a ```server_time``` in unix epoch int64 format, 
error ```code``` and ```message```, if any. For every successful request, a ```Content-Digest``` header is included. Client applications should verify this 
hash using sha256 algorithm.
```text
sha256=bc08aa0cbd668c66d1a40e447a64cf887824670c7d098f75fcd3d8e0280b158f
```








