// Minimal Feedback Form JavaScript (HTMX handles most functionality)
document.addEventListener('DOMContentLoaded', function() {
  // Initialize slider value display
  const slider = document.getElementById('setup-difficulty');
  if (slider) {
    // Set initial value
    const initialPercent = ((slider.value - slider.min) / (slider.max - slider.min)) * 100;
    slider.style.setProperty('--value', initialPercent + '%');
  }

  // Since backend is not implemented, intercept form submission for demo
  const form = document.getElementById('feedback-form');
  if (form) {
    form.addEventListener('submit', function(e) {
      e.preventDefault();
      
      // Collect form data
      const formData = new FormData(form);
      
      // Get all checked setup issues
      const setupIssues = [];
      form.querySelectorAll('input[name="setupIssues"]:checked').forEach(checkbox => {
        setupIssues.push(checkbox.value);
      });
      
      const data = {
        helpfulness: formData.get('helpfulness'),
        setupDifficulty: parseInt(formData.get('setupDifficulty')),
        docsQuality: formData.get('docsQuality'),
        setupIssues: setupIssues.join(', ') || '',
        additionalFeedback: formData.get('additionalFeedback') || '',
        email: formData.get('email') || '',
        source: 'landing-page',
        timestamp: new Date().toISOString()
      };
      
      console.log('Form data to submit:', data);
      
      // Show loading state
      form.classList.add('htmx-request');
      
      // Simulate API call
      setTimeout(() => {
        // Show success message
        const formMessage = document.getElementById('form-message');
        formMessage.innerHTML = `
          <div class="bg-green-500/10 border border-green-500/30 rounded-lg p-4 text-center animate-fade-in-up">
            <p class="text-green-400">Thank you for your feedback! Your insights will help us improve Vega AI for everyone.</p>
          </div>
        `;
        
        // Reset form
        form.reset();
        document.getElementById('difficulty-value').innerText = '5';
        document.getElementById('setup-difficulty').style.setProperty('--value', '50%');
        
        // Remove loading state
        form.classList.remove('htmx-request');
        
        // Hide form after 3 seconds
        setTimeout(() => {
          document.getElementById('feedback-form-container').classList.add('hidden');
          document.getElementById('expand-icon').classList.remove('rotate-180');
          formMessage.innerHTML = '';
        }, 3000);
      }, 1500);
    });
  }
});

// Copy to clipboard functionality is already in main script