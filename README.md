# ShareDocs

Secure document sharing application written in Go, inspired by https://github.com/mfts/papermark.

<img src="gopher-docs.png" width="200" />

## Plan
Foundation & Setup

- [x] Initialize Go module and install dependencies
- [x] Set up PostgreSQL connection with GORM
- [ ] Create database models and migrations
- [ ] Set up Gin/Echo router with basic middleware
- [ ] Implement configuration management (env vars)
- [ ] Basic project structure and error handling

Authentication & File Upload

- [ ] Implement user registration/login with JWT
- [ ] Create auth middleware for protected routes
- [ ] Build file upload endpoint with validation
- [ ] Implement file storage (local or S3/Heroku)
- [ ] Add file type validation (PDF, images,- [ ] presentations)
- [ ] Basic document CRUD operations

Document Versioning & Preview Generation

- [ ] Implement document versioning logic
- [ ] Build preview generation service

Shareable Links & Security

- [ ] Implement shareable link generation
- [ ] Add password protection for links
- [ ] Create link access validation
- [ ] Build public document viewing endpoint
- [ ] Add link expiration handling
- [ ] Security hardening (rate limiting, input validation)


## API Endpoints

__Auth__
```
POST /api/auth/register
POST /api/auth/login
POST /api/auth/refresh
```

__Documents__
```
GET    /api/documents              # List user documents
POST   /api/documents              # Upload new document
GET    /api/documents/:id          # Get document details
PUT    /api/documents/:id          # Update document
DELETE /api/documents/:id          # Delete document
POST   /api/documents/:id/version  # Upload new version
GET    /api/documents/:id/preview  # Get document preview
```

__Shareable Links__
```

```
POST   /api/documents/:id/links    # Create shareable link
GET    /api/documents/:id/links    # List document links
PUT    /api/links/:linkId          # Update link settings
DELETE /api/links/:linkId          # Delete link
GET    /api/shared/:token          # Access shared document (public)
POST   /api/shared/:token/verify   # Verify password for protected link
```