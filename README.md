# Vega Landing Page

This is the product landing page for Vega - an AI-powered job search platform.

## Features

- Modern, responsive design using Tailwind CSS
- Star navigation theme matching the main application
- Animated star field background
- Smooth scrolling navigation
- Browser extension installation instructions
- Docker setup guide
- Mobile-friendly layout

## Usage

### Option 1: Direct Opening

Simply open `index.html` in your web browser.

### Option 2: Local Server

For better performance and to avoid CORS issues:

```bash
# Using Python 3
python -m http.server 8080

# Using Node.js
npx http-server -p 8080

# Using PHP
php -S localhost:8080
```

Then visit `http://localhost:8080`

### Option 3: Deploy to Static Hosting

The landing page can be deployed to any static hosting service:

- GitHub Pages
- Netlify
- Vercel
- AWS S3
- CloudFlare Pages

## Customization

### Colors

The color scheme matches the main Vega application:

- Primary: `#0D9488` (Teal)
- Secondary: `#F59E0B` (Amber)
- Dark backgrounds: Slate color palette

### Content

Edit the following sections in `index.html`:

- Hero section text and statistics
- Features descriptions
- Browser extension instructions
- Installation steps
- Footer links

### Images/Videos

To add images or videos:

1. Create an `assets` folder
2. Add your media files
3. Update the src paths in the HTML

## Structure

```
landing-page/
├── index.html      # Main landing page
├── README.md       # This file
└── assets/         # (Create this for images/videos)
```

## SEO Optimization

The landing page includes:

- Semantic HTML structure
- Meta description
- Proper heading hierarchy
- Alt text placeholders for images
- Schema.org structured data (can be added)

## Performance

- Minimal JavaScript for animations
- Tailwind CSS via CDN (can be localized)
- No heavy frameworks
- Optimized for fast loading

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)
- Mobile browsers

## License

Same as the main Vega project - MIT License
