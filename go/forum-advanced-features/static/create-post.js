function toggleCategory(button) {
  const categoryId = button.getAttribute('data-category-id');
  const checkbox = document.getElementById('category-' + categoryId);
  
  // Check if button is currently selected
  const isCurrentlySelected = button.classList.contains('selected');
  
  if (!isCurrentlySelected) {
    // Check how many categories are currently selected
    const selectedCategories = document.querySelectorAll('.category-btn.selected');
    
    if (selectedCategories.length >= 3) {
      // Show error message
      showCategoryError("You can choose maximum 3 categories");
      return;
    }
    
    // Remove any existing error message
    removeCategoryError();
    
    // Select the category
    button.classList.add('selected');
    checkbox.checked = true;
  } else {
    // Deselect the category
    button.classList.remove('selected');
    checkbox.checked = false;
    
    // Remove any existing error message when deselecting
    removeCategoryError();
  }
}

function showCategoryError(message) {
  // Remove existing error if any
  removeCategoryError();
  
  // Create error element using existing validation-error class
  const errorDiv = document.createElement('div');
  errorDiv.className = 'validation-error';
  errorDiv.textContent = message;
  
  // Insert error message after the categories selection
  const categoriesSelection = document.querySelector('.categories-selection');
  categoriesSelection.parentNode.insertBefore(errorDiv, categoriesSelection.nextSibling);
}

function removeCategoryError() {
  // Look for validation-error that comes after categories-selection
  const categoriesSelection = document.querySelector('.categories-selection');
  const nextElement = categoriesSelection.nextSibling;
  if (nextElement && nextElement.classList && nextElement.classList.contains('validation-error')) {
    nextElement.remove();
  }
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

  // Handle file input change
  const fileInput = document.getElementById('image');
  const fileNameDisplay = document.getElementById('file-name');
  
  if (fileInput && fileNameDisplay) {
    fileInput.addEventListener('change', function(e) {
      const fileName = e.target.files[0]?.name || 'No file chosen';
      fileNameDisplay.textContent = fileName;
    });
  }
}; 


  let formToDelete = null;

  document.querySelectorAll('.delete__button').forEach((btn) => {
    btn.addEventListener('click', function () {
      const modal = document.getElementById('customConfirm');
      formToDelete = btn.closest('form'); // store the form reference
      modal.style.display = 'block';
    });
  });

  document.getElementById('confirmYes').addEventListener('click', function () {
    if (formToDelete) {
      formToDelete.submit();
    }
  });

  document.getElementById('confirmNo').addEventListener('click', function () {
    document.getElementById('customConfirm').style.display = 'none';
    formToDelete = null;
  });

  // Optional: close modal on background click
  window.onclick = function (event) {
    const modal = document.getElementById('customConfirm');
    if (event.target === modal) {
      modal.style.display = 'none';
      formToDelete = null;
    }
  };