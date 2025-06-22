# Vega AI Landing Page

Landing page for Vega AI, an AI-powered job search assistant.

## Project Structure

```plaintext
landing-page/
├── web/                 # Frontend (landing page)
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
├── api/                 # Backend (feedback collection)
│   └── main.go              # Cloud Function for feedback
└── README.md           # This file
```

## Tech Stack

**Frontend:**

- HTMX + Tailwind CSS
- Vanilla JavaScript
- Mobile-first responsive design
- Glass morphism effects and animations

**Backend:**

- Google Cloud Functions (2nd gen) with Go
- Google Sheets API for data storage
