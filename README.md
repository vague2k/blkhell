# Blkhell.
<img width="800" alt="Screenshot 2026-02-19 at 9 10 54 AM" src="https://github.com/user-attachments/assets/76873ebf-1f6b-413e-b2a4-de54bdf466c0" />


When working on blackheaven projects we're always having to message each other about which logo to use, or dig endless through the email to find the asset we're looking for. 

I took initiative to create a tool that lets us house all label assets in a single place, so we always know where to look.

Right now the scope of the project only handles images (label logos and other such assets.) but future possibilities may include specific tooling for keeping track of label projects.

## Main project structure
```
.
├── cmd/                # Application entrypoints
│   ├── cli/            # CLI binary
│   └── server/         # HTTP server binary
│
├── internal/           # Private packages
│   └── blkhell/        # blkhell CLI implementation details
│
├── server/             # Server application code
│   ├── auth/           # Authentication, sessions, user identity
│   ├── database/       # Database layer and persistence logic
│   ├── handlers/       # HTTP handlers grouped by route/action
│   ├── router.go       # Route definitions and middleware wiring
│   └── server.go       # Server bootstrap and configuration
│
└── views/              # Frontend (templ + static assets)
    ├── assets/         # CSS, JS, fonts, images
    ├── components/     # Reusable templ components
    ├── layouts/        # Page layouts / shells
    ├── pages/          # User-facing pages
    └── templui/        # TemplUI component library
```
