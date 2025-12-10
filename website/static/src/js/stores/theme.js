
import { getCookie, setCookie } from '../utils/cookie.js';

export default {
    on: false,
    
    init() {
        // 初始化时检查存储的主题或系统偏好
        const cookieTheme = getCookie('idiary-theme');
        if (cookieTheme) {
            this.on = cookieTheme === 'dark';
        } else {
            this.on = window.matchMedia('(prefers-color-scheme: dark)').matches;
        }
        
        // 立即应用
        this.apply();

        // 监听系统主题变化
        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', e => {
            if (!getCookie('idiary-theme')) {
                this.on = e.matches;
                this.apply();
            }
        });
    },

    toggle() {
        this.on = !this.on;
        this.apply();
        setCookie('idiary-theme', this.on ? 'dark' : 'light');
    },

    apply() {
        const html = document.documentElement;
        if (this.on) {
            html.classList.add('dark');
        } else {
            html.classList.remove('dark');
        }
        
        // 更新 meta theme-color
        let metaThemeColor = document.querySelector('meta[name="theme-color"]');
        if (!metaThemeColor) {
            metaThemeColor = document.createElement('meta');
            metaThemeColor.name = 'theme-color';
            document.head.appendChild(metaThemeColor);
        }
        metaThemeColor.content = this.on ? '#1f2937' : '#ffffff';
    }
};
