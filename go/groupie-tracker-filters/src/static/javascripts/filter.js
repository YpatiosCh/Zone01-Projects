document.addEventListener('DOMContentLoaded', () => {

    // Selecting the filter toggle button, sidebar, and content wrapper
    const filterToggleBtn = document.querySelector('.filter-toggle-btn');
    const sidebar = document.querySelector('.filters-sidebar');
    const contentWrapper = document.querySelector('.content-wrapper');
    
    // Create and append an overlay element that will be used for closing the sidebar
    const overlay = document.createElement('div');
    overlay.className = 'filters-overlay';
    document.body.appendChild(overlay);

     // Add event listener to toggle the filter sidebar visibility
     filterToggleBtn.addEventListener('click', () => {
        sidebar.classList.toggle('active');
        overlay.classList.toggle('active'); // Toggle overlay visibility
        contentWrapper.classList.toggle('with-sidebar');
    });

    // Close filters sidebar when clicking outside
    document.addEventListener('click', (e) => {
        if (!sidebar.contains(e.target) && 
            !filterToggleBtn.contains(e.target) && 
            sidebar.classList.contains('active')) {
            sidebar.classList.remove('active');
            overlay.classList.remove('active');
            contentWrapper.classList.remove('with-sidebar');
        }
    });

    // Select elements related to the artist grid and filter form
    const artistGrid = document.querySelector('.artist-grid');
    const filterForm = document.querySelector('.filter-content');
    const noResults = document.querySelector('.no-results');
    let debounceTimer;
    let hasUserInteracted = false;

    // Function to update the displayed results based on the filter form data
    const updateResults = () => {
        const formData = new FormData(filterForm);
        const queryString = new URLSearchParams(formData).toString();
    
        // Show loading state (reduce opacity of the artist grid)
        artistGrid.style.opacity = '0.5';
    
        // Fetch the filtered artist data from the server
        fetch(`/filter?${queryString}`)
            .then(response => response.json())
            .then(artists => {
                // Clear the existing artist grid content
                artistGrid.innerHTML = '';
                
                if (artists.length === 0) {
                    // If no artists were found, display "No results" message
                    artistGrid.style.display = 'none';
                    noResults.style.display = 'block';
                } else {
                    // If artists are found, create and display artist cards
                    const html = artists.map(artist => `
                        <a href="/artist/${artist.id}" class="artist-card">
                            <img src="${artist.image}" alt="${artist.name}" loading="lazy">
                            <div class="card-content">
                                <h2>${artist.name}</h2>
                                <p class="creation-date">Since ${artist.creationDate}</p>
                                <p class="first-album">First Album: ${artist.firstAlbum}</p>
                            </div>
                        </a>
                    `).join('');
                    
                    // Update grid and hide "No results" message
                    artistGrid.style.display = 'grid';
                    noResults.style.display = 'none';
                    artistGrid.innerHTML = html;
                }
                
                // Reset grid opacity after results are updated
                artistGrid.style.opacity = '1';
            })
            .catch(error => {
                console.error('Error:', error);
                // Handle errors by resetting the UI to the "No results" state
                artistGrid.style.opacity = '1';
                artistGrid.innerHTML = '';
                artistGrid.style.display = 'none';
                noResults.style.display = 'block';
            });
    };

    // Handle range inputs with dual slider functionality
    const rangeSections = document.querySelectorAll('.range-filter');
    rangeSections.forEach(section => {
        const minInput = section.querySelector('.range-min');
        const maxInput = section.querySelector('.range-max');
        const [minSpan, maxSpan] = section.querySelectorAll('.range-values span');
        const track = section.querySelector('.slider-track');

        // Function to update the visual track between the two range values
        const updateTrack = () => {
            const min = parseInt(minInput.value);
            const max = parseInt(maxInput.value);
            const percent1 = ((min - minInput.min) / (minInput.max - minInput.min)) * 100;
            const percent2 = ((max - minInput.min) / (minInput.max - minInput.min)) * 100;
            
            // Set the background style of the track based on the range values
            track.style.background = `linear-gradient(to right, 
                var(--primary-color) ${percent1}%, 
                var(--accent-color) ${percent1}%, 
                var(--accent-color) ${percent2}%, 
                var(--primary-color) ${percent2}%)`;
        };

        // Function to update the range values and the associated track when user interacts with input
        const updateRangeValues = () => {
            let minVal = parseInt(minInput.value);
            let maxVal = parseInt(maxInput.value);

            // Swap values if the min value is greater than the max value
            if (minVal > maxVal) {
                if (minInput === document.activeElement) {
                    maxVal = minVal;
                    maxInput.value = maxVal;
                } else {
                    minVal = maxVal;
                    minInput.value = minVal;
                }
            }

            minSpan.textContent = minVal;
            maxSpan.textContent = maxVal;
            updateTrack();
        };

        // Add event listeners to update range values and results on input changes
        minInput.addEventListener('input', () => {
            hasUserInteracted = true;
            updateRangeValues();
            clearTimeout(debounceTimer);
            debounceTimer = setTimeout(updateResults, 300);
        });

        maxInput.addEventListener('input', () => {
            hasUserInteracted = true;
            updateRangeValues();
            clearTimeout(debounceTimer);
            debounceTimer = setTimeout(updateResults, 300);
        });

        // Initialize track and values on page load
        updateRangeValues();
    });

    // Handle location checkboxes and update results on change
    const locationCheckboxes = document.querySelectorAll('input[type="checkbox"]');
    locationCheckboxes.forEach(checkbox => {
        checkbox.addEventListener('change', () => {
            hasUserInteracted = true;
            clearTimeout(debounceTimer);
            debounceTimer = setTimeout(updateResults, 300);
        });
    });

    // Handle country dropdown interactions (toggle cities visibility)
    const countryHeaders = document.querySelectorAll('.country-header');
    countryHeaders.forEach(header => {
        header.addEventListener('click', () => {
            // Toggle the visibility of cities under each country
            header.classList.toggle('active');
            const citiesContainer = header.nextElementSibling;
            citiesContainer.classList.toggle('active');
            
            // Close other countries if they're open
            countryHeaders.forEach(otherHeader => {
                if (otherHeader !== header && otherHeader.classList.contains('active')) {
                    otherHeader.classList.remove('active');
                    otherHeader.nextElementSibling.classList.remove('active');
                }
            });
        });
    });

    // Handle form reset (reset range values and update results)
filterForm.addEventListener('reset', (e) => {
    e.preventDefault(); // Prevent default reset behavior

    // Reset all range inputs to their default values
    rangeSections.forEach(section => {
        const minInput = section.querySelector('.range-min');
        const maxInput = section.querySelector('.range-max');
        const [minSpan, maxSpan] = section.querySelectorAll('.range-values span');
        const track = section.querySelector('.slider-track');
        
        // Reset to initial values
        minInput.value = minInput.min;
        maxInput.value = maxInput.max;
        
        // Update displays
        minSpan.textContent = minInput.value;
        maxSpan.textContent = maxInput.value;
        
        // Update track
        const min = parseInt(minInput.value);
        const max = parseInt(maxInput.value);
        const percent1 = ((min - minInput.min) / (minInput.max - minInput.min)) * 100;
        const percent2 = ((max - minInput.min) / (minInput.max - minInput.min)) * 100;
        
        track.style.background = `linear-gradient(to right, 
            var(--primary-color) ${percent1}%, 
            var(--accent-color) ${percent1}%, 
            var(--accent-color) ${percent2}%, 
            var(--primary-color) ${percent2}%)`;
    });

    // Reset all checkboxes
    const checkboxes = filterForm.querySelectorAll('input[type="checkbox"]');
    checkboxes.forEach(checkbox => {
        checkbox.checked = false;
    });

    // Trigger the filter update
    hasUserInteracted = true;
        updateResults();
    });
});
