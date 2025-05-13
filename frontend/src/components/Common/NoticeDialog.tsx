import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  Box,
  Link
} from '@mui/material';

interface NoticeDialogProps {
  open: boolean;
  onClose: () => void;
}

const NoticeDialog: React.FC<NoticeDialogProps> = ({ open, onClose }) => {
  return (
    <Dialog
      open={open}
      onClose={onClose}
      fullWidth
      maxWidth="md"
    >
      <DialogTitle>使用须知</DialogTitle>
      <DialogContent>
        <Typography variant="body1" paragraph>
          如果你觉得这个项目还不错，求给我一个Star吧！
        </Typography>
        <Box sx={{ mt: 2 }}>
          <Typography variant="body1">
            项目链接：
            <Link 
              href="https://github.com/ligen131/ToRead" 
              target="_blank" 
              rel="noopener noreferrer"
            >
              https://github.com/ligen131/ToRead
            </Link>
          </Typography>
        </Box>

        <Typography variant="body1" paragraph>
          使用方法：
        </Typography>
        <Typography variant="body1" paragraph>
          先注册并登录，用户名不能有特殊字符。登陆进去之后会跳转到收藏列表页面，点击右上角的加号（+）按钮可以添加你想要的收藏链接，可以添加任意链接，比如 https://lg.gl 。点击确认之后可能需要稍等一会，如果使用的人多的话可能需要等待一分钟，如果返回的结果是失败请换一个链接试试，可能链接内容需要通过人机验证才能获取；如果返回成功，直接刷新页面就可以看到已经被解析出来的文章，点击上方的标签按钮可以对文章进行筛选，点击卡片可以跳转原链接。
        </Typography>
        <Typography variant="body1" paragraph>
          ToRead 目前已实现的部分比较简单，前端是用 AI 写的，还存在不少 bug，后端也没有任何的速率限制和安全检查，所以求轻喷。
        </Typography>
        <Typography variant="body1" paragraph>
          如果你也想像ToRead一样丢一个链接给 API，然后 API 给你返回标题+简要描述+标签，你可以像这样请求 
          <Box component="pre" sx={{ 
              backgroundColor: '#f5f5f5', 
              padding: '8px', 
              borderRadius: '4px',
              overflowX: 'auto',
              fontSize: '0.9rem'
            }}>
              {`curl https://g.lg.gl --data '{"url": "https://lg.gl"}'`}
            </Box>
          把 url 字段改成任何你想要的链接就可以了。
        </Typography>
        <Typography variant="body1" paragraph>
          ToRead还有一个浏览器插件可以使用，你可以在上方的开源仓库中找到，在 extension 文件夹中。
        </Typography>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} variant="contained" color="primary">
          确认
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default NoticeDialog;
