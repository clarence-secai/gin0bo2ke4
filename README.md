# gin0bo2ke4
用gin框架开发的前后端分离的个人博客项目。
功能：包括注册登录和token校验、博客的发布修改删除、游客浏览博客、游客评论博客和删除评论。
特点：配置文件的加载、运行时用户上传的管理、规范的错误处理和错误日志、标准的代码分层、gorm与数据库交互、redis缓存。
待改进：实现博客缓存与评论缓存的低耦合，提高缓存的利用率（如不因删一条评论而需重新加载整个博客与评论的缓存）
        向登录用户邮箱发送验证码
