$(document).ready(function () {
    let lastRequestTime = 0;
    const RATE_LIMIT_MS = 1000; // 1 second

    $('#shortenForm').on('submit', function (e) {
        e.preventDefault();

        const now = Date.now();
        if (now - lastRequestTime < RATE_LIMIT_MS) {
            showError('Please wait a second before shortening another link');
            return;
        }

        const targetUrl = $('#longUrl').val();
        if (!isValidUrl(targetUrl)) {
            showError('Please enter a valid URL');
            return;
        }

        // Reset alerts
        $('.alert').hide();

        // Make API request
        $.ajax({
            url: 'https://api.tnfy.link/v1/links',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ link: { targetUrl } }),
            success: function (response) {
                lastRequestTime = Date.now();
                const shortUrl = response.link.url;
                $('#shortUrl').val(shortUrl);
                $('#resultCard').fadeIn();
            },
            error: function (xhr) {
                let errorMessage = 'An error occurred while shortening the URL';
                if (xhr.responseJSON && xhr.responseJSON.error) {
                    errorMessage = xhr.responseJSON.error;
                }
                showError(errorMessage);
            }
        });
    });

    $('#copyBtn').on('click', function () {
        const shortUrl = $('#shortUrl').val();
        navigator.clipboard.writeText(shortUrl).then(function () {
            const originalText = $('#copyBtn').text();
            $('#copyBtn').text('Copied!');
            setTimeout(() => $('#copyBtn').text(originalText), 2000);
        });
    });

    function showError(message) {
        $('#errorAlert').text(message).fadeIn();
    }

    function isValidUrl(url) {
        try {
            new URL(url);
            return true;
        } catch {
            return false;
        }
    }
});
