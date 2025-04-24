document.getElementById('submitToken').addEventListener('click', async () => {
    const token = document.getElementById('tokenInput').value.trim();
    const errorElement = document.getElementById('errorMessage');

    if (!token || token.length !== 5) {
        errorElement.textContent = 'Token must be exactly 5 characters';
        errorElement.style.display = 'block';
        return;
    }

    try {
        const response = await fetch(`/api/check-token`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ token })
        });

        if (response.ok) {
            window.location.href = `survey.html?token=${token}`;
        } else {
            throw new Error('Invalid token');
        }
    } catch (error) {
        errorElement.textContent = 'Invalid token. Please try again.';
        errorElement.style.display = 'block';
    }
});
document.getElementById('tokenInput').addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        document.getElementById('submitToken').click();
    }
});