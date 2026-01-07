document.addEventListener('DOMContentLoaded', () => {
    const shortenBtn = document.getElementById('shortenBtn');
    const statsBtn = document.getElementById('statsBtn');

    // --- ЛОГИКА СОКРАЩЕНИЯ ---
    shortenBtn.addEventListener('click', async () => {
        const longUrl = document.getElementById('longUrl').value;
        const alias = document.getElementById('alias').value;
        const resultDiv = document.getElementById('result');

        try {
        const response = await fetch('/shorten', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ url: longUrl, alias: alias })
        });

        const data = await response.json();
        if (response.ok) {

            resultDiv.innerHTML = `Код: ${data.short_url} <br> <a href="/s/${data.short_url}" target="_blank">Перейти</a>`;
        } else {
            resultDiv.innerHTML = `<span class="error">${data.error}</span>`;
        }
    } catch (e) {
        resultDiv.innerHTML = `<span class="error">Сервер недоступен</span>`;
    }
});

    // --- ЛОГИКА АНАЛИТИКИ ---
    statsBtn.addEventListener('click', async () => {
    let code = document.getElementById('statsCode').value.trim();
    
    if (code.includes('/s/')) {
        code = code.split('/s/').pop();
    }

    const statsResult = document.getElementById('statsResult');
    try {
        const response = await fetch(`/analytics/${code}`);
        if (!response.ok) throw new Error('Not found');
        
        const data = await response.json();
        statsResult.innerHTML = `
            <p>Всего кликов: ${data.total_clicks}</p>
            <pre>${JSON.stringify(data.by_browser, null, 2)}</pre>
        `;
    } catch (e) {
        statsResult.innerHTML = `<span class="error">Статистика не найдена</span>`;
    }
  }); 
});
