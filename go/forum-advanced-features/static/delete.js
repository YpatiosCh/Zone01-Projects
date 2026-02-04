
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