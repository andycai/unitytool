-- 系统管理权限
INSERT INTO permissions (id, name, description) VALUES
(1, 'system', '系统管理'),
(2, 'user:list', '用户列表'),
(3, 'user:create', '创建用户'),
(4, 'user:edit', '编辑用户'),
(5, 'user:delete', '删除用户'),
(6, 'role:list', '角色列表'),
(7, 'role:create', '创建角色'),
(8, 'role:edit', '编辑角色'),
(9, 'role:delete', '删除角色'),
(10, 'permission:list', '权限列表'),
(11, 'permission:create', '创建权限'),
(12, 'permission:edit', '编辑权限'),
(13, 'permission:delete', '删除权限'),
(14, 'menu:list', '菜单列表'),
(15, 'menu:create', '创建菜单'),
(16, 'menu:edit', '编辑菜单'),
(17, 'menu:delete', '删除菜单'),
(18, 'adminlog:list', '操作日志列表');

-- 游戏管理权限
INSERT INTO permissions (id, name, description) VALUES
(19, 'game', '游戏管理'),
(20, 'gamelog:list', '游戏日志列表'),
(21, 'stats:list', '统计列表');

-- 系统工具权限
INSERT INTO permissions (id, name, description) VALUES
(22, 'tools', '系统工具'),
(23, 'file:list', '文件列表'),
(24, 'tools:upload', 'FTP上传'),
(25, 'serverconf:list', '服务���配置'),
(26, 'tools:terminal', '命令执行'),
(27, 'package:list', '打包工具'),
(28, 'citask:list', '任务列表'),
(29, 'citask:create', '创建任务'),
(30, 'citask:edit', '编辑任务'),
(31, 'citask:delete', '删除任务'),
(32, 'citask:run', '执行任务'); 