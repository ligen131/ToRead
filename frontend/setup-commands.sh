# 创建项目目录
mkdir toread-app
cd toread-app

# 初始化npm项目
npm init -y

# 安装必要的依赖
npm install --save electron electron-builder react react-dom react-router-dom @mui/material @mui/icons-material @emotion/react @emotion/styled axios moment jwt-decode
npm install --save-dev electron-is-dev concurrently wait-on cross-env
npm install --save-dev @types/react @types/react-dom typescript