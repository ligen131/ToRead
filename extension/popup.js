document.addEventListener('DOMContentLoaded', async () => {
  const loginForm = document.getElementById('loginForm');
  const loggedInView = document.getElementById('loggedInView');
  const loginBtn = document.getElementById('loginBtn');
  const logoutBtn = document.getElementById('logoutBtn');
  const saveCustomBtn = document.getElementById('saveCustomBtn');
  const loginMessage = document.getElementById('loginMessage');
  const statusMessage = document.getElementById('statusMessage');
  const userNameDisplay = document.getElementById('userNameDisplay');
  
  // 检查登录状态
  const checkLoginStatus = async () => {
    const response = await chrome.runtime.sendMessage({ action: 'checkLogin' });
    
    if (response.isLoggedIn) {
      const { userName } = await chrome.storage.local.get('userName');
      userNameDisplay.textContent = `欢迎, ${userName || '用户'}`;
      loginForm.style.display = 'none';
      loggedInView.style.display = 'block';
      
      // 自动发送添加当前页面的请求
      await saveCurrentPage();
    } else {
      loginForm.style.display = 'block';
      loggedInView.style.display = 'none';
    }
  };
  
  // 保存当前页面
  const saveCurrentPage = async () => {
    statusMessage.textContent = '正在保存当前页面...';
    statusMessage.className = 'message info';
    
    const result = await chrome.runtime.sendMessage({ action: 'addCurrentPage' });
    
    if (result.success) {
      statusMessage.textContent = '页面已成功保存到ToRead';
      statusMessage.className = 'message success';
    } else {
      statusMessage.textContent = result.message || '保存失败';
      statusMessage.className = 'message error';
    }
  };
  
  // 初始检查登录状态
  await checkLoginStatus();
  
  // 登录按钮点击事件
  loginBtn.addEventListener('click', async () => {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    
    if (!username || !password) {
      loginMessage.textContent = '请输入用户名和密码';
      loginMessage.className = 'message error';
      return;
    }
    
    loginMessage.textContent = '登录中...';
    loginMessage.className = 'message info';
    
    const result = await chrome.runtime.sendMessage({
      action: 'login',
      username,
      password
    });
    
    if (result.success) {
      loginMessage.textContent = '登录成功';
      loginMessage.className = 'message success';
      await checkLoginStatus();
    } else {
      loginMessage.textContent = result.message || '登录失败';
      loginMessage.className = 'message error';
    }
  });
  
  // 退出登录按钮点击事件
  logoutBtn.addEventListener('click', async () => {
    await chrome.storage.local.remove(['token', 'tokenExpirationTime', 'userId', 'userName']);
    await checkLoginStatus();
  });
  
  // 保存自定义链接按钮点击事件
  saveCustomBtn.addEventListener('click', async () => {
    const url = document.getElementById('customUrl').value;
    
    if (!url) {
      statusMessage.textContent = '请输入URL';
      statusMessage.className = 'message error';
      return;
    }
    
    statusMessage.textContent = '正在保存链接...';
    statusMessage.className = 'message info';
    
    const result = await chrome.runtime.sendMessage({
      action: 'addCustomUrl',
      url
    });
    
    if (result.success) {
      statusMessage.textContent = '链接已成功保存到ToRead';
      statusMessage.className = 'message success';
      document.getElementById('customUrl').value = '';
    } else {
      statusMessage.textContent = result.message || '保存失败';
      statusMessage.className = 'message error';
    }
  });
  
  // 添加点击外部关闭弹窗功能
  document.addEventListener('mouseout', (event) => {
    // 检查鼠标是否离开了弹窗区域
    if (!event.relatedTarget || !document.body.contains(event.relatedTarget)) {
      // 添加一个点击事件监听器，当点击其他地方时关闭弹窗
      const closePopup = (e) => {
        window.close();
        document.removeEventListener('click', closePopup);
      };
      
      // 延迟添加点击事件，避免立即触发
      setTimeout(() => {
        document.addEventListener('click', closePopup);
      }, 300);
    }
  });
});
