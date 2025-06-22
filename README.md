# Vega AI Landing Page

Landing page for Vega AI, an AI-powered job search assistant.

## Project Structure

```plaintext
landing-page/
├── web/                 # Frontend
│   ├── assets/
│   │   ├── css/
│   │   │   └── styles.css    # Custom styles and animations
│   │   ├── images/           # Icons and favicon files
│   │   └── js/
│   │       ├── script.js     # Main JavaScript functionality
│   │       └── feedback.js   # Feedback form handling
│   ├── index.html            # Main landing page
│   ├── robots.txt
│   └── site.webmanifest
├── api/                 # Backend
│   ├── function.go           # Cloud Function entry point
│   ├── go.mod               # Go module dependencies
│   ├── go.sum               # Go module checksums
│   └── internal/            # Internal packages
│       ├── application.go   # HTTP handler and routing
│       ├── constants.go     # Action constants
│       └── actions/
│           ├── feedback.go  # Feedback handling logic
│           └── feedback_test.go
└── README.md           # This file
```

## Tech Stack

**Frontend:**

- HTML5 + Tailwind CSS
- Vanilla JavaScript
- Mobile-first responsive design
- Glass morphism effects and animations

**Backend:**

- Google Cloud Functions (2nd gen) with Go
- Google Sheets API for data storage
