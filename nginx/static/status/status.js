class StatusDashboard {
    constructor() {
        this.updateInterval = 3000;
        this.lastRequestCount = 0;
        this.init();
    }

    init() {
        this.updateTime();
        this.loadStatus();
        setInterval(() => {
            this.updateTime();
            this.loadStatus();
        }, this.updateInterval);
    }

    updateTime() {
        document.getElementById('server-time').textContent = 
            new Date().toLocaleTimeString();
        document.getElementById('last-update').textContent = 
            new Date().toLocaleString();
    }

    async loadStatus() {
        try {
            const response = await fetch('/nginx-status');
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const text = await response.text();
            this.parseStatus(text);
            this.hideError();
        } catch (error) {
            console.error('Error loading status:', error);
            this.showError('Failed to load server status data');
        }
    }

    parseStatus(text) {
        const lines = text.split('\n');
        const stats = {};

        lines.forEach(line => {
            if (line.includes('Active connections')) {
                stats.active = line.match(/\d+/)[0];
            } else if (line.includes('server accepts handled requests')) {
                const numbers = line.match(/\d+/g);
                if (numbers) {
                    stats.accepts = numbers[0];
                    stats.handled = numbers[1];
                    stats.requests = numbers[2];
                }
            } else if (line.includes('Reading')) {
                const numbers = line.match(/\d+/g);
                if (numbers) {
                    stats.reading = numbers[0];
                    stats.writing = numbers[1];
                    stats.waiting = numbers[2];
                }
            }
        });

        this.updateDisplay(stats);
    }

    updateDisplay(stats) {
        document.getElementById('active-conn').textContent = stats.active || '0';
        document.getElementById('reading').textContent = stats.reading || '0';
        document.getElementById('writing').textContent = stats.writing || '0';
        document.getElementById('waiting').textContent = stats.waiting || '0';
        document.getElementById('total-req').textContent = stats.requests || '0';

        if (this.lastRequestCount && stats.requests) {
            const rps = ((stats.requests - this.lastRequestCount) / (this.updateInterval / 1000)).toFixed(1);
            document.getElementById('requests-sec').textContent = rps;
        }
        this.lastRequestCount = stats.requests;
    }

    showError(message) {
        const errorEl = document.getElementById('error-message');
        errorEl.textContent = message;
        errorEl.style.display = 'block';
    }

    hideError() {
        const errorEl = document.getElementById('error-message');
        errorEl.style.display = 'none';
    }
}

document.addEventListener('DOMContentLoaded', () => {
    new StatusDashboard();
});