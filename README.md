# go-license-management
An open-source license management written in Go 

### Authentication and Authorization


### Device Registration

### Request


### Response
For every successful request (200 or 201), a Content-Digest header is included. Client application should verify this 
hash using sha256 algorithm.
```text
sha256=bc08aa0cbd668c66d1a40e447a64cf887824670c7d098f75fcd3d8e0280b158f
```




### Basic terminology
1. **Product**: A product represents a software application or a digital product. 
    It's essentially an identifier that groups together licenses, accounts, and policies related to a 
    specific software application. You can think of a product as a single version or line of software. 
    If your company offers multiple apps or product lines, you would create separate products for each.
    Each product has unique settings, including metadata (like version and description) and configurations,
    which can be customized for licensing purposes.

2. **Policy**: A policy defines the rules or restrictions for licensing within a product. It’s like a template or set of rules 
    that controls how the license for that product behaves. Policies are flexible and can be tailored to enforce specific 
    licensing conditions, such as: 
     * License duration (e.g., perpetual, trial, subscription-based)
     * User limits (e.g., single user or multiple users)
     * Feature gating (e.g., enabling specific features for different license tiers)
     * Activation limits (e.g., limiting the number of devices a license can activate)

    You can create multiple policies for the same product to support different license types, allowing you to offer various plans, 
    such as basic, pro, or enterprise versions of your software, with each policy defining unique rules for usage and access.

### Supported License Scheme
* **Ed25519**
* **RSA2048**
* **RSA-JWT-RS256**

### Supported License Types

| License Type            | Expiration Date | Activation Limits | Feature Limits | Offline Support |
|-------------------------|-----------------|-------------------|----------------|-----------------|
| **Perpetual License**   | No              | Optional          | Optional       | Yes             |
| **Timed License**       | Yes             | Optional          | Optional       | Yes             |
| **Floating License**    | Optional        | > 0               | Optional       | Yes             |
| **Node-locked License** | Optional        | 1                 | Optional       | Yes             |
| **Feature License**     | Optional        | Optional          | Yes            | Yes             |

### Roadmap
- [x] Product APIs 