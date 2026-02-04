function toggleCategory(button) {
  button.classList.toggle('selected');
  
  // Find the associated hidden checkbox and toggle its checked state
  const categoryId = button.getAttribute('data-category-id');
  const checkbox = document.getElementById('category-' + categoryId);
  checkbox.checked = !checkbox.checked;
}

// Initialize selected categories on page load
window.onload = function() {
  const checkboxes = document.querySelectorAll('input[name="categories"]:checked');
  checkboxes.forEach(function(checkbox) {
    const categoryId = checkbox.value;
    const button = document.querySelector(`button[data-category-id="${categoryId}"]`);
    if (button) {
      button.classList.add('selected');
    }
  });
}; 