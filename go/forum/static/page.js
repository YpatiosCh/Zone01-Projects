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

  // Restore scroll position after reload
  if (sessionStorage.getItem('scrollY') !== null) {
    window.scrollTo(0, parseInt(sessionStorage.getItem('scrollY')));
    sessionStorage.removeItem('scrollY');
  }
});