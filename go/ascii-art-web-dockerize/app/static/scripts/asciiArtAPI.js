document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("ascii-art-form");
    const asciiArtContainer = document.getElementById("ascii-art-container");

    form.addEventListener("submit", function (event) {
        event.preventDefault(); // Prevent default form submission

        // Get the input values from the form
        const text = document.getElementById("text").value;
        const banner = document.getElementById("banner").value;

        // Prepare the data to send in the request
        const requestData = {
            text: text,
            banner: banner
        };

        // Make the API call using fetch
        fetch("/api/ascii-art", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(requestData)
        })
        .then(response => response.json())
        .then(data => {
            if (data.ascii_art) {
                // If ASCII art is returned, display it in the container
                asciiArtContainer.innerHTML = "<pre>" + data.ascii_art.join("\n") + "</pre>";
            } else {
                // Handle the case where no ASCII art is returned
                asciiArtContainer.innerHTML = "<p>Error generating ASCII art</p>";
            }
        })
        .catch(error => {
            // Handle any errors that occurred during the fetch
            asciiArtContainer.innerHTML = "<p>Error: " + error.message + "</p>";
        });
    });
});
