const urlParams = new URLSearchParams(window.location.search);
const surveyToken = urlParams.get('token');

// Load survey responses
async function loadResponses() {
    try {
        const [surveyRes, responsesRes] = await Promise.all([
            fetch(`/api/surveys/${surveyToken}`),
            fetch(`/api/admin/responses/${surveyToken}`)
        ]);

        const survey = await surveyRes.json();
        const responses = await responsesRes.json();

        // Create question map using string IDs
        const questionMap = new Map();
        survey.Questions.forEach(q => {
            // Convert ObjectID to string
            const questionId = q.ID.$oid ? q.ID.$oid : q.ID;
            questionMap.set(questionId, q.Text);
        });

        const container = document.getElementById('responsesContainer');
        container.innerHTML = responses.map((response, index) => `
            <div class="response-container">
                <h3>Response #${index + 1}</h3>
                <p class="meta">
                    Submitted at: ${new Date(response.SubmittedAt).toLocaleString()}
                </p>
                <div class="answers">
                    ${response.Answers.map(answer => {
                        const answerId = answer.QuestionID.$oid ? answer.QuestionID.$oid : answer.QuestionID;
                        return `
                        <div class="answer-item">
                            <strong>${questionMap.get(answerId) || 'Unknown Question'}</strong>
                            <div>${formatAnswerValue(answer.Value)}</div>
                        </div>`;
                    }).join('')}
                </div>
            </div>
        `).join('');
    } catch (error) {
        console.error('Error:', error);
        document.getElementById('responsesContainer').innerHTML = `
            <div class="error">
                Error loading responses: ${error.message}
            </div>
        `;
    }
}

// Format answer value for display
function formatAnswerValue(value) {
    if (Array.isArray(value)) return value.join(', ');
    if (typeof value === 'object') return JSON.stringify(value);
    return value;
}

loadResponses();