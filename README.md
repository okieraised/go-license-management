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

To build this repository as a Docker image, clone this repository to your local computer and run:
```shell
make build
```
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






