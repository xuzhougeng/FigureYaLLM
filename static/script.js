// DOM Elements
const searchForm = document.getElementById('searchForm');
const searchInput = document.getElementById('searchInput');
const searchButton = document.getElementById('searchButton');
const loadingSpinner = document.getElementById('loadingSpinner');
const resultsSection = document.getElementById('resultsSection');
const errorSection = document.getElementById('errorSection');
const queryDisplay = document.getElementById('queryDisplay');
const explanation = document.getElementById('explanation');
const recommendationsGrid = document.getElementById('recommendationsGrid');
const errorMessage = document.getElementById('errorMessage');

// State
let isLoading = false;

// Event Listeners
searchForm.addEventListener('submit', handleSearch);
searchInput.addEventListener('keydown', function(e) {
    if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault();
        handleSearch(e);
    }
});

// Auto-resize textarea
searchInput.addEventListener('input', function() {
    this.style.height = 'auto';
    this.style.height = Math.min(this.scrollHeight, 120) + 'px';
});

// Fill search input with suggestion
function fillSearch(text) {
    searchInput.value = text;
    searchInput.focus();
    searchInput.style.height = 'auto';
    searchInput.style.height = Math.min(searchInput.scrollHeight, 120) + 'px';
}

// Handle search submission
async function handleSearch(e) {
    e.preventDefault();

    if (isLoading) return;

    const query = searchInput.value.trim();
    if (!query) return;

    setLoading(true);
    hideError();
    hideResults();

    try {
        const response = await fetch('/recommend', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ query })
        });

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }

        const data = await response.json();
        displayResults(data);

    } catch (error) {
        console.error('Search error:', error);
        showError(error.message || '搜索请求失败，请检查网络连接');
    } finally {
        setLoading(false);
    }
}

// Set loading state
function setLoading(loading) {
    isLoading = loading;
    searchButton.disabled = loading;

    if (loading) {
        searchButton.classList.add('loading');
        loadingSpinner.style.display = 'inline-block';
        searchButton.querySelector('.button-text').textContent = '分析中...';
    } else {
        searchButton.classList.remove('loading');
        loadingSpinner.style.display = 'none';
        searchButton.querySelector('.button-text').textContent = 'Ask';
    }
}

// Display search results
function displayResults(data) {
    queryDisplay.textContent = `"${data.query}"`;

    // Display explanation
    if (data.explanation) {
        explanation.innerHTML = `<p class="explanation-text">${escapeHtml(data.explanation)}</p>`;
    }

    // Display recommendations
    recommendationsGrid.innerHTML = '';

    if (data.recommendations && data.recommendations.length > 0) {
        data.recommendations.forEach((rec, index) => {
            const card = createRecommendationCard(rec, index + 1);
            recommendationsGrid.appendChild(card);
        });
    } else {
        recommendationsGrid.innerHTML = '<p class="no-results">未找到相关推荐</p>';
    }

    resultsSection.style.display = 'block';
    resultsSection.scrollIntoView({ behavior: 'smooth', block: 'start' });
}

// Create recommendation card
function createRecommendationCard(rec, rank) {
    const card = document.createElement('div');
    card.className = 'recommendation-card';

    // 计算评分星级
    const stars = getStarRating(rec.score);
    const confidence = getConfidenceLevel(rec.score);

    card.innerHTML = `
        <div class="card-header">
            <div class="card-rank">#${rank}</div>
            <div class="card-title">${escapeHtml(rec.module)}</div>
            <div class="card-score">
                <span class="stars">${stars}</span>
                <span class="score-value">${(rec.score * 100).toFixed(0)}%</span>
            </div>
        </div>

        <div class="card-content">
            <div class="card-field">
                <strong>图片类型：</strong>
                <span class="chart-type">${escapeHtml(rec.chartType || '未指定')}</span>
            </div>

            <div class="card-field">
                <strong>需求描述：</strong>
                <p class="description">${escapeHtml(rec.description || '无描述')}</p>
            </div>

            <div class="card-field">
                <strong>实用场景：</strong>
                <p class="scenario">${escapeHtml(rec.useCase || '无场景描述')}</p>
            </div>

            <div class="card-reason">
                <strong>推荐理由：</strong>
                <p class="reason-text">${escapeHtml(rec.reason || '无推荐理由')}</p>
            </div>
        </div>

        <div class="card-footer">
            <div class="confidence-badge ${confidence.class}">
                ${confidence.text}
            </div>
            <a href="${getModuleUrl(rec.module)}" target="_blank" class="module-link">
                查看模块详情 →
            </a>
        </div>
    `;

    return card;
}

// Get star rating based on score
function getStarRating(score) {
    const fullStars = Math.floor(score * 5);
    const hasHalfStar = (score * 5) % 1 >= 0.5;
    const emptyStars = 5 - fullStars - (hasHalfStar ? 1 : 0);

    return '★'.repeat(fullStars) +
           (hasHalfStar ? '☆' : '') +
           '☆'.repeat(emptyStars);
}

// Get confidence level
function getConfidenceLevel(score) {
    if (score >= 0.9) return { text: '高度匹配', class: 'high' };
    if (score >= 0.7) return { text: '较好匹配', class: 'medium' };
    if (score >= 0.5) return { text: '一般匹配', class: 'low' };
    return { text: '低匹配度', class: 'very-low' };
}

// Get module URL based on module name
function getModuleUrl(moduleName) {
    // Convert FigureYa36nSurvV3 -> https://ying-ge.github.io/FigureYa/FigureYa36nSurvV3/FigureYa36nSurvV3.html
    const baseUrl = 'https://ying-ge.github.io/FigureYa';
    return `${baseUrl}/${moduleName}/${moduleName}.html`;
}

// Show error message
function showError(message) {
    errorMessage.textContent = message;
    errorSection.style.display = 'block';
    errorSection.scrollIntoView({ behavior: 'smooth', block: 'start' });
}

// Hide error message
function hideError() {
    errorSection.style.display = 'none';
}

// Hide results
function hideResults() {
    resultsSection.style.display = 'none';
}

// Escape HTML to prevent XSS
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Initialize page
document.addEventListener('DOMContentLoaded', function() {
    searchInput.focus();

    // Add some visual feedback
    searchInput.addEventListener('focus', function() {
        this.parentElement.classList.add('focused');
    });

    searchInput.addEventListener('blur', function() {
        this.parentElement.classList.remove('focused');
    });
});