document.getElementById('submitToken').addEventListener('click', async function(e) {
    e.preventDefault();
    
    const tokenInput = document.getElementById('tokenInput').value;
    const errorMessage = document.getElementById('errorMessage');
    
    try {
        const response = await fetch('/api/post', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: `token=${encodeURIComponent(tokenInput)}`
        });

        const result = await response.text();
        
        if (!response.ok) {
            errorMessage.textContent = result;
            errorMessage.style.color = 'red';
            return;
        }

        errorMessage.textContent = result;
        errorMessage.style.color = 'green';
        
        // Add redirect logic here if needed
        // window.location.href = '/survey-page';

    } catch (error) {
        errorMessage.textContent = `Network error: ${error.message}`;
        errorMessage.style.color = 'red';
    }
});