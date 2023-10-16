# SSO (Single sign-on)

Is for authentication and authorizations of the users.

# Specifications
- User
  - Id
  - AvatarId
  - FirstName
  - LastName
  - Email
  - Password
  - TOTPSecret
  - EmailConfirmed
  - CreatedAt
  - UpdatedAt
  - IsActive
- Role
  - Id
  - Name
  - CreatedAt
  - UpdatedAt
  - IsActive
- Permissions
  - Id
  - RoleId
  - Name
  - CreatedAt
  - UpdatedAt
  - IsActive
- SignOns
  - Id
  - UserId
  - Timestamp
  - IPAddress
  - IPTrace
  - CreatedAt
  - UpdatedAt
  - IsActive