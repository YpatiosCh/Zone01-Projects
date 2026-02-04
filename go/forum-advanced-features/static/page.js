document.addEventListener('DOMContentLoaded', function () {
  document.querySelectorAll('.like-form').forEach(form => {
    form.addEventListener('submit', function (e) {
      e.preventDefault();
      // Save scroll position
      const scrollY = window.scrollY;
      fetch(form.action, {
        method: 'POST',
        credentials: 'same-origin'
      })
      .then(res => {
        if (res.ok) {
          // Reload the page and restore scroll position
          window.location.reload();
          // After reload, restore scroll position (see below)
          sessionStorage.setItem('scrollY', scrollY);
        }
      });
    });
  });

  // Check if we should autoscroll to a comment instead of restoring scroll position
  const urlParams = new URLSearchParams(window.location.search);
  const shouldAutoscroll = urlParams.get('u') || urlParams.get('c') || urlParams.get('n');
  
  // Restore scroll position after reload only if we don't need to autoscroll
  if (sessionStorage.getItem('scrollY') !== null && !shouldAutoscroll) {
    window.scrollTo(0, parseInt(sessionStorage.getItem('scrollY')));
    sessionStorage.removeItem('scrollY');
  } else if (sessionStorage.getItem('scrollY') !== null) {
    // Clear the saved scroll position if we're autoscrolling
    sessionStorage.removeItem('scrollY');
  }
});
