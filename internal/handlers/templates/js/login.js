document.getElementById('loginForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const remember = document.getElementById('remember').checked;
    const errorDiv = document.getElementById('error');
    const loginBtn = document.getElementById('loginBtn');
    const loginText = document.getElementById('loginText');
    const loginLoading = document.getElementById('loginLoading');

    loginBtn.disabled = true;
    loginText.style.display = 'none';
    loginLoading.classList.add('show');
    errorDiv.classList.remove('show');

    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password, remember }),
        });

        const data = await response.json();

        if (response.ok) {
            if (remember) {
                localStorage.setItem('auth_token', data.token);
                localStorage.setItem('auth_username', username);
            } else {
                sessionStorage.setItem('auth_token', data.token);
                sessionStorage.setItem('auth_username', username);
            }

            window.location.href = data.redirect || '/';
        } else {
            errorDiv.textContent = data.message || '登录失败，请检查用户名和密码';
            errorDiv.classList.add('show');

            loginBtn.disabled = false;
            loginText.style.display = 'inline';
            loginLoading.classList.remove('show');

            document.querySelector('.login-container').style.animation = 'none';
            setTimeout(() => {
                document.querySelector('.login-container').style.animation = 'shake 0.5s';
            }, 10);
        }
    } catch (error) {
        errorDiv.textContent = '网络错误，请稍后重试';
        errorDiv.classList.add('show');

        loginBtn.disabled = false;
        loginText.style.display = 'inline';
        loginLoading.classList.remove('show');
    }
});

window.addEventListener('DOMContentLoaded', () => {
    const rememberedUsername = localStorage.getItem('auth_username');
    if (rememberedUsername) {
        document.getElementById('username').value = rememberedUsername;
        document.getElementById('remember').checked = true;
    }
});

document.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        document.getElementById('loginForm').dispatchEvent(new Event('submit'));
    }
});
