class SearchHandler {
    constructor() {
        // Initialize references to DOM elements
        this.searchInput = document.getElementById('searchInput');
        this.suggestionsContainer = document.getElementById('suggestions');
        
        // Set up event listeners for user interactions
        this.setupEventListeners();
    }

    // Set up event listeners for the search input and other events
    setupEventListeners() {
        // Listen for input events with debouncing to avoid frequent API calls
        this.searchInput.addEventListener('input', this.debounce((e) => {
            this.fetchSuggestions(e.target.value); // Fetch suggestions when user types
        }, 300)); // Debounce time of 300ms

        // Listen for 'Enter' keypress to trigger search
        this.searchInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                this.performSearch(e.target.value); // Perform search when 'Enter' is pressed
            }
        });

        // Close the suggestions dropdown when clicking outside
        document.addEventListener('click', (e) => {
            if (!this.searchInput.contains(e.target) && !this.suggestionsContainer.contains(e.target)) {
                this.suggestionsContainer.style.display = 'none'; // Hide suggestions if clicked outside
            }
        });

        // Handle click on a suggestion
        this.suggestionsContainer.addEventListener('click', (e) => {
            e.preventDefault(); // Prevent default action (like navigating away immediately)
            const suggestionItem = e.target.closest('.suggestion-item');
            if (suggestionItem) {
                window.location.href = suggestionItem.href; // Navigate to the suggestion's URL
            }
        });
    }

    // Debounce function to delay execution of the search until after user stops typing
    debounce(func, wait) {
        let timeout;
        return function(...args) {
            clearTimeout(timeout); // Clear previous timeout
            timeout = setTimeout(() => func.apply(this, args), wait); // Set a new timeout
        };
    }

    // Fetch search suggestions from the server
    async fetchSuggestions(query) {
        if (!query.trim()) {
            this.suggestionsContainer.style.display = 'none'; // Hide suggestions if input is empty
            return;
        }

        try {
            // Fetch the suggestions from the server
            const response = await fetch(`/suggest?q=${encodeURIComponent(query)}`);
            if (!response.ok) throw new Error('Network response was not ok');
            const suggestions = await response.json();
            this.displaySuggestions(suggestions); // Display the fetched suggestions
        } catch (error) {
            console.error('Error fetching suggestions:', error); // Log any errors
        }
    }

    // Display the fetched suggestions in the UI
    displaySuggestions(suggestions) {
        if (!suggestions.length) {
            this.suggestionsContainer.style.display = 'none'; // Hide suggestions if none are returned
            return;
        }

        // Generate and insert the suggestion items into the suggestions container
        this.suggestionsContainer.innerHTML = suggestions.map(suggestion => `
            <a href="${suggestion.url}" class="suggestion-item">
                <span class="suggestion-text">${suggestion.text}</span>
                <span class="suggestion-type">${suggestion.type}</span>
            </a>
        `).join(''); // Create a list of suggestion items

        this.suggestionsContainer.style.display = 'block'; // Show the suggestions container
    }

    // Perform the search by redirecting to the search results page
    performSearch(query) {
        window.location.href = `/search?q=${encodeURIComponent(query)}`; // Redirect to the search page with the query
    }
}

// Initialize the SearchHandler class when the DOM is fully loaded
document.addEventListener('DOMContentLoaded', () => {
    new SearchHandler(); // Instantiate the search handler class to enable functionality
});
