# ShareDocs

Secure document sharing application written in Go, inspired by https://github.com/mfts/papermark.


## Plan
Foundation & Setup

- [ ] Initialize Go module and install dependencies
- [ ] Set up PostgreSQL connection with GORM
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