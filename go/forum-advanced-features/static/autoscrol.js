document.addEventListener('DOMContentLoaded', () => {
    // Get parameters from URL query string
    const urlParams = new URLSearchParams(window.location.search);
    let userToHighlight = urlParams.get('u'); // e.g. "You" or "someUsername"
    const commentToHighlight = urlParams.get('c'); // comment ID
    const notificationId = urlParams.get('n'); // notification ID
    
    // Function to clean up URL parameters after scrolling
    const cleanUpUrl = () => {
        const url = new URL(window.location);
        url.searchParams.delete('u');
        url.searchParams.delete('c');
        url.searchParams.delete('n');
        window.history.replaceState({}, document.title, url.toString());
    };
    
    // Handle direct comment scrolling (when c parameter is present)
    if (commentToHighlight) {
        const commentElement = document.getElementById(`comment-${commentToHighlight}`);
        
        if (commentElement) {
            // Scroll the comment into the center of the view smoothly
            commentElement.scrollIntoView({ behavior: 'smooth', block: 'center' });
            // Add a highlight class to make it visually distinct
            commentElement.classList.add('comment-card--highlighted');
            
            // Clean up URL parameters after a short delay to allow scrolling to complete
            setTimeout(cleanUpUrl, 1000);
            return; // Exit early since we found the specific comment
        }
    }
    
    // Handle user-based scrolling (when u parameter is present)
    if (!userToHighlight) {
        return; // no user to highlight, exit early
    }
    
    // If 'u=You', replace it with the actual logged-in username
    // We can get that from a global JS variable or embed it in a data attribute in HTML.
    // For now, let's assume your backend replaces 'You' with the actual username in a hidden span:
    
    if (userToHighlight === 'You') {
      const loggedInUsernameElem = document.getElementById('loggedInUsername');
      
      if (loggedInUsernameElem) {
        userToHighlight = loggedInUsernameElem.textContent.trim();
      } else {
        // fallback: do nothing or exit
        return;
      }
    }
  
    // Now try to find the comment with matching data-username attribute
    const commentElement = document.querySelector(`[data-username="${userToHighlight}"]`);
  
    if (commentElement) {
      // Scroll the comment into the center of the view smoothly
      commentElement.scrollIntoView({ behavior: 'smooth', block: 'center' });
  
      // Optional: add a highlight class to make it visually distinct
      commentElement.classList.add('comment-card--highlighted');
      
      // Clean up URL parameters after a short delay to allow scrolling to complete
      setTimeout(cleanUpUrl, 1000);
    }
  });
  