// Get token from URL
const urlParams = new URLSearchParams(window.location.search);
const surveyToken = urlParams.get('token');

if (!surveyToken) {
    window.location.href = '/';
}

// Load and display survey
async function loadSurvey() {
    try {
        const response = await fetch(`/api/surveys/${surveyToken}`);
        if (!response.ok) {
            throw new Error('Survey not found');
        }
        const survey = await response.json();
        renderSurvey(survey);
    } catch (error) {
        document.getElementById('surveyContainer').innerHTML = `
            <div class="error">
                <h2>Error Loading Survey</h2>
                <p>${error.message}</p>
                <a href="/">Return to home</a>
            </div>
        `;
    }
}

function renderSurvey(survey) {
    const container = document.getElementById('surveyContainer');
    container.innerHTML = `
        <h1>${survey.Title}</h1>
        <form id="surveyForm">
            ${survey.Questions.map((question, index) => `
                <div class="question-group" data-question-id="${question.ID}">
                    <div class="question-text">
                        ${index + 1}. ${question.Text}
                        ${question.IsRequired ? '<span class="required">*</span>' : ''}
                    </div>
                    ${renderQuestionInput(question)}
                </div>
            `).join('')}
            <button type="submit">Submit Survey</button>
        </form>
    `;

    document.getElementById('surveyForm').addEventListener('submit', handleSubmit);
}

function renderQuestionInput(question) {
    switch (question.Format) {
        case 'textbox':
            return `<input type="text" ${question.IsRequired ? 'required' : ''}>`;
        case 'multiple_choice':
            return question.Options.map(option => `
                <div class="input-group">
                    <label>
                        <input type="radio" name="q${question.ID}" 
                               value="${option}" ${question.IsRequired ? 'required' : ''}>
                        ${option}
                    </label>
                </div>
            `).join('');
        case 'likert':
            return `<select ${question.IsRequired ? 'required' : ''}>
                ${question.LikertScale.map(value => `
                    <option value="${value}">${value}</option>
                `).join('')}
            </select>`;
        default:
            return '<p>Unsupported question format</p>';
    }
}

async function handleSubmit(e) {
    e.preventDefault();
    
    const answers = Array.from(document.querySelectorAll('.question-group')).map(group => {
        const questionId = group.dataset.questionId;
        const input = group.querySelector('input, select');
        let value;

        if (input.type === 'radio') {
            const selected = group.querySelector('input:checked');
            value = selected ? selected.value : null;
        } else {
            value = input.value;
        }

        return { questionId, value };
    });

    try {
        const response = await fetch('/api/responses', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ surveyToken, answers })
        });

        if (response.ok) {
            alert('Thank you for completing the survey!');
            window.location.href = '/';
        } else {
            throw new Error('Submission failed');
        }
    } catch (error) {
        alert('Failed to submit survey. Please try again.');
    }
}

// Initialize
loadSurvey();