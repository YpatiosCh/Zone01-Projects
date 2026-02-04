document.addEventListener("DOMContentLoaded", () => {
    document.querySelectorAll('.edit-btn').forEach(button => {
        button.addEventListener('click', () => {
            // Disable all edit buttons while editing
            document.querySelectorAll('.edit-btn').forEach(btn => btn.disabled = true);
            
            const commentId = button.getAttribute('data-comment-id');
            const commentText = document.querySelector(`#comment-${commentId} .comment-card__text`);
            const currentContent = button.getAttribute('data-current-content');
            
            // Create form elements
            const form = document.createElement('form');
            form.action = `/comment/${commentId}/edit`;
            form.method = 'POST';
            form.className = 'inline-edit-form';
            
            const textarea = document.createElement('textarea');
            textarea.name = 'content';
            textarea.required = true;
            textarea.value = currentContent;
            textarea.className = 'inline-edit-textarea';
            
            const buttonGroup = document.createElement('div');
            buttonGroup.className = 'inline-edit-buttons';
            
            const updateBtn = document.createElement('button');
            updateBtn.type = 'submit';
            updateBtn.className = 'btn';
            updateBtn.textContent = 'Update';
            
            const cancelBtn = document.createElement('button');
            cancelBtn.type = 'button';
            cancelBtn.className = 'btn btn--secondary';
            cancelBtn.textContent = 'Cancel';
            
            // Add cancel functionality
            cancelBtn.addEventListener('click', () => {
                form.remove();
                commentText.style.display = 'block';
                // Re-enable all edit buttons after cancellation
                document.querySelectorAll('.edit-btn').forEach(btn => btn.disabled = false);
            });
            
            // Add form submit handler to re-enable buttons after successful update
            form.addEventListener('submit', () => {
                // Re-enable all edit buttons after successful update
                document.querySelectorAll('.edit-btn').forEach(btn => btn.disabled = false);
            });
            
            // Assemble the form
            buttonGroup.appendChild(updateBtn);
            buttonGroup.appendChild(cancelBtn);
            form.appendChild(textarea);
            form.appendChild(buttonGroup);
            
            // Replace the text with the form
            commentText.style.display = 'none';
            commentText.parentNode.insertBefore(form, commentText);
            
            // Focus the textarea
            textarea.focus();
        });
    });
});