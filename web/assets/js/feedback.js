// Minimal Feedback Form JavaScript (HTMX handles most functionality)
document.addEventListener('DOMContentLoaded', function() {
  // Initialize slider value display
  const slider = document.getElementById('setup-difficulty');
  if (slider) {
    // Set initial value
    const initialPercent = ((slider.value - slider.min) / (slider.max - slider.min)) * 100;
    slider.style.setProperty('--value', initialPercent + '%');
  }

  // HTMX handles form submission - no JavaScript override needed
  // Form data will be automatically serialized and sent to the hx-post endpoint
});

// Copy to clipboard functionality is already in main script