async function verifyAdmin() {
    // Verify admin token
    const token = document.getElementById('adminToken').value;
    const errorElement = document.getElementById('adminError');

    try {
        const response = await fetch('/api/admin/verify', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ token })
        });

        if (response.ok) {
            window.location.href = '/admin.html';
        } else {
            errorElement.textContent = 'Invalid admin token';
        }
    } catch (error) {
        errorElement.textContent = 'Authentication failed';
    }
}