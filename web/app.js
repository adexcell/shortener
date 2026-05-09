document.addEventListener('DOMContentLoaded', () => {
    const shortenBtn = document.getElementById('shortenBtn');
    const statsBtn = document.getElementById('statsBtn');
    const resultDiv = document.getElementById('result');
    const statsResult = document.getElementById('statsResult');

    // --- ЛОГИКА СОКРАЩЕНИЯ ---
    shortenBtn.addEventListener('click', async () => {
        const originalUrl = document.getElementById('longUrl').value;
        const alias = document.getElementById('alias').value;
        
        if (!originalUrl) {
            showError(resultDiv, 'Пожалуйста, введите URL');
            return;
        }

        shortenBtn.disabled = true;
        shortenBtn.innerHTML = 'Сокращаем...';

        try {
            const response = await fetch('/api/v1/shorten', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ original_url: originalUrl, shorten_code: alias })
            });

            const data = await response.json();
            
            if (response.ok) {
                const shortUrl = `/api/v1/s/${data.shorten_code}`;
                const fullUrl = window.location.origin + shortUrl;
                
                resultDiv.innerHTML = `
                    <div class="result-box animate-up">
                        <p style="margin-bottom: 8px; color: #94a3b8; font-size: 0.9rem;">Ссылка готова:</p>
                        <a href="${shortUrl}" class="result-link" target="_blank">${fullUrl}</a>
                    </div>
                `;
            } else {
                showError(resultDiv, data.error || 'Ошибка при сокращении');
            }
        } catch (e) {
            showError(resultDiv, 'Сервер недоступен. Попробуйте позже.');
        } finally {
            shortenBtn.disabled = false;
            shortenBtn.innerHTML = 'Сократить';
        }
    });

    // --- ЛОГИКА АНАЛИТИКИ ---
    statsBtn.addEventListener('click', async () => {
        let code = document.getElementById('statsCode').value.trim();
        
        if (!code) {
            showError(statsResult, 'Введите код или ссылку');
            return;
        }

        // Извлекаем код если вставили полную ссылку
        if (code.includes('/s/')) {
            code = code.split('/s/').pop();
        }

        statsBtn.disabled = true;
        statsBtn.innerHTML = 'Загрузка...';

        try {
            const response = await fetch(`/api/v1/analytics/${code}`);
            if (!response.ok) throw new Error('Статистика не найдена');
            
            const data = await response.json();
            statsResult.innerHTML = `
                <div class="stats-grid animate-up">
                    <div class="stat-item">
                        <p class="stat-label">Всего переходов</p>
                        <p class="stat-value">${data.total_clicks}</p>
                    </div>
                    <div class="stat-item">
                        <p class="stat-label">Браузеры</p>
                        <pre>${JSON.stringify(data.by_browser, null, 2)}</pre>
                    </div>
                </div>
            `;
        } catch (e) {
            showError(statsResult, e.message);
        } finally {
            statsBtn.disabled = false;
            statsBtn.innerHTML = 'Показать статистику';
        }
    });

    function showError(container, message) {
        container.innerHTML = `
            <div class="error-box animate-up">
                <span>⚠️ ${message}</span>
            </div>
        `;
    }
});
