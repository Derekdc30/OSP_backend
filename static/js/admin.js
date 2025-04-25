async function loadSurveys() {
    // Load and display all surveys
    try {
        const response = await fetch('/api/admin/surveys');
        const surveys = await response.json();
        
        const list = document.getElementById('surveysList');
        list.innerHTML = surveys.map(survey => `
            <div class="survey-item">
                <h3>${survey.Title} (${survey.Token})</h3>
                <p>Questions: ${survey.Questions.length}</p>
                <p>Created: ${new Date(survey.CreatedAt).toLocaleDateString()}</p>
                <button onclick="editSurvey('${survey.Token}')">Edit</button>
                <button onclick="viewResponses('${survey.Token}')">View Responses</button>
                <button onclick="deleteSurvey('${survey.Token}')">Delete</button>
            </div>
        `).join('');
    } catch (error) {
        console.error('Failed to load surveys:', error);
        alert('Failed to load surveys. Please try again.');
    }
}

async function editSurvey(surveyToken) {
    // Redirect to edit survey page
    window.location.href = `/edit-survey.html?token=${surveyToken}`;
}

async function viewResponses(surveyToken) {
    // Redirect to survey responses page
    window.location.href = `/survey-detail.html?token=${surveyToken}`;
}

async function deleteSurvey(surveyToken) {
    // Delete survey by token
    if (confirm('Are you sure you want to delete this survey?')) {
        try {
            const response = await fetch(`/api/admin/surveys/${surveyToken}`, { method: 'DELETE' });
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`Delete failed: ${response.status} - ${errorText}`);
            }
            loadSurveys();
        } catch (error) {
            console.error('Delete failed:', error);
            alert(`Error deleting survey: ${error.message}`);
        }
    }
}

// Initial load
loadSurveys();