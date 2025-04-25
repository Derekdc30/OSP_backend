let questions = [];

// Add new question
function addQuestion(type = 'textbox') {
    const question = {
        id: Date.now(),
        text: '',
        format: type,
        options: type === 'multiple_choice' ? ['', ''] : [],
        likertScale: type === 'likert' ? ['1', '2', '3', '4', '5'] : [],
        isRequired: true
    };
    
    questions.push(question);
    renderQuestions();
}

// Update question field
function updateQuestion(id, field, value) {
    const question = questions.find(q => q.id === id);
    if (!question) return;

    // Handle format changes
    if (field === 'format') {
        if (value === 'multiple_choice') {
            question.options = ['', '']; // Initialize with 2 empty options
        } else if (value === 'likert') {
            question.likertScale = ['1', '2', '3', '4', '5']; // Default scale
        }
    }
    
    question[field] = value;
    renderQuestions();
}

// Add multiple-choice option
function addOption(questionId) {
    const question = questions.find(q => q.id === questionId);
    if (question && question.format === 'multiple_choice') {
        question.options.push('');
        renderQuestions();
    }
}

// Update multiple-choice option
function updateOption(questionId, index, value) {
    const question = questions.find(q => q.id === questionId);
    if (question && question.options[index] !== undefined) {
        question.options[index] = value;
        renderQuestions();
    }
}

// Update Likert scale items
function updateLikertScale(questionId, value) {
    const question = questions.find(q => q.id === questionId);
    if (question) {
        question.likertScale = value.split(', ').map(item => item.trim()).filter(Boolean);
        renderQuestions();
    }
}

// Render question editors
function renderQuestions() {
    const container = document.getElementById('questionsContainer');
    container.innerHTML = questions.map((q, index) => `
        <div class="question-editor">
            <h3>Question ${index + 1}</h3>
            <input type="text" value="${q.text}" 
                   onchange="updateQuestion(${q.id}, 'text', this.value)"
                   placeholder="Question text">
            
            <select onchange="updateQuestion(${q.id}, 'format', this.value)">
                <option ${q.format === 'textbox' ? 'selected' : ''}>textbox</option>
                <option ${q.format === 'multiple_choice' ? 'selected' : ''}>multiple_choice</option>
                <option ${q.format === 'likert' ? 'selected' : ''}>likert</option>
            </select>

            ${q.format === 'multiple_choice' ? `
                <div class="options-config">
                    <h4>Multiple Choice Options (minimum 2)</h4>
                    ${q.options.map((opt, i) => `
                        <div class="option-item">
                            <input type="text" 
                                   value="${opt}" 
                                   placeholder="Option ${i + 1}"
                                   onchange="updateOption(${q.id}, ${i}, this.value)">
                            ${i >= 2 ? `<button onclick="question.options.splice(${i}, 1); renderQuestions()">×</button>` : ''}
                        </div>
                    `).join('')}
                    <button class="add-option" onclick="addOption(${q.id})">+ Add Option</button>
                </div>
            ` : ''}

            ${q.format === 'likert' ? `
                <div class="likert-config">
                    <h4>Likert Scale Items (comma-separated)</h4>
                    <input type="text" 
                           value="${q.likertScale.join(', ')}" 
                           placeholder="e.g., 1, 2, 3, 4, 5 or Strongly Disagree, Disagree, Neutral, Agree, Strongly Agree"
                           onchange="updateLikertScale(${q.id}, this.value)">
                </div>
            ` : ''}

            <label class="required-check">
                <input type="checkbox" 
                       ${q.isRequired ? 'checked' : ''}
                       onchange="updateQuestion(${q.id}, 'isRequired', this.checked)">
                Required
            </label>
        </div>
    `).join('');
}

// Save new survey
async function saveSurvey() {
    // Validation
    const errors = [];
    questions.forEach((q, index) => {
        if (q.format === 'multiple_choice') {
            const validOptions = q.options.filter(opt => opt.trim() !== '');
            if (validOptions.length < 2) {
                errors.push(`Question ${index + 1}: Multiple choice needs at least 2 options`);
            }
        }
        if (q.format === 'likert') {
            if (q.likertScale.length < 2) {
                errors.push(`Question ${index + 1}: Likert scale needs at least 2 items`);
            }
        }
    });

    if (errors.length > 0) {
        alert("Validation errors:\n" + errors.join("\n"));
        return;
    }

    // Prepare payload
    const survey = {
        title: document.getElementById('surveyTitle').value.trim(),
        questions: questions.map(q => {
            const base = {
                text: q.text.trim(),
                format: q.format,
                isRequired: q.isRequired
            };
            
            if (q.format === 'multiple_choice') {
                base.options = q.options.filter(opt => opt.trim() !== '');
            }
            
            if (q.format === 'likert') {
                base.likertScale = q.likertScale;
            }
            
            return base;
        })
    };

    try {
        const response = await fetch('/api/admin/surveys', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(survey, (key, value) => {
                if (Array.isArray(value) && value.length === 0) return undefined;
                return value;
            })
        });

        if (response.ok) {
            window.location.href = '/admin.html';
        } else {
            const rawResponse = await response.text();
            console.log('server response:', rawResponse);
            try {
                const error = await JSON.parse(rawResponse);
                throw new Error(error.message || 'Failed to save survey');
            } catch (parseError) {
                throw new Error(`Failed to parse server response: ${rawResponse}`);
            }
        }
    } catch (error) {
        alert(`Error saving survey: ${error.message}`);
    }
}

// Initialize with one question
addQuestion('textbox');