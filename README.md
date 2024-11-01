# go-license-management
An open-source license management written in Go 

### Authentication and Authorization


### Basic terminology

| Term        | Description                                                                                                     |
|-------------|-----------------------------------------------------------------------------------------------------------------|
| **License** | Licenses represent an entitlement, i.e. you grant a licensee permission to use something.                       |
| **Machine** | Machines represent a device or node that a license is allowed to be used with.                                  |
| **Policy**  | Policies define behavior for different license types, e.g., Timed Trial, Basic, and Pro types.                  |
| **User**    | Users represent an identity for an end-user, or licensee, of your software.                                     |



### Supported License Types

| License Type            | Expiration Date | Activation Limits | Feature Limits | Offline Support |
|-------------------------|-----------------|-------------------|----------------|-----------------|
| **Perpetual License**   | No              | Optional          | Optional       | Yes             |
| **Timed License**       | Yes             | Optional          | Optional       | Yes             |
| **Floating License**    | Optional        | > 0               | Optional       | Yes             |
| **Node-locked License** | Optional        | 1                 | Optional       | Yes             |
| **Feature License**     | Optional        | Optional          | Yes            | Yes             |
