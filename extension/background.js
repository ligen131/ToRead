// 后台脚本，处理API请求和登录状态
const API_BASE_URL = 'https://to-read.lg.gl/api/v1';

// 检查登录状态
async function checkLoginStatus() {
  const { token, tokenExpirationTime } = await chrome.storage.local.get(['token', 'tokenExpirationTime']);
  
  if (!token || !tokenExpirationTime) {
    return false;
  }
  
  // 检查token是否过期
  const currentTime = Math.floor(Date.now() / 1000);
  if (currentTime >= tokenExpirationTime) {
    return false;
  }
  
  return true;
}

// 登录函数
async function login(username, password) {
  try {
    const response = await fetch(`${API_BASE_URL}/user/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        user_name: username,
        password: password
      })
    });
    
    const data = await response.json();
    
    if (data.code === 200 && data.data && data.data.token) {
      // 保存token和过期时间
      await chrome.storage.local.set({
        token: data.data.token,
        tokenExpirationTime: data.data.token_expiration_time,
        userId: data.data.user_id,
        userName: data.data.user_name
      });
      return { success: true, data: data.data };
    } else {
      return { success: false, message: data.msg || '登录失败' };
    }
  } catch (error) {
    console.error('Login error:', error);
    return { success: false, message: '网络错误，请稍后再试' };
  }
}

// 添加收藏
async function addCollection(url) {
  try {
    const { token } = await chrome.storage.local.get('token');
    
    if (!token) {
      return { success: false, message: '未登录，请先登录' };
    }
    
    const response = await fetch(`${API_BASE_URL}/collection/add`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({ url })
    });
    
    const data = await response.json();
    
    if (data.code === 200) {
      return { success: true, data: data.data };
    } else {
      return { success: false, message: data.msg || '添加收藏失败' };
    }
  } catch (error) {
    console.error('Add collection error:', error);
    return { success: false, message: '网络错误，请稍后再试' };
  }
}

// 监听来自popup的消息
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.action === 'checkLogin') {
    checkLoginStatus().then(isLoggedIn => {
      sendResponse({ isLoggedIn });
    });
    return true; // 异步响应
  }
  
  if (message.action === 'login') {
    login(message.username, message.password).then(result => {
      sendResponse(result);
    });
    return true;
  }
  
  if (message.action === 'addCurrentPage') {
    chrome.tabs.query({ active: true, currentWindow: true }, async (tabs) => {
      if (tabs.length > 0) {
        const url = tabs[0].url;
        const result = await addCollection(url);
        sendResponse(result);
      } else {
        sendResponse({ success: false, message: '无法获取当前页面URL' });
      }
    });
    return true;
  }
  
  if (message.action === 'addCustomUrl') {
    addCollection(message.url).then(result => {
      sendResponse(result);
    });
    return true;
  }
});
